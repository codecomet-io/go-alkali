package commands

import (
	"context"

	"github.com/moby/buildkit/client"
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

	ctx, cancel := context.WithTimeout(ctx, node.ConnectionTimeout)
	defer cancel()

	return client.New(ctx, node.Address.String(), opts...)
}
