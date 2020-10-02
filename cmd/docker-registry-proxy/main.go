package main

import (
	"net/http"
	"strings"

	dockerregistryproxyv1 "github.com/VJftw/docker-registry-proxy/api/proto/v1"
	"github.com/VJftw/docker-registry-proxy/pkg/cmd"
	"github.com/VJftw/docker-registry-proxy/pkg/plugin"
	"github.com/VJftw/docker-registry-proxy/pkg/runtimes/docker/registryproxy"
	"github.com/spf13/viper"
)

const (
	flagAuthenticationVerifiers       = "authentication_verifier"
	flagUpstreamAuthenticationService = "upstream_authentication"
	flagUpstreamRepository            = "upstream_repository"
)

func main() {
	viper.SetEnvPrefix("drp")
	rootCmd := cmd.New("docker-registry-proxy")

	rootCmd.Flags().String(flagUpstreamRepository, "", "The upstream repository")
	if err := viper.BindPFlag(flagUpstreamRepository, rootCmd.Flags().Lookup(flagUpstreamRepository)); err != nil {
		cmd.HandleErr(err)
	}

	rootCmd.Flags().String(flagUpstreamAuthenticationService, "", "The upstream authentication service that returns credentials to pull images from the source repository")
	if err := viper.BindPFlag(flagUpstreamAuthenticationService, rootCmd.Flags().Lookup(flagUpstreamAuthenticationService)); err != nil {
		cmd.HandleErr(err)
	}
	// --upstream_authentication="static"

	rootCmd.Flags().StringSlice(flagAuthenticationVerifiers, []string{}, "The authentication verifiers in the format `<username>=<endpoint>`")
	if err := viper.BindPFlag(flagAuthenticationVerifiers, rootCmd.Flags().Lookup(flagAuthenticationVerifiers)); err != nil {
		cmd.HandleErr(err)
	}
	// --authentication_verifier="_gcp:google"
	// --authentication_verifier="_aws:aws"

	httpServer := cmd.NewHTTPServer()
	preFunc := func() error {
		var upstreamAuthService dockerregistryproxyv1.AuthenticationProviderAPIClient
		if upstreamAuth := viper.GetString(flagUpstreamAuthenticationService); upstreamAuth != "" {
			uAS, err := plugin.GetAuthProviderClient(upstreamAuth)
			if err != nil {
				return err
			}
			upstreamAuthService = uAS
		}
		authVerifiers, err := parseAuthenticationVerifierClients()
		if err != nil {
			return err
		}
		proxyOpts, err := registryproxy.GetProxyOpts(
			viper.GetString(flagUpstreamRepository),
			upstreamAuthService,
			authVerifiers,
		)
		if err != nil {
			return err
		}
		http.HandleFunc("/", registryproxy.ProxyHandler(proxyOpts))

		return nil
	}

	cmd.Execute(rootCmd, httpServer, preFunc)
}

func parseAuthenticationVerifierClients() (map[string]dockerregistryproxyv1.AuthenticationVerifierAPIClient, error) {
	res := map[string]dockerregistryproxyv1.AuthenticationVerifierAPIClient{}
	confs := viper.GetStringSlice(flagAuthenticationVerifiers)
	for _, conf := range confs {
		parts := strings.Split(conf, ":")
		username := parts[0]
		alias := parts[1]
		client, err := plugin.GetAuthVerifierClient(alias)
		if err != nil {
			return nil, err
		}

		res[username] = client
	}

	return res, nil
}
