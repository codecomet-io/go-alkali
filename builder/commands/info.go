package commands

import (
	"context"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"go.codecomet.dev/alkali/builder/builder"
)

func Info(ctx context.Context, node *builder.Node, writer io.Writer, format string) error {
	client, err := getClient(ctx, node)
	if err != nil {
		return err
	}

	res, err := client.Info(ctx)
	if err != nil {
		return err
	}

	if format != "" {
		tmpl, err := parseTemplate(format)
		if err != nil {
			return err
		}

		if err := tmpl.Execute(writer, res); err != nil {
			return err
		}

		_, err = fmt.Fprintf(writer, "\n")

		return err
	}

	tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	_, _ = fmt.Fprintf(tabWriter, "BuildKit:\t%s %s %s\n",
		res.BuildkitVersion.Package,
		res.BuildkitVersion.Version,
		res.BuildkitVersion.Revision)

	return tabWriter.Flush()
}
