package plugin

import (
	"fmt"
	"os/exec"
	"strings"

	v1 "github.com/VJftw/docker-registry-proxy/api/proto/v1"
	"github.com/hashicorp/go-plugin"
)

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	// This isn't required when using VersionedPlugins
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "DH*78Y@XgUXL54",
}

const (
	PluginTypeAuthProvider  = "auth-provider"
	PluginTypeAuthVerifier  = "auth-verifier"
	PluginTypeConfiguration = "configuration"
)

var PluginSearchPath = "./plugins"

// PluginMap is the map of plugins we can dispense.
var PluginMap = map[string]plugin.Plugin{
	PluginTypeAuthProvider:  &AuthProviderGRPCPlugin{},
	PluginTypeAuthVerifier:  &AuthVerifierGRPCPlugin{},
	PluginTypeConfiguration: &ConfigurationGRPCPlugin{},
}

var PluginClients = map[string]*plugin.Client{}

// Close cleans up and kills all of the plugin processes
func Close() {
	for _, c := range PluginClients {
		c.Kill()
	}
}

func LoadPluginsFromConfigSlice(configs []string) error {
	for _, pluginConfig := range configs {
		ptype, name, alias, err := ResolvePluginTypeNameAndAlias(pluginConfig)
		if err != nil {
			return err
		}
		pluginPath, err := FindPluginInPath(PluginSearchPath, name, ptype)
		if err != nil {
			return err
		}
		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: Handshake,
			Plugins:         PluginMap,
			Cmd:             exec.Command(pluginPath),
			AllowedProtocols: []plugin.Protocol{
				plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
		})

		PluginClients[alias] = client
	}
	return nil
}

func GetAuthProviderClient(alias string) (v1.AuthenticationProviderClient, error) {
	if _, ok := PluginClients[alias]; !ok {
		return nil, fmt.Errorf("plugin for alias '%s' is not loaded", alias)
	}
	// Connect via RPC
	rpcClient, err := PluginClients[alias].Client()
	if err != nil {
		return nil, err
	}

	// Request the plugin
	raw, err := rpcClient.Dispense(PluginTypeAuthProvider)
	if err != nil {
		return nil, err
	}

	return raw.(v1.AuthenticationProviderClient), nil
}

func GetAuthVerifierClient(alias string) (v1.AuthenticationVerifierClient, error) {
	if _, ok := PluginClients[alias]; !ok {
		return nil, fmt.Errorf("plugin for alias '%s' is not loaded", alias)
	}
	// Connect via RPC
	rpcClient, err := PluginClients[alias].Client()
	if err != nil {
		return nil, err
	}

	// Request the plugin
	raw, err := rpcClient.Dispense(PluginTypeAuthVerifier)
	if err != nil {
		return nil, err
	}

	return raw.(v1.AuthenticationVerifierClient), nil
}

func GetConfigurationClient(alias string) (v1.ConfigurationClient, error) {
	if _, ok := PluginClients[alias]; !ok {
		return nil, fmt.Errorf("plugin for alias '%s' is not loaded", alias)
	}
	// Connect via RPC
	rpcClient, err := PluginClients[alias].Client()
	if err != nil {
		return nil, err
	}

	// Request the plugin
	raw, err := rpcClient.Dispense(PluginTypeConfiguration)
	if err != nil {
		return nil, err
	}

	return raw.(v1.ConfigurationClient), nil
}

func ResolvePluginTypeNameAndAlias(config string) (string, string, string, error) {
	aliasParts := strings.Split(config, ":")
	if len(aliasParts) > 2 {
		return "", "", "", fmt.Errorf("could not resolve plugin alias for '%s'", config)
	}
	alias := aliasParts[len(aliasParts)-1]
	typeParts := strings.Split(aliasParts[0], "_")
	if len(typeParts) != 2 {
		return "", "", "", fmt.Errorf("could not resolve plugin type or name for '%s'", config)
	}
	pluginType := typeParts[0]
	pluginName := typeParts[1]

	return pluginType, pluginName, alias, nil
}
