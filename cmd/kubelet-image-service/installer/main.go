package main

import (
	"github.com/VJftw/docker-registry-proxy/pkg/cmd"
	"github.com/VJftw/docker-registry-proxy/pkg/installer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	viper.SetEnvPrefix("kisi")
	rootCmd := cmd.New("kubelet-image-service_installer")

	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {

		// 1. copy binary
		if err := installer.Copy(
			"/bin/kubelet-image-service",
			"/opt/bin/kubelet-image-service",
		); err != nil {
			return err
		}

		// 2. configure systemd
		if err := installer.SystemD(
			"/configs/kubelet-image-service/installer/systemd.tpl.service",
			"/opt/bin/kubelet-image-service",
		); err != nil {
			return err
		}

		// 3. Reconfigure kubelet
		if err := installer.ReconfigureKubelet(); err != nil {

		}

		return nil
	}

	cmd.HandleErr(rootCmd.Execute())
}
