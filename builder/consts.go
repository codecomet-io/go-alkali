package builder

import "time"

const DefaultConnectionTimeout = 10 * time.Second

// XXX make this configurable so that we can honor config from embedders.
const (
	DefaultDirPerms  = 0o700
	DefaultFilePerms = 0o600
)
