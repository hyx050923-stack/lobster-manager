package container

import (
	"encoding/json"
	"io"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/pkg/jsonmessage"
)

const OpenClawImage = "openclaw/openclaw:latest" // 假设官方镜像名，实际需要确认

// PullImage 拉取OpenClaw镜像，支持进度回调
func (m *Manager) PullImage(imageName string, progress func(current, total int64, status string)) error {
	// 检查本地是否已有镜像
	filter := filters.NewArgs()
	filter.Add("reference", imageName)
	images, err := m.cli.ImageList(m.ctx, types.ImageListOptions{Filters: filter})
	if err == nil && len(images) > 0 {
		// 已存在，直接返回
		if progress != nil {
			progress(100, 100, "镜像已存在")
		}
		return nil
	}

	// 拉取镜像
	resp, err := m.cli.ImagePull(m.ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer resp.Close()

	// 解析进度
	decoder := json.NewDecoder(resp)
	for {
		var msg jsonmessage.JSONMessage
		if err := decoder.Decode(&msg); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if msg.Error != nil {
			return msg.Error
		}
		if progress != nil && msg.Progress != nil {
			progress(msg.Progress.Current, msg.Progress.Total, msg.Status)
		}
	}
	return nil
}