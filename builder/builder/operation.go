package builder

import (
	"bytes"

	"go.codecomet.dev/alkali/builder"
	"go.codecomet.dev/alkali/builder/build"
	"go.codecomet.dev/alkali/builder/cache"
	"go.codecomet.dev/alkali/builder/exporter"
	"go.codecomet.dev/alkali/builder/registry"
	"go.codecomet.dev/alkali/builder/run"
	"go.codecomet.dev/alkali/machine"
)

type Operation struct {
	Node        *Node
	Credentials *registry.Authenticator
	Cache       *cache.Options
	Export      []exporter.Entry
	Options     *build.Options
	Run         *run.Data

	// XXX
	Progress string
}

func (o *Operation) Ingest(proto *bytes.Buffer) {
	o.Run = run.New(proto)
}

func NewOperation(path string) *Operation {
	return &Operation{
		Node: &Node{
			ConnectionTimeout: builder.DefaultConnectionTimeout,
			Address:           machine.GetSocket(path),
		},
		Credentials: registry.New(),
		Options:     build.New(),
		Cache: &cache.Options{
			Export: []cache.Entry{},
			Import: []cache.Entry{},
		},
		Export: []exporter.Entry{},
		Run:    run.New(nil),
	}
}
