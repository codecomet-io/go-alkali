package machine

import (
	"net/url"
	"os"
)

func GetSocket(path string) *url.URL {
	sock := os.Getenv("_UNSTABLE_CODECOMET_CUSTOM_BUILDER_SOCKET")
	if sock == "" {
		// Relationship to the way isovaline creates the VM is finicky
		// But then, shelling out is not necessarily good either
		sock = path
	}

	u, _ := url.Parse(sock)

	return u
}
