package commands

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/codecomet-io/go-alkali/builder/builder"
	"github.com/codecomet-io/go-alkali/builder/locals"
	"github.com/codecomet-io/go-alkali/builder/run"
	"github.com/codecomet-io/go-core/filesystem"
	"github.com/codecomet-io/go-core/log"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	gateway "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/moby/buildkit/solver/pb"
	"github.com/moby/buildkit/util/progress/progresswriter"
	digest "github.com/opencontainers/go-digest"
	"golang.org/x/sync/errgroup"
)

func read(r io.Reader, noCache bool) (*llb.Definition, error) {
	def, err := run.ToDefinition(r)

	if err != nil {
		return nil, err
	}

	if noCache {
		for _, dt := range def.Def {
			var op pb.Op
			if err := (&op).Unmarshal(dt); err != nil {
				return nil, fmt.Errorf("failed to parse llb proto op %w", err)
			}

			dgst := digest.FromBytes(dt)
			opMetadata, ok := def.Metadata[dgst]

			if !ok {
				opMetadata = pb.OpMetadata{}
			}
			c := llb.Constraints{Metadata: opMetadata}
			llb.IgnoreCache(&c)
			def.Metadata[dgst] = c.Metadata
		}
	}

	return def, nil
}

func Build(ctx context.Context, bo *builder.Operation) error {
	// Try and get a client
	cli, err := getClient(bo.Node)
	if err != nil {
		log.Fatal().Err(err).Msg("Builder node is down. Cannot recover.")
	}

	// Get exporters
	exporters := []client.ExportEntry{}
	for _, v := range bo.Export {
		exporters = append(exporters, v.GetEntry())
	}

	// Get SSH and Secrets
	attachable, _ := bo.Options.GetAttachable()
	// Get credentials
	auth := bo.Credentials.GetAttachable()
	attachable = append(attachable, auth...)

	// Create buildkit solve options
	solveOpt := client.SolveOpt{
		Exports:             exporters,
		CacheExports:        bo.Cache.ToClientExport(),
		CacheImports:        bo.Cache.ToClientImport(),
		Session:             attachable,
		AllowedEntitlements: bo.Options.GetEntitlements(),
		Ref:                 bo.Run.ID,
		LocalDirs:           locals.Dump(),
	}

	// Get tracer buffer
	traceEnc := json.NewEncoder(bo.Run.Trace)

	// Get error group
	eg, ctx := errgroup.WithContext(ctx)

	// Read protobuf into a definition
	var def *llb.Definition

	def, err = read(bo.Run.Protobuf, bo.Cache.NoCache)
	if err != nil {
		return err
	}

	if len(def.Def) == 0 {
		return fmt.Errorf("empty definition sent to build")
	}

	// not using shared context to not disrupt display but let is finish reporting errors
	pw, err := progresswriter.NewPrinter(context.TODO(), os.Stderr, bo.Progress)

	if err != nil {
		return err
	}

	if traceEnc != nil {
		traceCh := make(chan *client.SolveStatus)
		pw = progresswriter.Tee(pw, traceCh)
		eg.Go(func() error {
			for s := range traceCh {
				if err := traceEnc.Encode(s); err != nil {
					return err
				}
			}
			return nil
		})
	}

	mw := progresswriter.NewMultiWriter(pw)

	var writers []progresswriter.Writer

	for _, at := range attachable {
		if s, ok := at.(interface {
			SetLogger(progresswriter.Logger)
		}); ok {
			w := mw.WithPrefix("", false)
			s.SetLogger(func(s *client.SolveStatus) {
				w.Status() <- s
			})
			writers = append(writers, w)
		}
	}

	var subMetadata map[string][]byte

	eg.Go(func() error {
		defer func() {
			for _, w := range writers {
				close(w.Status())
			}
		}()

		sreq := gateway.SolveRequest{
			Frontend:    solveOpt.Frontend,
			FrontendOpt: solveOpt.FrontendAttrs,
		}

		if def != nil {
			sreq.Definition = def.ToPB()
		}
		resp, err := cli.Build(ctx, solveOpt, "buildctl", func(ctx context.Context, gwClient gateway.Client) (*gateway.Result, error) {
			_, isSubRequest := sreq.FrontendOpt["requestid"]
			if isSubRequest {
				if _, ok := sreq.FrontendOpt["frontend.caps"]; !ok {
					sreq.FrontendOpt["frontend.caps"] = "moby.buildkit.frontend.subrequests"
				}
			}
			res, err := gwClient.Solve(ctx, sreq)
			if err != nil {
				return nil, err
			}
			if isSubRequest && res != nil {
				subMetadata = res.Metadata
			}
			return res, err
		}, progresswriter.ResetTime(mw.WithPrefix("", false)).Status())
		if err != nil {
			return err
		}
		/*
			for k, v := range resp.ExporterResponse {
				logrus.Debugf("exporter response: %s=%s", k, v)
			}

		*/

		if /*bo.MetadataFile != "" &&*/ resp.ExporterResponse != nil {
			if err := writeMetadataFile("localmeta.json", resp.ExporterResponse); err != nil {
				return err
			}
		}

		return nil
	})

	eg.Go(func() error {
		<-pw.Done()
		return pw.Err()
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	if txt, ok := subMetadata["result.txt"]; ok {
		fmt.Print(string(txt))
	} else {
		for k, v := range subMetadata {
			if strings.HasPrefix(k, "result.") {
				fmt.Printf("%s\n%s\n", k, v)
			}
		}
	}
	return nil
}

func writeMetadataFile(filename string, exporterResponse map[string]string) error {
	var err error
	out := make(map[string]interface{})
	for k, v := range exporterResponse {
		dt, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			out[k] = v
			continue
		}
		var raw map[string]interface{}
		if err = json.Unmarshal(dt, &raw); err != nil || len(raw) == 0 {
			out[k] = v
			continue
		}
		out[k] = json.RawMessage(dt)
	}
	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}
	return filesystem.WriteFile(filename, b, 0o600)
}
