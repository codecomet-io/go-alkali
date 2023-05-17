package registry

import (
	"os"

	"github.com/docker/cli/cli/config"
	"github.com/docker/cli/cli/config/configfile"
	"github.com/docker/cli/cli/config/types"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/auth/authprovider"
)

// TODO figure out how to get mTLS for registries.

func New() *Authenticator {
	return &Authenticator{
		dckr: config.LoadDefaultConfigFile(os.Stderr), // &configfile.ConfigFile{},
	}
}

type Authenticator struct {
	dckr *configfile.ConfigFile
}

type Credentials struct {
	ServerAddress string
	Username      string
	Password      string
	Auth          string
}

func (o *Authenticator) Login(auth *Credentials) {
	o.dckr.AuthConfigs[auth.ServerAddress] = types.AuthConfig{
		ServerAddress: auth.ServerAddress,
		Username:      auth.Username,
		Password:      auth.Password,
		Auth:          auth.Auth,
	}
}

func (o *Authenticator) GetAttachable() []session.Attachable {
	return []session.Attachable{
		authprovider.NewDockerAuthProvider(o.dckr),
	}
}

/*
func (o *Authenticator) UseDockerConfig() {
	dockerConfig := config.LoadDefaultConfigFile(os.Stderr)
	for k, v := range dockerConfig.AuthConfigs {
		o.dckr.AuthConfigs[k] = v
	}
}
*/
