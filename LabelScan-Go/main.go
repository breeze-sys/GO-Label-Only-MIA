package main

import (
	"LabelScan-Go/core"
	"LabelScan-Go/dataset"
	"LabelScan-Go/worker"
)

// --- 模拟对象 (等到联调时换成队长的真实代码) ---
type MockModel struct{}

func (m *MockModel) Predict(img core.Image) (int, error) { return 1, nil }
func (m *MockModel) PredictBatch(imgs []core.Image) ([]int, error) {
	res := make([]int, len(imgs))
	for i := range res {
		res[i] = 1
	}
	return res, nil
}

type MockAttacker struct{}

func (a *MockAttacker) Attack(m core.Model, s core.Sample) core.AttackResult {
	return core.AttackResult{
		SampleID: s.ID, Distance: 0.85, Queries: 100, IsMember: s.IsMember, IsSuccess: true,
	}
}

func main() {
	// 1. 加载 100 张成员和 100 张非成员 (Task 1)
	mLoader := &dataset.CifarLoader{IsMemberSet: true}
	nmLoader := &dataset.CifarLoader{IsMemberSet: false}
	m, _ := mLoader.LoadBatch("data/data_batch_1.bin", 100)
	nm, _ := nmLoader.LoadBatch("data/test_batch.bin", 100)
	allSamples := append(m, nm...)

	// 2. 并发批量重标 (Task 2)
	relabeler := worker.NewRelabeler(&MockModel{}, 128)
	relabeler.RelabelAll(allSamples)

	// 3. 通用高并发审计 (Task 3)
	auditor := worker.NewAuditor(&MockModel{}, &MockAttacker{}, 20)
	finalResults := auditor.RunAudit(allSamples)

	// 4. 导出 CSV (持久化)
	ExportAttackResults(finalResults, "final_audit_score.csv")
}
