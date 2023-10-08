package console

import (
	"errors"
	"fmt"
	"github.com/longyuan/docker.v3/client"
	"github.com/longyuan/lib.v3/ctl"
	"os"
	"path"
)

type S3Config struct {
	AccessKeyId     string
	SecretAccessKey string
	RootPath        string
}

type DockerConfig struct {
	Address string
}

type Read interface {
	types() string
}
type Write interface {
	types() string
}

type S3Read struct {
	Value *S3Config
}

func (d *S3Read) types() string {
	return "s3"
}

type DockerRead struct {
	Docker *DockerConfig
}

func (d *DockerRead) types() string {
	return "docker"
}

type S3Write struct {
	Value *S3Config
}

func (d *S3Write) types() string {
	return "s3"
}

type DockerWrite struct {
	Docker *DockerConfig
}

func (d *DockerWrite) types() string {
	return "docker"
}

func Pipeline(input Read, write Write) error {
	var inputType = input.types()
	var writeType = write.types()
	handle, ok := instance()[inputType+"_"+writeType]
	if ok {
		err := handle(input, write)
		if err != nil {
			return err
		}
	}
	return nil
}

func instance() map[string]func(input Read, write Write) error {
	return map[string]func(input Read, write Write) error{
		"s3_docker": s3ToDocker,
		"docker_s3": dockerToS3,
	}
}

func s3ToDocker(input Read, write Write) error {
	return nil
}

func dockerToS3(input Read, write Write) error {
	inputDocker, ok := input.(*DockerRead)
	if !ok {
		return errors.New("input must be a DockerRead")
	}
	writeS3, ok := write.(*S3Write)
	if !ok {
		return errors.New("write must be a DockerRead")
	}

	// Docker Client
	dockerClient, err := client.NewDockerClient(inputDocker.Docker.Address, "", "")
	if err != nil {
		return err
	}
	imageList, err := dockerClient.ImageList()
	if err != nil {
		return err
	}

	// 备份目录
	backupDirectory, err := ctl.CreateTempDirectory("docker")
	for _, item := range imageList {
		// 导出镜像
		var outputFile = path.Join(*backupDirectory, item.ID+".docker.cache")
		var outputCacheFile = outputFile + ".cache"
		err = dockerClient.ImageExport(item.ID, outputCacheFile)
		if err != nil {
			return err
		}

		// 防止文件损坏, 重命名文件
		err = os.Rename(outputCacheFile, outputFile)
		if err != nil {
			return err
		}

		// 推送到S3

		// 删除临时目录
		err = os.Remove(outputFile)
		if err != nil {
			return err
		}
	}
	fmt.Println(dockerClient)
	fmt.Println(writeS3)
	return nil
}
