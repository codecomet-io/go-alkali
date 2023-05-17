package run

import (
	"fmt"
	"io"
	"strings"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/solver/pb"
	"go.codecomet.dev/containers/digest"
)

type llbOp struct {
	Op         pb.Op         `json:"op"`
	Digest     digest.Digest `json:"digest"`
	OpMetadata pb.OpMetadata `json:"opMetadata"`
}

func readLLB(reader io.Reader) ([]llbOp, error) {
	// Read it
	byt, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	// Unmarshal into the protobuf definition
	var pbDef pb.Definition
	if err := pbDef.Unmarshal(byt); err != nil {
		return nil, err
	}
	// Convert to the LLB definition
	var def llb.Definition

	def.FromPB(&pbDef)

	// Stuff everything in the ad-hoc llbOp struct
	ops := []llbOp{}

	for _, definition := range def.Def {
		var operation pb.Op
		if err := (&operation).Unmarshal(definition); err != nil {
			return nil, fmt.Errorf("failed to parse op %w", err)
		}

		dgst := digest.FromBytes(definition)
		ent := llbOp{Op: operation, Digest: dgst, OpMetadata: def.Metadata[dgst]}
		ops = append(ops, ent)
	}

	return ops, nil
}

func toDOT(reader io.Reader, writer io.Writer) error {
	ops, err := readLLB(reader)
	if err != nil {
		return err
	}

	// TODO: print OpMetadata
	fmt.Fprintln(writer, "digraph {")
	defer fmt.Fprintln(writer, "}")

	for _, op := range ops {
		name, shape := attr(op.Digest, op.Op)
		fmt.Fprintf(writer, "  %q [label=%q shape=%q];\n", op.Digest, name, shape)
	}

	for _, operation := range ops {
		for i, inp := range operation.Op.Inputs {
			label := ""

			if eo, ok := operation.Op.Op.(*pb.Op_Exec); ok {
				for _, m := range eo.Exec.Mounts {
					if int(m.Input) == i && m.Dest != "/" {
						label = m.Dest
					}
				}
			}

			fmt.Fprintf(writer, "  %q -> %q [label=%q];\n", inp.Digest, operation.Digest, label)
		}
	}

	return nil
}

func attr(dgst digest.Digest, op pb.Op) (string, string) {
	switch operation := op.Op.(type) {
	case *pb.Op_Source:
		return operation.Source.Identifier, "ellipse"
	case *pb.Op_Exec:
		return strings.Join(operation.Exec.Meta.Args, " "), "box"
	case *pb.Op_Build:
		return "build", "box3d"
	case *pb.Op_Merge:
		return "merge", "invtriangle"
	case *pb.Op_Diff:
		return "diff", "doublecircle"
	case *pb.Op_File:
		names := []string{}

		for _, action := range operation.File.Actions {
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
