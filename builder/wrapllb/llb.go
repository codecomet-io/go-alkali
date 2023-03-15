package wrapllb

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/solver/pb"
	"github.com/opencontainers/go-digest"
)

type llbOp struct {
	Op         pb.Op
	Digest     digest.Digest
	OpMetadata pb.OpMetadata
}

func ToDefinition(r io.Reader) (*llb.Definition, error) {
	// Read it
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	// Unmarshal into the protobuf definition
	var pbDef pb.Definition
	if err := pbDef.Unmarshal(b); err != nil {
		return nil, err
	}
	// Convert to the LLB definition
	var def llb.Definition
	def.FromPB(&pbDef)
	return &def, nil
}

func readLLB(r io.Reader) ([]llbOp, error) {
	def, err := ToDefinition(r)
	if err != nil {
		return nil, err
	}
	// Stuff everything in the ad-hoc llbOp struct
	var ops []llbOp
	for _, dt := range def.Def {
		var op pb.Op
		if err := (&op).Unmarshal(dt); err != nil {
			return nil, fmt.Errorf("failed to parse op %w", err)
		}
		dgst := digest.FromBytes(dt)
		ent := llbOp{Op: op, Digest: dgst, OpMetadata: def.Metadata[dgst]}
		ops = append(ops, ent)
	}
	return ops, nil
}

func ToJSON(r io.Reader, w io.Writer) error {
	ops, err := readLLB(r)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(w)
	for _, op := range ops {
		if err := enc.Encode(op); err != nil {
			return err
		}
	}
	return nil
}

func ToDOT(r io.Reader, w io.Writer) error {
	ops, err := readLLB(r)
	if err != nil {
		return err
	}

	// TODO: print OpMetadata
	fmt.Fprintln(w, "digraph {")
	defer fmt.Fprintln(w, "}")
	for _, op := range ops {
		name, shape := attr(op.Digest, op.Op)
		fmt.Fprintf(w, "  %q [label=%q shape=%q];\n", op.Digest, name, shape)
	}
	for _, op := range ops {
		for i, inp := range op.Op.Inputs {
			label := ""
			if eo, ok := op.Op.Op.(*pb.Op_Exec); ok {
				for _, m := range eo.Exec.Mounts {
					if int(m.Input) == i && m.Dest != "/" {
						label = m.Dest
					}
				}
			}
			fmt.Fprintf(w, "  %q -> %q [label=%q];\n", inp.Digest, op.Digest, label)
		}
	}
	return nil
}

func attr(dgst digest.Digest, op pb.Op) (string, string) {
	switch op := op.Op.(type) {
	case *pb.Op_Source:
		return op.Source.Identifier, "ellipse"
	case *pb.Op_Exec:
		return strings.Join(op.Exec.Meta.Args, " "), "box"
	case *pb.Op_Build:
		return "build", "box3d"
	case *pb.Op_Merge:
		return "merge", "invtriangle"
	case *pb.Op_Diff:
		return "diff", "doublecircle"
	case *pb.Op_File:
		names := []string{}

		for _, action := range op.File.Actions {
			var name string

			switch act := action.Action.(type) {
			case *pb.FileAction_Copy:
				name = fmt.Sprintf("copy{src=%s, dest=%s}", act.Copy.Src, act.Copy.Dest)
			case *pb.FileAction_Mkfile:
				name = fmt.Sprintf("mkfile{path=%s}", act.Mkfile.Path)
			case *pb.FileAction_Mkdir:
				name = fmt.Sprintf("mkdir{path=%s}", act.Mkdir.Path)
			case *pb.FileAction_Rm:
				name = fmt.Sprintf("rm{path=%s}", act.Rm.Path)
			}

			names = append(names, name)
		}
		return strings.Join(names, ","), "note"
	default:
		return dgst.String(), "plaintext"
	}
}
