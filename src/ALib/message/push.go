package message

import (
	"bytes"
	"encoding/json"
	"github.com/fatih/color"
	"net/http"
	"strings"
	"time"
)

var client = http.Client{}

const DomainType = "DOMAIN"

var handler = map[string]func(config string, title string, message string) error{
	"DOMAIN:CP_WECHAT": func(config string, title string, message string) error {
		var content = "## " + title + "\n" +
			"> 程序版本号：**1.0.2** \n" +
			"> 检查时间：**#{now}**\n"
		content += message + "\n"
		content = ParseContent(content)

		color.Green("[企业微信机器人] 开始推送消息: " + content)

		var rows = strings.Split(content, "\n")
		var length = len(rows)
		var index int
		if length%60 == 0 {
			index = length / 60
		} else {
			index = length/60 + 1
		}

		for i := 0; i < index; i++ {
			if index-1 == i {
				content = strings.Join(rows[i*60:], "\n")
			} else {
				content = strings.Join(rows[i*60:(i+1)*60], "\n")
			}
			var requestBody = map[string]interface{}{}
			var contentRequest = map[string]string{}
			contentRequest["content"] = content
			requestBody["msgtype"] = "markdown"
			requestBody["markdown"] = contentRequest
			requestByteData, err := json.Marshal(requestBody)
			if err != nil {
				panic(err)
			}
			request, err := http.NewRequest("POST", config, bytes.NewReader(requestByteData))
			if err != nil {
				panic(err)
			}
			_, err = client.Do(request)
			if err != nil {
				return err
			}
			time.Sleep(time.Second)
		}
		return nil
	},
}

func Push(messageFormatType, config string, title string, message string) error {
	var configs = strings.Split(config, ",")
	var messageType = configs[0]
	config = config[len(messageType)+1:]
	if fun, ok := handler[messageFormatType+":"+messageType]; ok {
		return fun(config, title, message)
	}
	return nil
}

func ParseContent(message string) string {
	return strings.ReplaceAll(message, "#{now}", time.Now().Format(time.DateOnly))
}
