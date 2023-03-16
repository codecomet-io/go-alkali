package machine

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/codecomet-io/isovaline/isovaline/config"
)

// XXX home and runner name should be passed down from the implementor side to remove the dependency on isovaline.
func GetSocket() *url.URL {
	sock := os.Getenv("_UNSTABLE_CODECOMET_CUSTOM_BUILDER_SOCKET")
	if sock == "" {
		// sock = lima_cli.New(filepath.Join(config.Get().GetRunRoot(), "vm"), "runner").GetSock()
		sock = fmt.Sprintf("unix://%s/%s/sock/buildkitd.sock", filepath.Join(config.Get().GetRunRoot(), "vm"), "runner")
	}

	u, _ := url.Parse(sock)

	return u
}
