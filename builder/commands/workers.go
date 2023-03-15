package commands

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/codecomet-io/go-alkali/builder/builder"
	"github.com/containerd/containerd/platforms"
	"github.com/moby/buildkit/client"
	ocispecs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/tonistiigi/units"
)

func Workers(node *builder.Node, w io.Writer, ctx context.Context, verbose bool, format string, filter []string) error {
	c, err := getClient(node)
	if err != nil {
		return err
	}
	workers, err := c.ListWorkers(ctx, client.WithFilter(filter))
	if err != nil {
		return err
	}
	if format != "" {
		tmpl, err := parseTemplate(format)
		if err != nil {
			return err
		}
		if err := tmpl.Execute(w, workers); err != nil {
			return err
		}
		_, err = fmt.Fprintf(w, "\n")
		return err
	}

	tw := tabwriter.NewWriter(w, 1, 8, 1, '\t', 0)

	if verbose {
		printWorkersVerbose(tw, workers)
	} else {
		printWorkersTable(tw, workers)
	}
	return nil
}

func printWorkersVerbose(tw *tabwriter.Writer, winfo []*client.WorkerInfo) {
	for _, wi := range winfo {
		fmt.Fprintf(tw, "ID:\t%s\n", wi.ID)
		fmt.Fprintf(tw, "Platforms:\t%s\n", joinPlatforms(wi.Platforms))
		fmt.Fprintf(tw, "BuildKit:\t%s %s %s\n", wi.BuildkitVersion.Package, wi.BuildkitVersion.Version, wi.BuildkitVersion.Revision)
		fmt.Fprintf(tw, "Labels:\n")
		for _, k := range sortedKeys(wi.Labels) {
			v := wi.Labels[k]
			fmt.Fprintf(tw, "\t%s:\t%s\n", k, v)
		}
		for i, rule := range wi.GCPolicy {
			fmt.Fprintf(tw, "GC Policy rule#%d:\n", i)
			fmt.Fprintf(tw, "\tAll:\t%v\n", rule.All)
			if len(rule.Filter) > 0 {
				fmt.Fprintf(tw, "\tFilters:\t%s\n", strings.Join(rule.Filter, " "))
			}
			if rule.KeepDuration > 0 {
				fmt.Fprintf(tw, "\tKeep Duration:\t%v\n", rule.KeepDuration.String())
			}
			if rule.KeepBytes > 0 {
				fmt.Fprintf(tw, "\tKeep Bytes:\t%g\n", units.Bytes(rule.KeepBytes))
			}
		}
		fmt.Fprintf(tw, "\n")
	}

	tw.Flush()
}

func printWorkersTable(tw *tabwriter.Writer, winfo []*client.WorkerInfo) {
	fmt.Fprintln(tw, "ID\tPLATFORMS")

	for _, wi := range winfo {
		id := wi.ID
		fmt.Fprintf(tw, "%s\t%s\n", id, joinPlatforms(wi.Platforms))
	}

	tw.Flush()
}

func sortedKeys(m map[string]string) []string {
	s := make([]string, len(m))
	i := 0
	for k := range m {
		s[i] = k
		i++
	}
	sort.Strings(s)
	return s
}

func joinPlatforms(p []ocispecs.Platform) string {
	str := make([]string, 0, len(p))
	for _, pp := range p {
		str = append(str, platforms.Format(platforms.Normalize(pp)))
	}
	return strings.Join(str, ",")
}
