package run

import (
	"bytes"
	"encoding/json"

	"github.com/moby/buildkit/identity"
	"go.codecomet.dev/core/log"
)

func New(proto *bytes.Buffer) *Data {
	return &Data{
		Protobuf: proto,
		ID:       identity.NewID(),
		Trace:    new(bytes.Buffer),
		Meta:     new(bytes.Buffer),
	}
}

type Data struct {
	ID       string
	Trace    *bytes.Buffer
	Protobuf *bytes.Buffer
	Meta     *bytes.Buffer
	Locals   map[string]string
}

func (o *Data) GetJSON() *bytes.Buffer {
	out := new(bytes.Buffer)

	if o.Protobuf == nil {
		log.Fatal().Msg("Uninitialized protobuf buffer")
	}

	ops, err := readLLB(o.Protobuf)
	if err != nil {
		log.Fatal().Msg("Failed reading protobuf")
	}

	enc := json.NewEncoder(out)

	for _, op := range ops {
		if err := enc.Encode(op); err != nil {
			log.Fatal().Msg("Failed json encoding op")
		}
	}

	return out
}

func (o *Data) GetDOT() *bytes.Buffer {
	out := new(bytes.Buffer)

	if o.Protobuf == nil {
		log.Fatal().Msg("Uninitialized protobuf buffer")
	}

	_ = toDOT(o.Protobuf, out)

	return out
}
