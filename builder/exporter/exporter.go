package exporter

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/moby/buildkit/client"
	"go.codecomet.dev/core/filesystem"
)

type exporterType string

const (
	typeImage exporterType = client.ExporterImage
	typeLocal exporterType = client.ExporterLocal
	typeTar   exporterType = client.ExporterTar
	typeOCI   exporterType = client.ExporterOCI

	// Not supported.
	// typeDocker exporterType = client.ExporterDocker.
)

type Entry interface {
	GetEntry() client.ExportEntry
}

type Local struct {
	Path string
	OCI  bool
}

func (o *Local) GetEntry() client.ExportEntry {
	clientExport := client.ExportEntry{
		Type:      string(typeLocal),
		Attrs:     map[string]string{},
		OutputDir: o.Path,
	}

	if o.OCI {
		clientExport.Type = string(typeOCI)
		clientExport.Attrs["tar"] = "false"
	}

	if strings.HasSuffix(o.Path, ".tar") {
		if clientExport.Type == string(typeLocal) {
			clientExport.Type = string(typeTar)
		} else {
			clientExport.Attrs["tar"] = "true"
		}

		clientExport.OutputDir = ""
		clientExport.Output = func(m map[string]string) (io.WriteCloser, error) {
			// XXX this is problematic right now as this is not using the config permissions
			// Still wondering if this should be static instead
			if err := os.MkdirAll(path.Dir(o.Path), filesystem.DirPermissionsDefault); err != nil {
				return nil, err
			}

			ret, err := os.Create(o.Path)
			if err != nil {
				return nil, err
			}

			return io.WriteCloser(ret), nil
		}
	}

	return clientExport
}

type Image struct {
	Name         string
	Push         bool
	PushByDigest bool

	/*
		name=<value>: specify image name(s)
		push=true: push after creating the image
		push-by-digest=true: push unnamed image
		registry.insecure=true: push to insecure HTTP registry
		oci-mediatypes=true: use OCI mediatypes in configuration JSON instead of Docker's
		unpack=true: unpack image after creation (for use with containerd)
		dangling-name-prefix=<value>: name image with prefix@<digest>, used for anonymous images
		name-canonical=true: add additional canonical name name@<digest>
		compression=<uncompressed|gzip|estargz|zstd>: choose compression type for layers newly created and
			cached, gzip is default value. estargz should be used with oci-mediatypes=true.
		compression-level=<value>: compression level for gzip, estargz (0-9) and zstd (0-22)
		force-compression=true: forcefully apply compression option to all layers (including already existing layers)
		store=true: store the result images to the worker's (e.g. containerd) image store as well as ensures that the
			image has all blobs in the content store (default true). Ignored if the worker doesn't have image store
			(e.g. OCI worker).
		annotation.<key>=<value>: attach an annotation with the respective key and value to the built image
		Using the extended syntaxes, annotation-<type>.<key>=<value>, annotation[<platform>].<key>=<value> and both
			combined with annotation-<type>[<platform>].<key>=<value>, allows configuring exactly where to attach the annotation.
		<type> specifies what object to attach to, and can be any of manifest (the default), manifest-descriptor,
			index and index-descriptor
		<platform> specifies which objects to attach to (by default, all), and is the same key passed into the platform
			opt, see docs/multi-platform.md.
		See docs/annotations.md for more details.
	*/
}

func (o *Image) GetEntry() client.ExportEntry {
	return client.ExportEntry{
		Type: string(typeImage),
		Attrs: map[string]string{
			"name":           o.Name,
			"push":           fmt.Sprintf("%t", o.Push),
			"compression":    "estargz",
			"oci-mediatypes": "true",
		},
	}
}
