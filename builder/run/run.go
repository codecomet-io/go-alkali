package run

import (
	"bytes"

	"github.com/codecomet-io/go-alkali/builder/wrapllb"
	"github.com/codecomet-io/go-core/log"
	"github.com/moby/buildkit/identity"
)

func New(proto *bytes.Buffer /*state llb.State*/) *Data {
	bd := &Data{
		// state: state,

		Protobuf: proto, // new(bytes.Buffer),
		ID:       identity.NewID(),
		Trace:    new(bytes.Buffer),
		Meta:     new(bytes.Buffer),
	}
	// _ = wrapllb.Write(state, bd.Protobuf)
	return bd
}

type Data struct {
	// state llb.State

	ID       string
	Trace    *bytes.Buffer
	Protobuf *bytes.Buffer
	Meta     *bytes.Buffer
	// JSON     io.ReadWriter
}

func (o *Data) GetJSON() *bytes.Buffer {
	out := new(bytes.Buffer)
	if o.Protobuf == nil {
		log.Fatal().Msg("Uninitialized protobuf buffer")
	}
	_ = wrapllb.ToJSON(o.Protobuf, out)
	return out
}

func (o *Data) GetDOT() *bytes.Buffer {
	out := new(bytes.Buffer)
	if o.Protobuf == nil {
		log.Fatal().Msg("Uninitialized protobuf buffer")
	}
	_ = wrapllb.ToDOT(o.Protobuf, out)
	return out
}
