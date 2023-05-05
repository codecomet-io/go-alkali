package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	gateway "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/moby/buildkit/solver/pb"
	"github.com/moby/buildkit/util/progress/progresswriter"
	"go.codecomet.dev/alkali/builder/builder"
	"go.codecomet.dev/containers/digest"
	"golang.org/x/sync/errgroup"
)

var errEmptyDefinition = errors.New("empty definition sent to build")

func read(reader io.Reader, noCache bool) (*llb.Definition, error) {
	// Read it
	byt, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	// Unmarshal into the protobuf definition
	var pbDef pb.Definition
	if err := pbDef.Unmarshal(byt); err != nil {
		return nil, err
	}

	// Convert to the LLB definition
	var def llb.Definition

	def.FromPB(&pbDef)

	if noCache {
		for _, definition := range def.Def {
			var op pb.Op
			if err := (&op).Unmarshal(definition); err != nil {
				return nil, fmt.Errorf("failed to parse llb proto op %w", err)
			}

			dgst := digest.FromBytes(definition)
			opMetadata, ok := def.Metadata[dgst]

			if !ok {
				opMetadata = pb.OpMetadata{}
			}

			c := llb.Constraints{Metadata: opMetadata}
			llb.IgnoreCache(&c)
			def.Metadata[dgst] = c.Metadata
		}
	}

	return &def, nil
}

func Run(ctx context.Context, buildOp *builder.Operation) (map[string]string, error) { //nolint:gocognit
	// Try and get a client
	cli, err := getClient(ctx, buildOp.Node)
	if err != nil {
		return nil, fmt.Errorf("builder node down: %w", err)
	}

	// Get exporters
	exporters := []client.ExportEntry{}
	for _, v := range buildOp.Export {
		exporters = append(exporters, v.GetEntry())
	}

	// Get SSH and Secrets
	attachable, _ := buildOp.Options.GetAttachable()
	// Get credentials
	auth := buildOp.Credentials.GetAttachable()
	attachable = append(attachable, auth...)

	// Create buildkit solve options
	solveOpt := client.SolveOpt{
		Exports:             exporters,
		CacheExports:        buildOp.Cache.ToClientExport(),
		CacheImports:        buildOp.Cache.ToClientImport(),
		Session:             attachable,
		AllowedEntitlements: buildOp.Options.GetEntitlements(),
		Ref:                 buildOp.Run.ID,
		LocalDirs:           buildOp.Run.Locals,
	}

	// Get tracer buffer
	traceEnc := json.NewEncoder(buildOp.Run.Trace)

	// Get error group
	errGroup, ctx := errgroup.WithContext(ctx)

	// Read protobuf into a definition
	var def *llb.Definition

	def, err = read(buildOp.Run.Protobuf, buildOp.Cache.NoCache)
	if err != nil {
		return nil, err
	}

	if len(def.Def) == 0 {
		return nil, errEmptyDefinition
	}

	// not using shared context to not disrupt display but let is finish reporting errors
	progWriter, err := progresswriter.NewPrinter(context.TODO(), os.Stderr, buildOp.Progress) //nolint:contextcheck
	if err != nil {
		return nil, err
	}

	if traceEnc != nil {
		traceCh := make(chan *client.SolveStatus)
		progWriter = progresswriter.Tee(progWriter, traceCh)

		errGroup.Go(func() error {
			for s := range traceCh {
				if err := traceEnc.Encode(s); err != nil {
					return err
				}
			}

			return nil
		})
	}

	multiWriter := progresswriter.NewMultiWriter(progWriter)

	var writers []progresswriter.Writer

	// This will log events from the authenticators and ssh agent
	/*
		for _, at := range attachable {
			if s, ok := at.(interface {
				SetLogger(progresswriter.Logger)
			}); ok {
				prefixedWriter := multiWriter.WithPrefix("", false)

				s.SetLogger(func(s *client.SolveStatus) {
					prefixedWriter.Status() <- s
				})

				writers = append(writers, prefixedWriter)
			}
		}
	*/

	var subMetadata map[string][]byte

	exportResponse := make(map[string]string)

	errGroup.Go(func() error {
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
		resp, err := cli.Build(
			ctx,
			solveOpt,
			"codecomet-alkali",
			func(ctx context.Context, gwClient gateway.Client) (*gateway.Result, error) {
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
			}, progresswriter.ResetTime(multiWriter.WithPrefix("", false)).Status())
		if err != nil {
			return err
		}

		/*
			if resp.ExporterResponse != nil {
				if err := writeMetadataFile("localmeta.json", resp.ExporterResponse); err != nil {
					return err
				}
			}
		*/
		exportResponse = resp.ExporterResponse

		return nil
	})

	errGroup.Go(func() error {
		<-progWriter.Done()

		return progWriter.Err()
	})

	if err := errGroup.Wait(); err != nil {
		return nil, err
	}

	if txt, ok := subMetadata["result.txt"]; ok {
		fmt.Print(string(txt)) //nolint:forbidigo
	} else {
		for k, v := range subMetadata {
			if strings.HasPrefix(k, "result.") {
				fmt.Printf("%s\n%s\n", k, v) //nolint:forbidigo
			}
		}
	}

	return exportResponse, nil
}
