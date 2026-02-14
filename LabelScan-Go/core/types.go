package core

type Image []float32

type Sample struct {
	ID          int    `json:"id"`
	Data        Image  `json:"pixels"`
	Label       int    `json:"label"`        // 真标签
	TargetLabel int    `json:"target_label"` // 影子重标标签
	IsMember    bool   `json:"is_member"`
	Filename    string `json:"filename"`
}

// AttackResult 审计战报：由队长 A 填写，你负责回收
type AttackResult struct {
	SampleID      int
	OriginalLabel int
	FinalLabel    int
	IsSuccess     bool
	Queries       int     // 攻击这张图查了多少次 API
	Distance      float64 // 到决策边界的距离（核心指标）
	IsMember      bool    // 样本真身
}

// Model 接口：队长实现的对讲机（你现在要求他必须支持批量）
type Model interface {
	Predict(img Image) (int, error)
	PredictBatch(imgs []Image) ([]int, error) // 涡轮增压接口
}

// Attacker 接口：队长实现的攻击逻辑（如：边界攻击）
type Attacker interface {
	Attack(model Model, sample Sample) AttackResult
}
