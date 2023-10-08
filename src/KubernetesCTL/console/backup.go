package console

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/longyuan/kubernetes.v3/client"
	"github.com/longyuan/lib.v3/compress"
	"github.com/longyuan/lib.v3/ctl"
	"github.com/longyuan/storage.v3/storage"
	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v3"
	v1 "k8s.io/api/core/v1"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type BackupClient struct {
	rootPath string
	client   *client.KClient
}

// Backup 备份
func Backup(configPath, outputFile string, filter func(namespace v1.Namespace) bool) (*string, error) {
	// 创建客户端
	kClient, err := client.NewKClient(configPath)
	if err != nil {
		return nil, err
	}

	// 备份目录
	backupDirectory, err := ctl.CreateTempDirectory(kClient.Name)
	if err != nil {
		return nil, err
	}

	// 读取命名空间备份
	namespaces, err := kClient.Namespaces()
	if err != nil {
		return nil, err
	}
	var backupClient = &BackupClient{
		client:   kClient,
		rootPath: *backupDirectory,
	}
	for _, namespace := range namespaces {
		if filter != nil && !filter(namespace) {
			continue
		}
		var namespaceName = namespace.ObjectMeta.Name
		err = backupClient.backupNamespace(namespaceName)
		if err != nil {
			return nil, err
		}
		err = backupClient.backupDeployment(namespaceName)
		if err != nil {
			return nil, err
		}
		err = backupClient.backupStatefulSet(namespaceName)
		if err != nil {
			return nil, err
		}
		err = backupClient.backupDaemonSet(namespaceName)
		if err != nil {
			return nil, err
		}
		err = backupClient.backupJob(namespaceName)
		if err != nil {
			return nil, err
		}
		err = backupClient.backupCronJob(namespaceName)
		if err != nil {
			return nil, err
		}
		err = backupClient.backupConfigmap(namespaceName)
		if err != nil {
			return nil, err
		}
		err = backupClient.backupSecret(namespaceName)
		if err != nil {
			return nil, err
		}
		err = backupClient.backupService(namespaceName)
		if err != nil {
			return nil, err
		}
		err = backupClient.backupIngress(namespaceName)
		if err != nil {
			return nil, err
		}
		err = backupClient.backupPersistentVolumeClaim(namespaceName)
		if err != nil {
			return nil, err
		}
		err = backupClient.backupPersistentVolume()
		if err != nil {
			return nil, err
		}
	}

	// 压缩文件 zip
	if outputFile == "" {
		outputFile = kClient.Name + ".zip"
	}
	err = compress.Zip(*backupDirectory, outputFile, true)
	if err != nil {
		return nil, err
	}
	return &outputFile, nil
}

