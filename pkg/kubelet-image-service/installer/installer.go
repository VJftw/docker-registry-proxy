package installer

import (
	"fmt"
	"os"
)

// Install installs the SRC_BIN env var into DEST_BIN and configures SystemD as well as kubelet to use it.
func Install() error {

	srcBinPath := os.Getenv("SRC_BIN")
	destBinPath := os.Getenv("DEST_BIN")

	if err := Copy(srcBinPath, destBinPath); err != nil {
		return fmt.Errorf("could not copy binary: %w", err)
	}

	if err := SystemD(destBinPath); err != nil {
		return fmt.Errorf("could not configure SystemD: %w", err)
	}

	// do a health check on kubelet-image-service before reconfiguring and restarting kubelet

	if err := ReconfigureKubelet(); err != nil {
		return fmt.Errorf("could not reconfigure kubelet: %w", err)
	}

	if err := SystemDRestartKubelet(); err != nil {
		RollbackKubeletConfiguration()
	}

	// do a health check on kubelet. If health check fails, roll back.

	return nil
}
