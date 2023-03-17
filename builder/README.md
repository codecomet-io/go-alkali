# Builder

In a shell: "a toolkit for developing buildkitd clients".

This is evidently built on top of https://github.com/moby/buildkit and 
reuses code from buildctl.

## Example

```go
package main

import (
	"context"
	"bytes"

	"github.com/codecomet-io/go-core/log"
	"github.com/codecomet-io/go-alkali/builder/builder"
	"github.com/codecomet-io/go-alkali/builder/commands"
	"github.com/codecomet-io/go-alkali/builder/locals"
)

func main() {
	// A protobuf message (eg: a marshalled llb.State)
	var proto *bytes.Buffer

	locals.Reset()
	bo := builder.NewOperation()
	bo.Ingest(proto)

	err := commands.Run(context.Background(), bo)
	if err != nil {
		log.Error().Err(err).Msg("failed to run pipeline")
	}
}

```

## Caveats

Current design is work in progress, specifically the `BuildOperation`
top-level struct. `Run` and `Export` clearly belong to a notion of one operation,
while the rest belongs to a `controller`.
Furthermore, `Secrets` (inside `Options`) is also run-dependent.
