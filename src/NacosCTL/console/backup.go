package console

import (
	"fmt"
	mse "github.com/alibabacloud-go/mse-20190531/v3/client"
	"github.com/fatih/color"
	"github.com/longyuan/lib.v3/compress"
	"github.com/longyuan/lib.v3/ctl"
	"github.com/longyuan/nacos.v3/client"
	"github.com/longyuan/storage.v3/storage"
	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func Backup(host, username, password, outputFile string) (*string, error) {
	// 创建客户端
	nacosClient, err := client.NewNacosClient(host, username, password)
	if err != nil {
		return nil, err
	}

	// 备份目录
	backupDirectory, err := ctl.CreateTempDirectory(
		strings.ReplaceAll(strings.ReplaceAll(host, "/", "-"), "\\", "-"))
	if err != nil {
		return nil, err
	}

	// 打印命名空间
	namespaces, err := nacosClient.Namespaces()
	if err != nil {
		return nil, err
	}
	var namespaceTable [][]string
	for _, item := range namespaces.Data {
		namespaceTable = append(namespaceTable, []string{
			item.Namespace, item.NamespaceShowName, strconv.Itoa(item.ConfigCount),
		})
	}
	ctl.PrintTable([]string{"Nacos 命名空间ID", "Nacos 命名空间名称", "配置数量"}, namespaceTable)

	// 执行备份
	for _, it := range namespaces.Data {
		color.Blue(fmt.Sprintf("[Nacos] 备份命名空间 %s  /  %s", it.Namespace, it.NamespaceShowName))
		var namespaceId = it.Namespace
		if namespaceId == "" {
			continue
		}
		var items []client.NacosConfigItem
		items, err = nacosClient.NamespaceItems(namespaceId)
		if err != nil {
			return nil, err
		}
		// 创建目录
		var namespacePath = path.Join(*backupDirectory, namespaceId)
		if _, err = os.Stat(namespacePath); err != nil || os.IsNotExist(err) {
			err = os.MkdirAll(namespacePath, 644)
			if err != nil {
				return nil, err
			}
		}

		// Items
		for _, item := range items {
			var outputPath = path.Join(namespacePath, item.DataId+"."+item.Type)
			err = os.WriteFile(outputPath, []byte(item.Content), 644)
			if err != nil {
				return nil, err
			}
		}
	}

	// 备份完成
	colorPrint := color.New()
	colorPrint.Add(color.Bold)
	colorPrint.Add(color.FgGreen)
	_, err = colorPrint.Println("[Nacos] " + time.Now().Format("2006-01-02 15:04:05") + " 备份完成")
	if err != nil {
		return nil, err
	}
	// 压缩文件 zip
	if outputFile == "" {
		outputFile = "nacos.zip"
	}
	err = compress.Zip(*backupDirectory, outputFile, true)
	if err != nil {
		return nil, err
	}
	return &outputFile, nil
}

func AliBackup(accessKeyId, accessKeySecret, instanceId, namespace, outputFile string) (*string, error) {
	// 创建客户端
	nacosClient, err := client.NewAliyunNacosClient(accessKeyId, accessKeySecret, instanceId)
	if err != nil {
		return nil, err
	}

	// 备份目录
	backupDirectory, err := ctl.CreateTempDirectory(
		strings.ReplaceAll(strings.ReplaceAll(accessKeyId, "/", "-"), "\\", "-"))
	if err != nil {
		return nil, err
	}

	// 命名空间
	namespaces := strings.Split(namespace, ",")

	// 执行备份
	for _, it := range namespaces {
		color.Blue(fmt.Sprintf("[Nacos] 备份命名空间: %s", it))
		var namespaceId = it
		if namespaceId == "" {
			continue
		}
		var items []mse.ListNacosConfigsResponseBodyConfigurations
		items, err = nacosClient.GetNacosConfigList(namespaceId)
		if err != nil {
			return nil, err
		}
		// 创建目录
		var namespacePath = path.Join(*backupDirectory, namespaceId)
		if _, err = os.Stat(namespacePath); err != nil || os.IsNotExist(err) {
			err = os.MkdirAll(namespacePath, 644)
			if err != nil {
				return nil, err
			}
		}

		// Items
		for _, item := range items {
			itemDetail, err := nacosClient.GetNacosConfig(namespaceId, *item.Group, *item.DataId)
			log.Println(fmt.Sprintf("[Nacos 阿里云] %s / %s", namespaceId, *item.DataId))
			if err != nil {
				return nil, err
			}
			var outputPath = path.Join(namespacePath, *itemDetail.DataId+"."+*itemDetail.Type)
			err = os.WriteFile(outputPath, []byte(*itemDetail.Content), 644)
			if err != nil {
				return nil, err
			}
		}
	}

	// 备份完成
	colorPrint := color.New()
	colorPrint.Add(color.Bold)
	colorPrint.Add(color.FgGreen)
	_, err = colorPrint.Println("[Nacos] " + time.Now().Format("2006-01-02 15:04:05") + " 备份完成")
	if err != nil {
		return nil, err
	}
	// 压缩文件 zip
	if outputFile == "" {
		outputFile = "nacos.zip"
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
		tempDirectory, err := ctl.CreateTempDirectory("cron_nacos", dateTimeFormat)
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
			var nacosConfig = struct {
				Host     string `yaml:"host"`
				Username string `yaml:"username"`
				Password string `yaml:"password"`

				AccessKeyId     string `yaml:"accessKeyId"`
				AccessKeySecret string `yaml:"accessKeySecret"`
				InstanceId      string `yaml:"instanceId"`
				Namespace       string `yaml:"namespace"`
			}{}
			err = yaml.Unmarshal(fileBytes, &nacosConfig)
			if err != nil {
				return err
			}
			// 备份任务
			var outFileName = path.Base(fileName) + "_" + dateTimeFormat + ".zip"
			var outputFile = path.Join(*tempDirectory, outFileName)
			var backupZipFile *string
			if nacosConfig.InstanceId == "" {
				backupZipFile, err = Backup(nacosConfig.Host, nacosConfig.Username, nacosConfig.Password, outputFile)
				if err != nil {
					return err
				}
			} else {
				backupZipFile, err = AliBackup(nacosConfig.AccessKeyId, nacosConfig.AccessKeySecret, nacosConfig.InstanceId, nacosConfig.Namespace, outputFile)
				if err != nil {
					return err
				}
			}
			_, err = cosClient.Put(*backupZipFile, "nacos/"+dateFormat+"/"+outFileName)
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
