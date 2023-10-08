package storage

import (
	"context"
	"fmt"
	"github.com/fatih/color"
)

func (c *TencentCosClient) Put(localPath, cloudPath string) (*string, error) {
	color.Blue(fmt.Sprintf("[Cloud Storage] Put: %s -> %s", localPath, cloudPath))
	_, _, err := c.client.Object.Upload(context.Background(), cloudPath, localPath, nil)
	if err != nil {
		return nil, err
	}
	return &cloudPath, err
}
