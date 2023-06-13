package commands

import (
	"context"
	"errors"
	"github.com/docker/cli/cli/connhelper/commandconn"
	"github.com/moby/buildkit/client/connhelper"
	"net"
	"net/url"

	"github.com/moby/buildkit/client"
	cst "go.codecomet.dev/alkali/builder"
	"go.codecomet.dev/alkali/builder/builder"
	"go.codecomet.dev/core/telemetry"
)

func getClient(ctx context.Context, node *builder.Node) (*client.Client, error) {
	opts := []client.ClientOpt{
		client.WithFailFast(),
		// TODO: investigate tracing in detail
		// client.WithTracerDelegate
		client.WithTracerProvider(telemetry.GetTracerProvider()),
	}
	if node.CACert != "" || node.Cert != "" || node.Key != "" {
		opts = append(opts, client.WithCredentials(node.Address.Host, node.CACert, node.Cert, node.Key))
	}

	// XXX all bad
	timeout := node.ConnectionTimeout
	if node.ConnectionTimeout == 0 {
		timeout = cst.DefaultConnectionTimeout
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)

	defer cancel()

	return client.New(ctx, node.Address.String(), opts...)
}

// XXX just testing
func init() {
	connhelper.Register("docker-container", Helper)
}

// Helper returns helper for connecting to a Docker container.
// Requires BuildKit v0.5.0 or later in the container.
func Helper(u *url.URL) (*connhelper.ConnectionHelper, error) {
	sp, err := SpecFromURL(u)
	if err != nil {
		return nil, err
	}
	return &connhelper.ConnectionHelper{
		ContextDialer: func(ctx context.Context, addr string) (net.Conn, error) {
			ctxFlags := []string{}
			if sp.Context != "" {
				ctxFlags = append(ctxFlags, "--context="+sp.Context)
			}
			// using background context because context remains active for the duration of the process, after dial has completed
			return commandconn.New(context.Background(), "docker", append(ctxFlags, []string{"exec", "-i", sp.Container, "buildctl", "dial-stdio"}...)...)
		},
	}, nil
}

// Spec
type Spec struct {
	Context   string
	Container string
}

// SpecFromURL creates Spec from URL.
// URL is like docker-container://<container>?context=<context>
// Only <container> part is mandatory.
func SpecFromURL(u *url.URL) (*Spec, error) {
	q := u.Query()
	sp := Spec{
		Context:   q.Get("context"),
		Container: u.Hostname(),
	}
	if sp.Container == "" {
		return nil, errors.New("url lacks container name")
	}
	return &sp, nil
}
