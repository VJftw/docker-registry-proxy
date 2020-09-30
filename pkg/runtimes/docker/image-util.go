package docker

import (
	"fmt"
	"log"
	"strings"

	dockerregistryproxyv1 "github.com/VJftw/docker-registry-proxy/api/proto/v1"
	"github.com/docker/distribution/reference"
)

// GetRepository returns the given image's repository
func GetRepository(imageRef string) (string, error) {
	matches := reference.NameRegexp.FindStringSubmatch(imageRef)

	if len(matches) > 0 {
		if len(strings.Split(matches[0], "/")) == 1 {
			return fmt.Sprintf("index.docker.io/%s", matches[0]), nil
		}
		return matches[0], nil
	}

	return "", fmt.Errorf("invalid repository")
}

// ResolveRepositoryAuthProvider resolves the authentication provider to use for the given image
func ResolveRepositoryAuthProvider(repository string, authProviders map[string]dockerregistryproxyv1.AuthenticationProviderClient) (dockerregistryproxyv1.AuthenticationProviderClient, error) {

	repoParts := strings.Split(repository, "/")
	for i := len(repoParts); i > 0; i-- {
		repoAuth := strings.Join(repoParts[:i], "/")
		log.Println(repository, repoAuth)
		if val, ok := authProviders[repoAuth]; ok {
			return val, nil
		}
	}

	return nil, fmt.Errorf("could not resolve auth provider for repository: %s", repository)
}
