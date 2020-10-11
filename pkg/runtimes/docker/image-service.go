package docker

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	dockerregistryproxyv1 "github.com/VJftw/docker-registry-proxy/api/proto/v1"
	"github.com/docker/docker/api/types"
	"github.com/google/cadvisor/fs"

	dockerclient "github.com/docker/docker/client"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// ImageService represents a custom Docker image service for K8s
type ImageService struct {
	client *dockerclient.Client

	authProviders map[string]dockerregistryproxyv1.AuthenticationProviderAPIClient
}

// NewImageService returns a new instance of the Docker Image Service with the given authentication providers
func NewImageService(authProviders map[string]dockerregistryproxyv1.AuthenticationProviderAPIClient) (*ImageService, error) {
	client, err := dockerclient.NewClientWithOpts()
	if err != nil {
		return nil, fmt.Errorf("could not create new docker client: %w", err)
	}
	if err := dockerclient.FromEnv(client); err != nil {
		return nil, fmt.Errorf("could not configure docker client from environment: %w", err)
	}
	return &ImageService{
		client:        client,
		authProviders: authProviders,
	}, nil
}

// ListImages returns a List of Images
// TODO: filter images based on request
func (s *ImageService) ListImages(ctx context.Context, req *runtimeapi.ListImagesRequest) (*runtimeapi.ListImagesResponse, error) {
	images, err := s.client.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return nil, err
	}
	respImages := []*runtimeapi.Image{}
	for _, image := range images {
		respImages = append(respImages, &runtimeapi.Image{
			Id:          image.ID,
			RepoTags:    image.RepoTags,
			RepoDigests: image.RepoDigests,
			Size_:       uint64(image.Size),
		})
	}
	return &runtimeapi.ListImagesResponse{
		Images: respImages,
	}, nil
}

// ImageStatus returns the status of a given image
func (s *ImageService) ImageStatus(ctx context.Context, req *runtimeapi.ImageStatusRequest) (*runtimeapi.ImageStatusResponse, error) {
	image, _, err := s.client.ImageInspectWithRaw(ctx, req.GetImage().GetImage())
	resp := &runtimeapi.ImageStatusResponse{}
	if err != nil {
		return resp, nil
	}
	resp.Image = &runtimeapi.Image{
		Id:          image.ID,
		RepoTags:    image.RepoTags,
		RepoDigests: image.RepoDigests,
		Size_:       uint64(image.Size),
	}
	return resp, nil

}

// PullImage pulls a Docker Image
func (s *ImageService) PullImage(ctx context.Context, req *runtimeapi.PullImageRequest) (*runtimeapi.PullImageResponse, error) {
	img := req.GetImage().GetImage()
	imgPullOpts := types.ImagePullOptions{All: false}

	repository, err := GetRepository(img)
	if err != nil {
		return nil, err
	}
	authProvider, err := ResolveRepositoryAuthProvider(repository, s.authProviders)
	if err != nil {
		log.Println(err) // Only warn for images that may not need authentication
	}
	if authProvider != nil {
		authResp, err := authProvider.Provide(ctx, &dockerregistryproxyv1.ProvideRequest{})
		if err != nil {
			return nil, err
		}
		regAuth := base64.StdEncoding.EncodeToString(
			[]byte(fmt.Sprintf(
				`{"username": "%s", "password": "%s"}`,
				authResp.GetUsername(),
				authResp.GetPassword(),
			)),
		)
		imgPullOpts.RegistryAuth = regAuth
	}

	out, err := s.client.ImagePull(ctx, img, imgPullOpts)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		os.Stdout.Write(scanner.Bytes())
	}
	out.Close()

	image, _, err := s.client.ImageInspectWithRaw(ctx, req.GetImage().GetImage())
	if err != nil {
		return nil, err
	}

	return &runtimeapi.PullImageResponse{
		ImageRef: image.ID,
	}, nil
}

// RemoveImage Removes the given Docker image
func (s *ImageService) RemoveImage(ctx context.Context, req *runtimeapi.RemoveImageRequest) (*runtimeapi.RemoveImageResponse, error) {
	_, err := s.client.ImageRemove(ctx, req.GetImage().GetImage(), types.ImageRemoveOptions{PruneChildren: true})
	if err != nil {
		return nil, err
	}
	return &runtimeapi.RemoveImageResponse{}, nil
}

// ImageFsInfo returns status about the filesystem that images are stored on
func (s *ImageService) ImageFsInfo(ctx context.Context, _ *runtimeapi.ImageFsInfoRequest) (*runtimeapi.ImageFsInfoResponse, error) {
	info, err := s.client.Info(ctx)
	if err != nil {
		return nil, err
	}
	usageInfo, err := fs.GetDirUsage(info.DockerRootDir)
	if err != nil {
		return nil, err
	}
	return &runtimeapi.ImageFsInfoResponse{
		ImageFilesystems: []*runtimeapi.FilesystemUsage{
			{
				Timestamp:  time.Now().UnixNano(),
				FsId:       &runtimeapi.FilesystemIdentifier{Mountpoint: info.DockerRootDir},
				UsedBytes:  &runtimeapi.UInt64Value{Value: usageInfo.Bytes},
				InodesUsed: &runtimeapi.UInt64Value{Value: usageInfo.Inodes},
			},
		},
	}, nil
}
