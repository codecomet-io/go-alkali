package commands

import (
	"context"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/codecomet-io/go-alkali/builder/builder"
)

func Info(node *builder.Node, writer io.Writer, format string, ctx context.Context) error {
	client, err := getClient(node)
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

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	_, _ = fmt.Fprintf(w, "BuildKit:\t%s %s %s\n", res.BuildkitVersion.Package, res.BuildkitVersion.Version, res.BuildkitVersion.Revision)
	return w.Flush()
}
