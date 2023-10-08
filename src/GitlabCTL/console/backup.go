package console

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/longyuan/gitlab.v3/client"
	"github.com/longyuan/lib.v3/compress"
	"github.com/longyuan/lib.v3/ctl"
	"github.com/longyuan/storage.v3/storage"
	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func Backup(host, token, outputFile string) (*string, error) {
	gitlabClient, err := client.NewGitlabClient(host, token)
	if err != nil {
		return nil, err
	}
	color.Cyan("[Gitlab] 扫描项目列表 ...")
	projects, err := gitlabClient.Projects()
	if err != nil {
		return nil, err
	}
	color.Blue("[Gitlab] 扫描完成，当前授权可访问项目共: %d", len(projects))

	colorPrint := color.New()
	colorPrint.Add(color.Bold)
	colorPrint.Add(color.FgYellow)
	_, err = colorPrint.Println("备份时间会很长，请耐心等待（不要关闭正在执行的程序）")
	if err != nil {
		return nil, err
	}

	// 临时目录
	backupDirectory, err := ctl.CreateTempDirectory("gitlab")
	if err != nil {
		return nil, err
	}
	for _, project := range projects {
		for {
			// 导出配置
			jsonFile, err := json.Marshal(&project)
			if err != nil {
				return nil, err
			}
			color.Blue(fmt.Sprintf("[Gitlab] 导出项目: %s", project.Name))
			var projectConfigPath = path.Join(*backupDirectory, fmt.Sprintf("project.%d.json", project.ID))
			if _, err := os.Stat(projectConfigPath); err == nil || !os.IsNotExist(err) {
				err := os.Remove(projectConfigPath)
				if err != nil {
					return nil, err
				}
			}
			err = os.WriteFile(projectConfigPath, jsonFile, 644)
			if err != nil {
				return nil, err
			}
			// 导出项目
			var projectId = project.ID
			// 判断是否已经下载过
			var projectOutputPath = path.Join(*backupDirectory, strconv.Itoa(projectId)+".gitlab")
			if _, err = os.Stat(projectOutputPath); err == nil || !os.IsNotExist(err) {
				break
			}
			err = gitlabClient.Export(projectId, projectOutputPath)
			if err != nil {
				color.Red("[Gitlab] 导出时发生异常 (等待3s): " + fmt.Sprint(err))
				time.Sleep(3 * time.Second)
				continue
			}
			break
		}
	}

	// 压缩文件 zip
	if outputFile == "" {
		outputFile = "gitlab.zip"
	}
	err = compress.Zip(*backupDirectory, outputFile, true)
	if err != nil {
		return nil, err
	}
	return &outputFile, nil
}

func CronBackup(configPath, backupCron, cloudStorageConfig string) error {
	var values = strings.Split(cloudStorageConfig, ",")
	if len(values) < 3 {
		return fmt.Errorf("CloudStorageConfig error")
	}
	var url = values[0]
	var secretId = values[1]
	var secretKey = values[2]
	cosClient, err := storage.NewTencentCOS(url, secretId, secretKey)
	if err != nil {
		return err
	}

	c := cron.New()
	_, err = c.AddFunc(backupCron, func() {
		var dateFormat = time.Now().Format("2006_01_02")
		var dateTimeFormat = time.Now().Format("2006_01_02_15_04_05")
		tempDirectory, err := ctl.CreateTempDirectory("cron_gitlab", dateTimeFormat)
		if err != nil {
			color.Red(fmt.Sprint(err))
			return
		}
		err = filepath.Walk(configPath, func(configItemPath string, fi os.FileInfo, errBack error) (err error) {
			var fileName = fi.Name()
			if !strings.HasSuffix(fileName, ".yaml") {
				return nil
			}
			// 读取Yaml 文件
			fileBytes, err := os.ReadFile(configItemPath)
			if err != nil {
				return err
			}
			var gitlabConfig = struct {
				Host  string `yaml:"host"`
				Token string `yaml:"token"`
			}{}
			err = yaml.Unmarshal(fileBytes, &gitlabConfig)
			if err != nil {
				return err
			}
			// 备份
			var outFileName = path.Base(fileName) + "_" + dateTimeFormat + ".zip"
			var outputFile = path.Join(*tempDirectory, outFileName)
			backupZipFile, err := Backup(gitlabConfig.Host, gitlabConfig.Token, outputFile)
			if err != nil {
				return err
			}
			_, err = cosClient.Put(*backupZipFile, "gitlab/"+dateFormat+"/"+outFileName)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			color.Red(fmt.Sprint(err))
		}
	})
	if err != nil {
		return err
	}
	color.Blue(fmt.Sprintf("Cron (%s) Start Success ...", backupCron))
	color.Blue(fmt.Sprintf("ConfigPath: %s", configPath))
	color.Blue(fmt.Sprintf("S3 Config: %s", cloudStorageConfig))
	c.Start()
	select {}
}
