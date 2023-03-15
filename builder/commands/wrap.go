package commands

import (
	"context"

	"github.com/codecomet-io/go-alkali/builder/builder"
	"github.com/moby/buildkit/client"
)

func getClient(node *builder.Node) (*client.Client, error) {
	opts := []client.ClientOpt{client.WithFailFast()}
	// TODO: investigate tracing here
	// client.WithTracerProvider, client.WithTracerDelegate
	if node.CACert != "" || node.Cert != "" || node.Key != "" {
		opts = append(opts, client.WithCredentials(node.Address.Host, node.CACert, node.Cert, node.Key))
	}
	ctx, cancel := context.WithTimeout(context.TODO(), node.ConnectionTimeout)
	defer cancel()

	return client.New(ctx, node.Address.String(), opts...)
}
