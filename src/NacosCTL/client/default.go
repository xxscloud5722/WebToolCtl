package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type NacosClient struct {
	Host  string
	Token string
}

type NacosResponse[T []NacosNamespace] struct {
	Code int `json:"code"`
	Data T   `json:"data"`
}
type NacosPageResponse[T []NacosConfigItem] struct {
	TotalCount int `json:"totalCount"`
	PageNumber int `json:"pageNumber"`
	PageItems  T   `json:"pageItems"`
}
type NacosNamespace struct {
	Namespace         string `json:"namespace"`
	NamespaceShowName string `json:"namespaceShowName"`
	ConfigCount       int    `json:"configCount"`
}
type NacosConfigItem struct {
	Namespace string `json:"tenant"`
	Id        string `json:"id"`
	DataId    string `json:"dataId"`
	Group     string `json:"group"`
	Content   string `json:"content"`
	Md5       string `json:"md5"`
	Type      string `json:"type"`
}

func NewNacosClient(nacosHost, username, password string) (*NacosClient, error) {
	if !strings.HasPrefix(nacosHost, "http") {
		nacosHost = "http://" + nacosHost
	}
	if !strings.HasSuffix(nacosHost, ":8080") {
		nacosHost = nacosHost + ":8080"
	}
	response, err := http.PostForm(nacosHost+"/nacos/v1/auth/users/login", url.Values{
		"username": {username},
		"password": {password},
	})
	if err != nil {
		return nil, err
	}
	var result struct {
		AccessToken string `json:"accessToken"`
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(response.Body)
	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(responseBytes, &result)
	if err != nil {
		return nil, err
	}
	return &NacosClient{Host: nacosHost, Token: result.AccessToken}, nil
}

func (nacos *NacosClient) Namespaces() (*NacosResponse[[]NacosNamespace], error) {
	response, err := http.Get(fmt.Sprintf("%s/nacos/v1/console/namespaces?&accessToken=%s&namespaceId=", nacos.Host, nacos.Token))
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(response.Body)
	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Response Status Code: %d", response.StatusCode))
	}
	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var result NacosResponse[[]NacosNamespace]
	err = json.Unmarshal(responseBytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (nacos *NacosClient) NamespaceItems(namespaceId string) ([]NacosConfigItem, error) {
	if namespaceId == "" {
		return nil, nil
	}
	response, err := http.Get(fmt.Sprintf("%s/nacos/v1/cs/configs?dataId=&group=&appName=&config_tags=&"+
		"pageNo=1&pageSize=500&tenant=%s&search=accurate&accessToken=%s&username=nacos", nacos.Host, namespaceId, nacos.Token))
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(response.Body)
	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Response Status Code: %d", response.StatusCode))
	}
	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var result NacosPageResponse[[]NacosConfigItem]
	err = json.Unmarshal(responseBytes, &result)
	if err != nil {
		return nil, err
	}
	return result.PageItems, nil
}
