package main

import (
	"LabelScan-Go/core"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

// ExportAttackResults 导出最终审计成绩单 (Task 3 结果)
func ExportAttackResults(results []core.AttackResult, filename string) {
	file, _ := os.Create(filename)
	defer file.Close()
	w := csv.NewWriter(file)
	defer w.Flush()

	w.Write([]string{"id", "orig", "final", "success", "queries", "distance", "is_member"})
	for _, r := range results {
		w.Write([]string{
			strconv.Itoa(r.SampleID),
			strconv.Itoa(r.OriginalLabel),
			strconv.Itoa(r.FinalLabel),
			strconv.FormatBool(r.IsSuccess),
			strconv.Itoa(r.Queries),
			fmt.Sprintf("%.6f", r.Distance),
			strconv.FormatBool(r.IsMember),
		})
	}
	fmt.Printf("💾 审计报告已保存至: %s\n", filename)
}