// CronBackup 备份
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
		tempDirectory, err := ctl.CreateTempDirectory("cron_kubernetes", dateTimeFormat)
		if err != nil {
			color.Red(fmt.Sprint(err))
			return
		}
		err = filepath.Walk(configPath, func(configItemPath string, fi os.FileInfo, errBack error) (err error) {
			var fileName = fi.Name()
			if !strings.HasSuffix(fileName, ".yaml") {
				return nil
			}
			var outFileName = path.Base(fileName) + "_" + dateTimeFormat + ".zip"
			var outputFile = path.Join(*tempDirectory, outFileName)
			backupZipFile, err := Backup(configItemPath, outputFile, nil)
			if err != nil {
				return err
			}
			_, err = cosClient.Put(*backupZipFile, "kubernetes/"+dateFormat+"/"+outFileName)
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
	color.Green(fmt.Sprintf("Cron (%s) Start Success ...", backupCron))
	color.Blue(fmt.Sprintf("ConfigPath: %s", configPath))
	color.Blue(fmt.Sprintf("S3 Config: %s", cloudStorageConfig))
	c.Start()
	select {}
}

func (backup *BackupClient) createDirectory(values ...string) (*string, error) {
	var localDirectoryPath = path.Join(backup.rootPath, path.Join(values...))
	if _, err := os.Stat(localDirectoryPath); err != nil || os.IsNotExist(err) {
		err := os.MkdirAll(localDirectoryPath, 644)
		if err != nil {
			return nil, err
		}
	}
	return &localDirectoryPath, nil
}

func (backup *BackupClient) output(value any, output ...string) error {
	itemJson, err := json.Marshal(value)
	if err != nil {
		return err
	}
	m := make(map[string]any)
	err = json.Unmarshal(itemJson, &m)
	if err != nil {
		return err
	}
	data, err := yaml.Marshal(m)
	if err != nil {
		return err
	}
	err = os.WriteFile(path.Join(output...), []byte(strings.ReplaceAll(string(data), "    ", "  ")), 644)
	if err != nil {
		return err
	}
	return nil
}

func (backup *BackupClient) backupNamespace(namespaceName string) error {
	color.Green(fmt.Sprintf("[Kubernetes] Backup Namespace: %s", namespaceName))
	localPath, err := backup.createDirectory("namespaces", namespaceName)
	if err != err {
		return err
	}
	namespace, err := backup.client.Namespace(namespaceName)
	if err != nil {
		return err
	}
	namespace.Kind = "Namespace"
	namespace.APIVersion = "v1"
	err = backup.output(namespace, *localPath, "namespace.yaml")
	if err != nil {
		return err
	}
	return nil
}

func (backup *BackupClient) backupDeployment(namespaceName string) error {
	localPath, err := backup.createDirectory("namespaces", namespaceName, "deployment")
	if err != err {
		return err
	}
	deployments, err := backup.client.Deployments(namespaceName)
	if err != err {
		return err
	}
	for _, item := range deployments {
		item.Kind = "Deployment"
		item.APIVersion = "apps/v1"
		color.Green(fmt.Sprintf("[Kubernetes] Backup Deployment: %s / %s", namespaceName, item.ObjectMeta.Name))
		err = backup.output(item, *localPath, item.ObjectMeta.Name+".yaml")
		if err != nil {
			return err
		}
	}
	return nil
}

func (backup *BackupClient) backupStatefulSet(namespaceName string) error {
	localPath, err := backup.createDirectory("namespaces", namespaceName, "statefulSet")
	if err != err {
		return err
	}
	statefulSets, err := backup.client.StatefulSets(namespaceName)
	if err != err {
		return err
	}
	for _, item := range statefulSets {
		item.Kind = "StatefulSet"
		item.APIVersion = "apps/v1"
		color.Green(fmt.Sprintf("[Kubernetes] Backup StatefulSet: %s / %s", namespaceName, item.ObjectMeta.Name))
		err = backup.output(item, *localPath, item.ObjectMeta.Name+".yaml")
		if err != nil {
			return err
		}
	}
	return nil
}

func (backup *BackupClient) backupDaemonSet(namespaceName string) error {
	localPath, err := backup.createDirectory("namespaces", namespaceName, "daemonSet")
	if err != err {
		return err
	}
	daemonSets, err := backup.client.DaemonSets(namespaceName)
	if err != err {
		return err
	}
	for _, item := range daemonSets {
		item.Kind = "DaemonSet"
		item.APIVersion = "apps/v1"
		color.Green(fmt.Sprintf("[Kubernetes] Backup DaemonSet: %s / %s", namespaceName, item.ObjectMeta.Name))
		err = backup.output(item, *localPath, item.ObjectMeta.Name+".yaml")
		if err != nil {
			return err
		}
	}
	return nil
}

func (backup *BackupClient) backupJob(namespaceName string) error {
	localPath, err := backup.createDirectory("namespaces", namespaceName, "job")
	if err != err {
		return err
	}
	jobs, err := backup.client.Jobs(namespaceName)
	if err != err {
		return err
	}
	for _, item := range jobs {
		item.Kind = "Job"
		item.APIVersion = "batch/v1"
		color.Green(fmt.Sprintf("[Kubernetes] Backup Job: %s / %s", namespaceName, item.ObjectMeta.Name))
		err = backup.output(item, *localPath, item.ObjectMeta.Name+".yaml")
		if err != nil {
			return err
		}
	}
	return nil
}

func (backup *BackupClient) backupCronJob(namespaceName string) error {
	localPath, err := backup.createDirectory("namespaces", namespaceName, "cronJob")
	if err != err {
		return err
	}
	cronJobs, err := backup.client.CronJobs(namespaceName)
	if err != err {
		return err
	}
	for _, item := range cronJobs {
		item.Kind = "CronJob"
		item.APIVersion = "batch/v1"
		color.Green(fmt.Sprintf("[Kubernetes] Backup CronJob: %s / %s", namespaceName, item.ObjectMeta.Name))
		err = backup.output(item, *localPath, item.ObjectMeta.Name+".yaml")
		if err != nil {
			return err
		}
	}
	return nil
}

func (backup *BackupClient) backupConfigmap(namespaceName string) error {
	localPath, err := backup.createDirectory("namespaces", namespaceName, "configmap")
	if err != err {
		return err
	}
	configmaps, err := backup.client.Configmaps(namespaceName)
	if err != err {
		return err
	}
	for _, item := range configmaps {
		item.Kind = "ConfigMap"
		item.APIVersion = "v1"
		color.Green(fmt.Sprintf("[Kubernetes] Backup ConfigMap: %s / %s", namespaceName, item.ObjectMeta.Name))
		err = backup.output(item, *localPath, item.ObjectMeta.Name+".yaml")
		if err != nil {
			return err
		}
	}
	return nil
}

func (backup *BackupClient) backupSecret(namespaceName string) error {
	localPath, err := backup.createDirectory("namespaces", namespaceName, "secret")
	if err != err {
		return err
	}
	secrets, err := backup.client.Secrets(namespaceName)
	if err != err {
		return err
	}
	for _, item := range secrets {
		item.Kind = "Secret"
		item.APIVersion = "v1"
		color.Green(fmt.Sprintf("[Kubernetes] Backup Secret: %s / %s", namespaceName, item.ObjectMeta.Name))
		err = backup.output(item, *localPath, item.ObjectMeta.Name+".yaml")
		if err != nil {
			return err
		}
	}
	return nil
}

func (backup *BackupClient) backupService(namespaceName string) error {
	localPath, err := backup.createDirectory("namespaces", namespaceName, "service")
	if err != err {
		return err
	}
	services, err := backup.client.Services(namespaceName)
	if err != err {
		return err
	}
	for _, item := range services {
		item.Kind = "Service"
		item.APIVersion = "v1"
		color.Green(fmt.Sprintf("[Kubernetes] Backup Service: %s / %s", namespaceName, item.ObjectMeta.Name))
		err = backup.output(item, *localPath, item.ObjectMeta.Name+".yaml")
		if err != nil {
			return err
		}
	}
	return nil
}

func (backup *BackupClient) backupIngress(namespaceName string) error {
	localPath, err := backup.createDirectory("namespaces", namespaceName, "ingress")
	if err != err {
		return err
	}
	secrets, err := backup.client.Secrets(namespaceName)
	if err != err {
		return err
	}
	for _, item := range secrets {
		item.Kind = "Ingress"
		item.APIVersion = "networking.k8s.io/v1"
		color.Green(fmt.Sprintf("[Kubernetes] Backup Ingress: %s / %s", namespaceName, item.ObjectMeta.Name))
		err = backup.output(item, *localPath, item.ObjectMeta.Name+".yaml")
		if err != nil {
			return err
		}
	}
	return nil
}

func (backup *BackupClient) backupPersistentVolumeClaim(namespaceName string) error {
	localPath, err := backup.createDirectory("namespaces", namespaceName, "persistentVolumeClaim")
	if err != err {
		return err
	}
	persistentVolumeClaims, err := backup.client.PersistentVolumeClaims(namespaceName)
	if err != err {
		return err
	}
	for _, item := range persistentVolumeClaims {
		item.Kind = "PersistentVolumeClaim"
		item.APIVersion = "v1"
		color.Green(fmt.Sprintf("[Kubernetes] Backup PersistentVolumeClaim: %s / %s", namespaceName, item.ObjectMeta.Name))
		err = backup.output(item, *localPath, item.ObjectMeta.Name+".yaml")
		if err != nil {
			return err
		}
	}
	return nil
}

func (backup *BackupClient) backupPersistentVolume() error {
	localPath, err := backup.createDirectory("persistentVolume")
	if err != err {
		return err
	}
	persistentVolumes, err := backup.client.PersistentVolumes()
	if err != err {
		return err
	}
	for _, item := range persistentVolumes {
		item.Kind = "PersistentVolume"
		item.APIVersion = "v1"
		color.Green(fmt.Sprintf("[Kubernetes] Backup PersistentVolume: %s", item.ObjectMeta.Name))
		err = backup.output(item, *localPath, item.ObjectMeta.Name+".yaml")
		if err != nil {
			return err
		}
	}
	return nil
}
