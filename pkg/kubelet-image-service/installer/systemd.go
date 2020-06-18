package installer

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/coreos/go-systemd/v22/unit"
)

// SystemDUnitDir is where to place the systemD configuration
const SystemDUnitDir = "/etc/systemd/system/"

// SystemD configures SystemD and starts the new Kubelet Service
func SystemD(binaryPath string) error {
	// write systemd unit file
	systemdOpts := []*unit.UnitOption{
		{Section: "Unit", Name: "Description", Value: "Kubernetes Image Service"},
		{Section: "Service", Name: "ExecStart", Value: binaryPath},
		{Section: "Service", Name: "Restart", Value: "always"},
		{Section: "Service", Name: "StartLimitInterval", Value: "0"},
		{Section: "Service", Name: "RestartSec", Value: "3s"},
		{Section: "Install", Name: "WantedBy", Value: "multi-user.target"},
	}
	serializedReader := unit.Serialize(systemdOpts)
	serializedBytes, err := ioutil.ReadAll(serializedReader)
	if err != nil {
		return fmt.Errorf("could not generate SystemD Unit file: %w", err)
	}
	if err := ioutil.WriteFile(
		filepath.Join(SystemDUnitDir, "kubelet-image-service.service"),
		serializedBytes,
		0755,
	); err != nil {
		return fmt.Errorf("could not write SystemD Unit file: %w", err)
	}

	// systemctl daemon-reload
	// https://godoc.org/github.com/coreos/go-systemd/dbus#Conn.Reload
	systemDConn, err := dbus.NewSystemdConnection()
	if err != nil {
		return fmt.Errorf("could not get systemd connection: %w", err)
	}
	defer systemDConn.Close()
	if err := systemDConn.Reload(); err != nil {
		return fmt.Errorf("unable to reload system configuration")
	}

	// systemctl restart kubelet-image-service
	// https://godoc.org/github.com/coreos/go-systemd/dbus#Conn.RestartUnit
	resCh := make(chan string)
	systemDConn.RestartUnit("kubelet-image-service.service", "replace", resCh)
	job := <-resCh
	if job != "done" {
		return fmt.Errorf("unable to restart kubelet-image-service as job != 'done': %s", job)
	}

	return nil
}
