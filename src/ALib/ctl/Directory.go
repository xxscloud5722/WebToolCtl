package ctl

import (
	"github.com/olekukonko/tablewriter"
	"os"
	"path"
	"runtime"
)

func CreateLocalRootDirectory(rootPath, name string) (*string, error) {
	if rootPath == "" {
		workspace, err := GetWorkspace()
		if err != nil {
			return nil, err
		}
		rootPath = *workspace
	}
	var localDirectoryPath = path.Join(rootPath, name)
	if _, err := os.Stat(localDirectoryPath); err != nil || os.IsNotExist(err) {
		err = os.MkdirAll(localDirectoryPath, 644)
		if err != nil {
			return nil, err
		}
	}
	return &localDirectoryPath, nil
}

func CreateTempDirectory(values ...string) (*string, error) {
	var tempPath string
	if runtime.GOOS == "windows" {
		tempPath = "C:\\Windows\\Temp"
	} else {
		tempPath = "/tmp"
	}
	var tempDirectory = path.Join(tempPath, path.Join(values...))
	if _, err := os.Stat(tempDirectory); err != nil || os.IsNotExist(err) {
		err = os.MkdirAll(tempDirectory, 644)
		if err != nil {
			return nil, err
		}
	}
	return &tempDirectory, nil
}

// GetWorkspace 获取程序的工作空间
func GetWorkspace() (*string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return &dir, nil
}

// PrintTable 输出Table
func PrintTable(header []string, dataSources [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	for _, v := range dataSources {
		table.Append(v)
	}
	table.Render()
}
