package builder

import (
	"github.com/opencontainers/go-digest"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

// A Digest is... a Digest... which is a Digest.
type Digest = digest.Digest

type Platform = specs.Platform
