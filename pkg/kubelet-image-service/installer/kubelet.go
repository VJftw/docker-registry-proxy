package installer

import (
	"fmt"

	"github.com/coreos/go-systemd/v22/dbus"
)

// ReconfigureKubelet attempts to reconfigure kubelet on the system
func ReconfigureKubelet() error {
	return nil
}

// ReconfigureKubeletGKE reconfigures kubelet in GKE managed nodes
func ReconfigureKubeletGKE() error {
	return nil
}

// ReconfigureKubeletEKS reconfigures kubelet in EKS managed nodes
func ReconfigureKubeletEKS() error {
	return nil
}

// SystemDRestartKubelet restarts Kubelet
func SystemDRestartKubelet() error {
	systemDConn, err := dbus.NewSystemdConnection()
	if err != nil {
		return fmt.Errorf("could not get systemd connection: %w", err)
	}
	defer systemDConn.Close()
	resCh := make(chan string)
	systemDConn.RestartUnit("kubelet.service", "replace", resCh)
	job := <-resCh
	if job != "done" {
		return fmt.Errorf("unable to restart kubelet as job != 'done': %s", job)
	}

	return nil
}

// RollbackKubeletConfiguration rolls back kubelet configuration on the system
func RollbackKubeletConfiguration() error {
	return nil
}

// RollbackKubeletGKE reconfigures kubelet in GKE managed nodes
func RollbackKubeletGKE() error {
	return nil
}

// RollbackKubeletEKS reconfigures kubelet in GKE managed nodes
func RollbackKubeletEKS() error {
	return nil
}
