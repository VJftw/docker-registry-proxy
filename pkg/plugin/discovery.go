package plugin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func FindPluginInPath(path string, rName string, rType string) (string, error) {
	pluginBinary := ""
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		// ./plugin_<type>_<name>
		// e.g. ./plugin_auth-provider_static
		if !f.IsDir() && strings.HasPrefix(f.Name(), "plugin_") {
			parts := strings.Split(f.Name(), "_")
			pType := parts[1]
			name := parts[2]
			if name == rName && pType == rType {
				pluginBinary = path
				return nil
			}
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if pluginBinary == "" {
		return "", fmt.Errorf("could not find binary: plugin_%s_%s", rType, rName)
	}
	return pluginBinary, nil
}
