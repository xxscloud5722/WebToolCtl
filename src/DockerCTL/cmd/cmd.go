package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/longyuan/docker.v3/console"
	"github.com/spf13/cobra"
)

func Pipeline() []*cobra.Command {
	var pipeCmd = &cobra.Command{
		Use:     "pipe",
		Short:   "Docker Image Pipeline",
		Example: "pipe -iS3.ak [AK0001] -iS3.sk [SK0001] -oS3.ak [AK0001] -oS3.sk [SK0001] -iDocker.tcp -oDocker.tcp",
		Run: func(cmd *cobra.Command, args []string) {
			inputAccessKeyId, err := cmd.Flags().GetString("iS3.ak")
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
			inputSecretKeyId, err := cmd.Flags().GetString("iS3.sk")
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
			inputDockerTcp, err := cmd.Flags().GetString("iDocker.tcp")
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}

			writeAccessKeyId, err := cmd.Flags().GetString("oS3.ak")
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
			writeSecretKeyId, err := cmd.Flags().GetString("oS3.sk")
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
			writeDockerTcp, err := cmd.Flags().GetString("oDocker.tcp")
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}

			var read console.Read
			if inputAccessKeyId != "" && inputSecretKeyId != "" {
				read = &console.S3Read{
					Value: &console.S3Config{
						AccessKeyId:     inputAccessKeyId,
						SecretAccessKey: inputSecretKeyId,
						RootPath:        "",
					},
				}
			} else {
				read = &console.DockerRead{Docker: &console.DockerConfig{Address: inputDockerTcp}}
			}

			var write console.Write
			if writeAccessKeyId != "" && writeSecretKeyId != "" {
				write = &console.S3Write{
					Value: &console.S3Config{
						AccessKeyId:     writeAccessKeyId,
						SecretAccessKey: writeSecretKeyId,
						RootPath:        "",
					},
				}
			} else {
				write = &console.DockerWrite{Docker: &console.DockerConfig{Address: writeDockerTcp}}
			}
			err = console.Pipeline(read, write)
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
		},
	}
	pipeCmd.Flags().String("iS3.ak", "", "")
	pipeCmd.Flags().String("iS3.sk", "", "")
	pipeCmd.Flags().String("iDocker.tcp", "", "")

	pipeCmd.Flags().String("oS3.ak", "", "")
	pipeCmd.Flags().String("oS3.sk", "", "")
	pipeCmd.Flags().String("oDocker.tcp", "", "")

	return []*cobra.Command{
		pipeCmd,
	}
}

func Console() []*cobra.Command {
	var pushCmd = &cobra.Command{
		Use:     "push",
		Short:   "Docker Image Push (Auth)",
		Example: "push -u docker -p 12345678 -i nginx:1.21",
		Run: func(cmd *cobra.Command, args []string) {
			username, err := cmd.Flags().GetString("username")
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
			password, err := cmd.Flags().GetString("password")
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
			tcp, err := cmd.Flags().GetString("tcp")
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
			image, err := cmd.Flags().GetString("image")
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
			if image == "" {
				color.Red(fmt.Sprint("image is null ?"))
				return
			}

			err = console.Push(tcp, username, password, image)
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
		},
	}
	pushCmd.Flags().StringP("username", "u", "", "")
	pushCmd.Flags().StringP("password", "p", "", "")
	pushCmd.Flags().StringP("tcp", "t", "", "")
	pushCmd.Flags().StringP("image", "i", "", "")

	var removeCmd = &cobra.Command{
		Use:     "remove",
		Short:   "Docker Image Remove (Force)",
		Example: "remove -i nginx:1.21",
		Run: func(cmd *cobra.Command, args []string) {
			tcp, err := cmd.Flags().GetString("tcp")
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
			image, err := cmd.Flags().GetString("image")
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
			if image == "" {
				color.Red(fmt.Sprint("image is null ?"))
				return
			}

			err = console.Remove(tcp, image)
			if err != nil {
				color.Red(fmt.Sprint(err))
				return
			}
		},
	}
	removeCmd.Flags().StringP("tcp", "t", "", "")
	removeCmd.Flags().StringP("image", "i", "", "")

	return []*cobra.Command{
		pushCmd,
		removeCmd,
	}
}
