package main

import (
	"fmt"
	"math/rand"
	"time"

	// ⚠️ 注意：这里要换成你 go.mod 里的模块名
	"Label-Only-MIA-Go/pkg/attack"
	"Label-Only-MIA-Go/pkg/core"
	"Label-Only-MIA-Go/pkg/mathutils"
)

// ==========================================
// 1. 定义一个简单的虚拟模型 (Mock Model)
// ==========================================
type SimpleModel struct{}

// GetInputSize 假设输入只有 10 个像素，方便观察
func (m *SimpleModel) GetInputSize() int {
	return 10
}

// Predict 简单的线性决策边界
// 规则：如果 input[0] (第一个像素) > 0.5，则是 Label 1，否则是 Label 0
func (m *SimpleModel) Predict(img core.Image) (int, error) {
	if len(img) > 0 && img[0] > 0.5 {
		return 1, nil
	}
	return 0, nil
}

// ==========================================
// 2. 主函数
// ==========================================
func main() {
	// A. 设置随机种子，保证每次运行结果一致 (方便调试)
	rand.Seed(time.Now().UnixNano())
	mathutils.SetSeed(42) // 如果你的 mathutils 有 SetSeed

	fmt.Println("=== 开始测试 HSJA 攻击算法 ===")

	// B. 准备数据
	// 创建一个全为 0.2 的图片 (它在 SimpleModel 中应该被判为 Label 0)
	inputSize := 10
	originalData := make(core.Image, inputSize)
	for i := range originalData {
		originalData[i] = 0.2
	}

	sample := core.Sample{
		ID:    1,
		Data:  originalData,
		Label: 0, // 真实标签是 0
	}

	// C. 初始化模型和攻击者
	model := &SimpleModel{}
	
	// 配置参数：为了测试快一点，迭代次数设少一点
	config := attack.HSJAConfig{
		MaxQueries:    1000,
		MaxIterations: 10,  // 迭代 10 轮
		NumEvals:      20,  // 每次梯度估计采样 20 次
		InitEvals:     20,  // 初始化采样 20 次
		ClipMin:       0.0,
		ClipMax:       1.0,
	}
	
	hsja := attack.NewHSJA(config)

	// D. 执行攻击
	fmt.Printf("原始数据[0]: %.4f, 原始标签: %d\n", sample.Data[0], sample.Label)
	startTime := time.Now()
	
	result := hsja.Attack(sample, model)
	
	duration := time.Since(startTime)

	// E. 验证结果
	fmt.Println("\n=== 攻击结果分析 ===")
	fmt.Printf("攻击耗时: %v\n", duration)
	fmt.Printf("是否成功: %v\n", result.IsSuccess)
	fmt.Printf("最终标签: %d\n", result.FinalLabel)
	fmt.Printf("查询次数: %d\n", result.Queries)
	fmt.Printf("L2 距离:  %.6f\n", result.Distance)
	
	// 我们知道边界是 0.5
	// 攻击成功的样本，其第一个像素应该略大于 0.5 (例如 0.501)
	// 如果是 0.8 或 0.9，说明攻击虽然成功了，但还没收敛到最优 (HSJA 的目的是贴近边界)
	// 如果是 0.5001，说明效果非常好
	// 此时还需要拿出攻击后的数据来看看
	// 注意：Result 里通常不存 Data，我们需要改一下 Attack 代码让它返回 Data，
	// 或者在 Attack 函数里加个 Log 打印最终的 xAdv[0]。
	// 这里我们假设我们信任 Distance。
	
	if result.IsSuccess {
		fmt.Println("\n✅ 测试通过！算法能够跨越决策边界。")
		if result.Distance < 0.35 { 
			// 0.5 - 0.2 = 0.3，理论最小距离是 0.3
			fmt.Println("🌟 优秀！结果非常接近理论最小距离 (0.3)。")
		} else {
			fmt.Println("⚠️ 注意：结果虽然成功，但距离较大，可能需要增加 MaxIterations。")
		}
	} else {
		fmt.Println("❌ 测试失败：未能改变标签。请检查 initialize 或 binarySearch 逻辑。")
	}
}
