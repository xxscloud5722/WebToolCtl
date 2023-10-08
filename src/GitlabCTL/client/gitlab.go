package client

import (
	"github.com/fatih/color"
	"github.com/xanzy/go-gitlab"
	"os"
	"strings"
	"time"
)

const (
	gitLabExportStatusFinished = "finished"
	gitlabWaitStatus           = 429
)

type GitlabClient struct {
	client *gitlab.Client
	URL    string `yaml:"url"`
	Token  string `yaml:"token"`
}

func NewGitlabClient(host, token string) (*GitlabClient, error) {
	var url = host
	if !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}
	if !strings.HasSuffix(url, "/api/v4") {
		url = url + "/api/v4"
	}
	git, err := gitlab.NewClient(token, gitlab.WithBaseURL(url))
	if err != nil {
		return nil, err
	}
	return &GitlabClient{client: git, URL: url, Token: token}, nil
}

// Projects 扫描授权下所有项目列表
func (g *GitlabClient) Projects() ([]*gitlab.Project, error) {
	var cursor = 1
	var allProjects []*gitlab.Project
	for {
		options := &gitlab.ListProjectsOptions{}
		options.OrderBy = gitlab.String("id")
		options.ListOptions.Page = cursor
		projects, response, err := g.client.Projects.ListProjects(options)
		if err != nil {
			return nil, err
		}
		for _, project := range projects {
			color.Green("[Gitlab] 项目ID: (%d) 项目路径: %s", project.ID, project.NameWithNamespace)
		}
		allProjects = append(allProjects, projects...)
		if response.NextPage == 0 {
			break
		}
		cursor++
	}
	return allProjects, nil
}

func (g *GitlabClient) Export(projectId int, backupFile string) error {
	_, err := g.client.ProjectImportExport.ScheduleExport(projectId, &gitlab.ScheduleExportOptions{})
	if err != nil {
		return err
	}
	// 查询导出状态 轮训到成功为止
	for {
		status, _, err := g.client.ProjectImportExport.ExportStatus(projectId)
		// 如果请求异常再次请求
		if err != nil {
			time.Sleep(time.Second * 1)
			continue
		}
		if status.ExportStatus == gitLabExportStatusFinished {
			break
		}
		time.Sleep(time.Second * 3)
	}

	// 下载导出文件
	for {
		download, response, err := g.client.ProjectImportExport.ExportDownload(projectId)
		if response.StatusCode == gitlabWaitStatus {
			time.Sleep(time.Second * 30)
			continue
		}
		// 如果请求异常再次请求
		if err != nil {
			time.Sleep(time.Second * 1)
			continue
		}
		// 先写入缓存然后在重命名, 防止文件损坏
		var cacheFile = backupFile + ".backup"
		err = os.WriteFile(cacheFile, download, 644)
		if err != nil {
			return err
		}
		// 重命名文件
		err = os.Rename(cacheFile, backupFile)
		if err != nil {
			return err
		}
		break
	}
	return nil
}
