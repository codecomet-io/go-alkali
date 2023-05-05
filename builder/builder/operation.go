package builder

import (
	"bytes"

	"github.com/codecomet-io/go-alkali/builder"
	"github.com/codecomet-io/go-alkali/builder/build"
	"github.com/codecomet-io/go-alkali/builder/cache"
	"github.com/codecomet-io/go-alkali/builder/exporter"
	"github.com/codecomet-io/go-alkali/builder/registry"
	"github.com/codecomet-io/go-alkali/builder/run"
	"github.com/codecomet-io/go-alkali/machine"
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
