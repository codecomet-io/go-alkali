package build

import (
	"os"

	"github.com/codecomet-io/go-core/log"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/secrets/secretsprovider"
	"github.com/moby/buildkit/session/sshforward/sshprovider"
	"github.com/moby/buildkit/util/entitlements"
)

type Entitlement = entitlements.Entitlement

const (
	entitlementSecurityInsecure Entitlement = entitlements.EntitlementSecurityInsecure
	entitlementNetworkHost      Entitlement = entitlements.EntitlementNetworkHost
)

func New() *Options {
	return &Options{}
}

type Options struct {
	ssh          []sshprovider.AgentConfig
	entitlements []Entitlement

	// XXX should be moved to the main API. This should be a property of Actions or Filesets
	secrets []secretsprovider.Source
}

func (o *Options) AddSSH(id string, paths []string) {
	o.ssh = append(o.ssh, sshprovider.AgentConfig{
		ID:    id,
		Paths: paths,
	})
}

func (o *Options) AddSecret(id string, filepath string, envName string) {
	o.secrets = append(o.secrets, secretsprovider.Source{
		ID:       id,
		FilePath: filepath,
		Env:      envName,
	})
}

func (o *Options) AllowNetworkHost(allow bool) {
	o.entitlements = toggle(o.entitlements, entitlementNetworkHost, allow)
}

func (o *Options) AllowInsecure(allow bool) {
	o.entitlements = toggle(o.entitlements, entitlementSecurityInsecure, allow)
}

func (o *Options) GetEntitlements() []Entitlement {
	return o.entitlements
}

func (o *Options) GetAttachable() ([]session.Attachable, error) {
	attachable := []session.Attachable{}

	if soc := os.Getenv("SSH_AUTH_SOCK"); soc != "" {
		o.AddSSH("default", []string{soc})
	}

	if len(o.ssh) > 0 {
		sp, err := sshprovider.NewSSHAgentProvider(o.ssh)
		if err != nil {
			return nil, err
		}

		attachable = append(attachable, sp)
	} else {
		log.Error().Msg("SSH_AUTH_SOCK is not set and you have not specified any SSH agent configuration. " +
			"SSH forwarding will not work if it is needed (for example, for git+ssh operations). " +
			"Recommend that you start an ssh-agent and set SSH_AUTH_SOCK to the socket path.")
	}

	if len(o.secrets) > 0 {
		store, err := secretsprovider.NewStore(o.secrets)
		if err != nil {
			return nil, err
		}

		attachable = append(attachable, secretsprovider.NewSecretProvider(store))
	}

	return attachable, nil
}

func toggle(slice []Entitlement, value Entitlement, allow bool) []Entitlement {
	position := -1

	for p, v := range slice {
		if v == value {
			position = p

			break
		}
	}
	// No change requested
	if position != -1 && allow || position == -1 && !allow {
		return slice
	}
	// Want it, add it
	if allow {
		return append(slice, value)
	}
	// Otherwise, remove it
	return append(slice[:position], slice[position+1:]...)
}
