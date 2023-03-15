package machine

import (
	"net/url"
	"os"
	"path/filepath"

	// XXX Need to rip these out somehow and get them here.
	"github.com/codecomet-io/isovaline/isovaline/config"
	lima_cli "github.com/codecomet-io/isovaline/upstream/lima-cli"
)

func GetSocket() *url.URL {
	sock := os.Getenv("_UNSTABLE_CODECOMET_CUSTOM_BUILDER_SOCKET")
	if sock == "" {
		sock = lima_cli.New(filepath.Join(config.Get().GetRunRoot(), "vm"), "runner").GetSock()
	}

	u, _ := url.Parse(sock)

	return u
}
