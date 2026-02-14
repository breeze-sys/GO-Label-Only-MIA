package dataset

import (
	"LabelScan-Go/core"
	"bufio"
	"fmt"
	"io"
	"os"
)

type CifarLoader struct {
	IsMemberSet bool // 标记当前加载的文件是否属于训练集
}

// LoadBatch 加载指定数量的数据
func (l *CifarLoader) LoadBatch(path string, limit int) ([]core.Sample, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var samples []core.Sample
	reader := bufio.NewReader(file)

	// CIFAR-10 每条记录固定 3073 字节 (1 label + 3072 pixels)
	const recordSize = 3073
	buffer := make([]byte, recordSize)

	count := 0
	for {
		// 达到 limit 则停止读取 (-1 代表读取全量)
		if limit != -1 && count >= limit {
			break
		}

		_, err := io.ReadFull(reader, buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// 1. 提取标签 (第1个字节)
		label := int(buffer[0])

		// 2. 提取并转换像素 (归一化到 0.0 - 1.0)
		pixels := make(core.Image, 3072)
		for i := 0; i < 3072; i++ {
			pixels[i] = float32(buffer[i+1]) / 255.0
		}

		// 3. 组装并生成 Filename
		samples = append(samples, core.Sample{
			ID:       count,
			Data:     pixels,
			Label:    label,
			IsMember: l.IsMemberSet,
			// 自动生成 Filename 方便队长 Debug: "路径_序号"
			Filename: fmt.Sprintf("%s_#%d", path, count),
		})
		count++
	}

	return samples, nil
}
