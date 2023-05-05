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

	"go.codecomet.dev/core/log"
	"go.codecomet.dev/alkali/builder/builder"
	"go.codecomet.dev/alkali/builder/commands"
	"go.codecomet.dev/alkali/builder/locals"
)

func main() {
	// A protobuf message (eg: a marshalled llb.State)
	var proto *bytes.Buffer

	locals.Reset()
	bo := builder.NewOperation("socket_path")
	bo.Ingest(proto)

	_, err := commands.Run(context.Background(), bo)
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

### About init

Unfortunately, buildkit does:
* use environment variables as the only mean to customize certain things (eg: `BUILDKIT_COLORS`)
* also buildkit is setting these from the environment inside an `init()` function

This makes it impossible to override them, unless we manipulate these env variables *before* the `init()` function is called.
Given go import ordering...
