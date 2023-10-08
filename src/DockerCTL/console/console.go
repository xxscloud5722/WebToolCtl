package console

import "github.com/longyuan/docker.v3/client"

func Push(tcp, username, password, image string) error {
	// 创建
	dockerClient, err := client.NewDockerClient(tcp, username, password)
	if err != nil {
		return err
	}

	// 推送到仓库
	err = dockerClient.DockerPush(image)
	if err != nil {
		return err
	}

	return nil
}

func Remove(tcp, image string) error {
	// 创建
	dockerClient, err := client.NewDockerClient(tcp, "", "")
	if err != nil {
		return err
	}

	// 删除镜像
	err = dockerClient.DockerRemove(image)
	if err != nil {
		return err
	}

	return nil
}
