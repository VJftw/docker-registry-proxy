package main

import (
	"fmt"
	"strings"

	"github.com/VJftw/docker-registry-proxy/pkg/cmd"
	v1 "github.com/VJftw/docker-registry-proxy/api/proto/v1"
	"github.com/VJftw/docker-registry-proxy/pkg/plugin"
	"github.com/VJftw/docker-registry-proxy/pkg/runtimes/docker"
	"github.com/spf13/viper"

	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

const (
	flagAuthenticationProviders = "authentication_provider"
)

func main() {
	viper.SetEnvPrefix("kis")
	rootCmd := cmd.New("kubelet-image-service")

	rootCmd.Flags().StringSlice(flagAuthenticationProviders, []string{}, "The authentication providers in the format `<image_prefix>=<endpoint>`")
	if err := viper.BindPFlag(flagAuthenticationProviders, rootCmd.Flags().Lookup(flagAuthenticationProviders)); err != nil {
		cmd.HandleErr(err)
	}
	grpcServer := cmd.NewGRPCServer()

	preFunc := func() error {
		authProviders, err := parseAuthenticationProviderClients()
		if err != nil {
			return fmt.Errorf("could not parse authentication providers: %w", err)
		}
		imageService, err := docker.NewImageService(authProviders)
		if err != nil {
			return fmt.Errorf("could not create docker image service: %w", err)
		}
		runtimeapi.RegisterImageServiceServer(grpcServer, imageService)
		return nil
	}

	cmd.Execute(rootCmd, grpcServer, preFunc)
}

func parseAuthenticationProviderClients() (map[string]v1.AuthenticationProviderClient, error) {
	res := map[string]v1.AuthenticationProviderClient{}
	confs := viper.GetStringSlice(flagAuthenticationProviders)
	for _, conf := range confs {
		parts := strings.Split(conf, "=")
		registryPrefix := parts[0]
		alias := parts[1]
		client, err := plugin.GetAuthProviderClient(alias)
		if err != nil {
			return nil, err
		}
		res[registryPrefix] = client
	}
	return res, nil
}
