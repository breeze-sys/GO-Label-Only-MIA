package worker

import (
	"LabelScan-Go/core"
	"fmt"
	"sync"
)

type Relabeler struct {
	Model     core.Model
	BatchSize int // 建议 128
}

func NewRelabeler(m core.Model, batchSize int) *Relabeler {
	return &Relabeler{Model: m, BatchSize: batchSize}
}

// RelabelAll 使用“大卡车运输”模式，一次处理 128 张图
func (r *Relabeler) RelabelAll(samples []core.Sample) {
	var wg sync.WaitGroup
	n := len(samples)
	numBatches := (n + r.BatchSize - 1) / r.BatchSize
	batchChan := make(chan int, numBatches)

	fmt.Printf("🌀 涡轮重标引擎：每批 %d 张，共 %d 批\n", r.BatchSize, numBatches)

	for i := 0; i < 5; i++ { // 开启 5 个批量工人即可，防止网络拥堵
		wg.Add(1)
		go func() {
			defer wg.Done()
			for bIdx := range batchChan {
				start := bIdx * r.BatchSize
				end := start + r.BatchSize
				if end > n {
					end = n
				}

				// 1. 打包一批图
				var batchImgs []core.Image
				for j := start; j < end; j++ {
					batchImgs = append(batchImgs, samples[j].Data)
				}

				// 2. 一次性问 Python 拿回一堆标签
				labels, err := r.Model.PredictBatch(batchImgs)
				if err == nil {
					for k, label := range labels {
						samples[start+k].TargetLabel = label
					}
				}
			}
		}()
	}

	for i := 0; i < numBatches; i++ {
		batchChan <- i
	}
	close(batchChan)
	wg.Wait()
}
