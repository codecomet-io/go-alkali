package commands

import (
	"context"
	"fmt"
	"io"
	"path"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/codecomet-io/go-alkali/builder/builder"
	"github.com/moby/buildkit/client"
	ocispecs "github.com/opencontainers/image-spec/specs-go/v1"
)

const defaultTabWidth = 8

func Workers(ctx context.Context, node *builder.Node, writer io.Writer, verbose bool,
	format string, filter []string,
) error {
	cli, err := getClient(ctx, node)
	if err != nil {
		return err
	}

	workers, err := cli.ListWorkers(ctx, client.WithFilter(filter))
	if err != nil {
		return err
	}

	if format != "" {
		tmpl, err := parseTemplate(format)
		if err != nil {
			return err
		}

		if err := tmpl.Execute(writer, workers); err != nil {
			return err
		}

		_, err = fmt.Fprintf(writer, "\n")

		return err
	}

	tabWriter := tabwriter.NewWriter(writer, 1, defaultTabWidth, 1, '\t', 0)

	if verbose {
		printWorkersVerbose(tabWriter, workers)
	} else {
		printWorkersTable(tabWriter, workers)
	}

	return nil
}

func printWorkersVerbose(tabWriter *tabwriter.Writer, winfo []*client.WorkerInfo) {
	for _, workerInfo := range winfo {
		fmt.Fprintf(tabWriter, "ID:\t%s\n", workerInfo.ID)
		fmt.Fprintf(tabWriter, "Platforms:\t%s\n", joinPlatforms(workerInfo.Platforms))
		fmt.Fprintf(tabWriter, "BuildKit:\t%s %s %s\n", workerInfo.BuildkitVersion.Package,
			workerInfo.BuildkitVersion.Version, workerInfo.BuildkitVersion.Revision)
		fmt.Fprintf(tabWriter, "Labels:\n")

		for _, k := range sortedKeys(workerInfo.Labels) {
			v := workerInfo.Labels[k]
			fmt.Fprintf(tabWriter, "\t%s:\t%s\n", k, v)
		}

		for i, rule := range workerInfo.GCPolicy {
			fmt.Fprintf(tabWriter, "GC Policy rule#%d:\n", i)
			fmt.Fprintf(tabWriter, "\tAll:\t%v\n", rule.All)

			if len(rule.Filter) > 0 {
				fmt.Fprintf(tabWriter, "\tFilters:\t%s\n", strings.Join(rule.Filter, " "))
			}

			if rule.KeepDuration > 0 {
				fmt.Fprintf(tabWriter, "\tKeep Duration:\t%v\n", rule.KeepDuration.String())
			}

			if rule.KeepBytes > 0 {
				fmt.Fprintf(tabWriter, "\tKeep Bytes:\t%d\n", rule.KeepBytes)
			}
		}

		fmt.Fprintf(tabWriter, "\n")
	}

	tabWriter.Flush()
}

func printWorkersTable(tabWriter *tabwriter.Writer, winfo []*client.WorkerInfo) {
	fmt.Fprintln(tabWriter, "ID\tPLATFORMS")

	for _, workerInfo := range winfo {
		id := workerInfo.ID
		fmt.Fprintf(tabWriter, "%s\t%s\n", id, joinPlatforms(workerInfo.Platforms))
	}

	tabWriter.Flush()
}

func sortedKeys(m map[string]string) []string {
	sorted := make([]string, len(m))
	i := 0

	for k := range m {
		sorted[i] = k
		i++
	}

	sort.Strings(sorted)

	return sorted
}

func joinPlatforms(p []ocispecs.Platform) string {
	str := make([]string, 0, len(p))
	for _, pp := range p {
		// XXX normalizing and formatting here makes no sense
		// str = append(str, platforms.Format(normalize(pp)))
		str = append(str, path.Join(pp.OS, pp.Architecture, pp.Variant))
	}

	return strings.Join(str, ",")
}

/*
func format(platform ocispecs.Platform) string {
	if platform.OS == "" {
		return "unknown"
	}

	return path.Join(platform.OS, platform.Architecture, platform.Variant)
}

// XXX obviously the host platform is NOT necessarily the daemon platform, so, this stuff here is just faulty
func normalize(platform ocispecs.Platform) ocispecs.Platform {
	platform.OS = normalizeOS(platform.OS)
	platform.Architecture, platform.Variant = normalizeArch(platform.Architecture, platform.Variant)

	return platform
}

func normalizeOS(os string) string {
	if os == "" {
		return runtime.GOOS
	}
	os = strings.ToLower(os)

	switch os {
	case "macos":
		os = "darwin"
	}
	return os
}

func normalizeArch(arch, variant string) (string, string) {
	arch, variant = strings.ToLower(arch), strings.ToLower(variant)
	switch arch {
	case "i386":
		arch = "386"
		variant = ""
	case "x86_64", "x86-64", "amd64":
		arch = "amd64"
		if variant == "v1" {
			variant = ""
		}
	case "aarch64", "arm64":
		arch = "arm64"
		switch variant {
		case "8", "v8":
			variant = ""
		}
	case "armhf":
		arch = "arm"
		variant = "v7"
	case "armel":
		arch = "arm"
		variant = "v6"
	case "arm":
		switch variant {
		case "", "7":
			variant = "v7"
		case "5", "6", "8":
			variant = "v" + variant
		}
	}

	return arch, variant
}


*/
