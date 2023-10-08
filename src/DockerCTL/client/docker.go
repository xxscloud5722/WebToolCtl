package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"io"
	"os"
)

const dockerUnix = "unix:///var/run/docker.sock"

type DockerClient struct {
	dockerClient *client.Client // Docker 镜像
	auth         string         // 授权信息
}

// NewDockerClient 创建Docker 客户端
func NewDockerClient(address, username, password string) (*DockerClient, error) {
	if address == "" {
		address = dockerUnix
	}
	dockerClient, err := client.NewClientWithOpts(client.WithHost(address),
		client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	if username != "" && password != "" {
		authJson, err := json.Marshal(registry.AuthConfig{
			Username: username,
			Password: password,
		})
		if err != nil {
			return nil, err
		}
		return &DockerClient{dockerClient: dockerClient, auth: base64.URLEncoding.EncodeToString(authJson)}, nil
	} else {
		return &DockerClient{dockerClient: dockerClient}, nil
	}
}

// DockerRemove 删除容器镜像
func (engine *DockerClient) DockerRemove(imageName string) error {
	imageList, err := engine.dockerClient.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}
	var imageId = func() *string {
		for _, image := range imageList {
			for _, it := range image.RepoTags {
				if it == imageName {
					return &image.ID
				}
			}
		}
		return nil
	}()
	if imageId == nil {
		return nil
	}
	_, err = engine.dockerClient.ImageRemove(context.Background(), *imageId, types.ImageRemoveOptions{
		Force: true,
	})
	if err != nil {
		return err
	}
	return nil
}

// DockerTag 删除容器重命名
func (engine *DockerClient) DockerTag(sourceImage, targetImage string) error {
	err := engine.dockerClient.ImageTag(context.Background(), sourceImage, targetImage)
	if err != nil {
		return err
	}
	return nil
}

// DockerPull 远程仓库拉取最新镜像
func (engine *DockerClient) DockerPull(image string) error {
	reader, err := engine.dockerClient.ImagePull(context.Background(), image, types.ImagePullOptions{
		RegistryAuth: engine.auth,
	})
	if err != nil {
		return err
	}
	defer func(reader io.ReadCloser) {
		err := reader.Close()
		if err != nil {
			panic(err)
		}
	}(reader)
	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		return err
	}
	return nil
}

// DockerPush 推送镜像到远程仓库
func (engine *DockerClient) DockerPush(image string) error {
	reader, err := engine.dockerClient.ImagePush(context.Background(), image, types.ImagePushOptions{
		RegistryAuth: engine.auth,
	})
	if err != nil {
		return err
	}
	defer func(reader io.ReadCloser) {
		err := reader.Close()
		if err != nil {
			panic(err)
		}
	}(reader)
	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		return err
	}
	return nil
}

// ImageList 镜像列表
func (engine *DockerClient) ImageList() ([]types.ImageSummary, error) {
	result, err := engine.dockerClient.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ImageExport 镜像导出
func (engine *DockerClient) ImageExport(imageId, outFile string) error {
	readCloser, err := engine.dockerClient.ImageSave(context.Background(), []string{imageId})
	if err != nil {
		return err
	}
	defer func(readCloser io.ReadCloser) {
		err := readCloser.Close()
		if err != nil {
			panic(err)
		}
	}(readCloser)
	outImageFile, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer func(outFile *os.File) {
		err := outFile.Close()
		if err != nil {
			panic(err)
		}
	}(outImageFile)
	_, err = io.Copy(outImageFile, readCloser)
	if err != nil {
		return err
	}
	return nil
}

// ImageImport 镜像导入
func (engine *DockerClient) ImageImport(imageRef, importPath string) error {
	imageFile, err := os.Open(importPath)
	if err != nil {
		return nil
	}
	defer func(imageFile *os.File) {
		err := imageFile.Close()
		if err != nil {
			panic(err)
		}
	}(imageFile)
	reader, err := engine.dockerClient.ImageImport(context.Background(), types.ImageImportSource{
		Source: imageFile,
	}, imageRef, types.ImageImportOptions{})
	if err != nil {
		return err
	}
	defer func(reader io.ReadCloser) {
		err := reader.Close()
		if err != nil {
			panic(err)
		}
	}(reader)

	// 打印输出
	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		return err
	}

	return nil
}
