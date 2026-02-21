package model

import (
	"fmt"
	"math"
)

// StorageInfo 存储信息模型
type StorageInfo struct {
	Limit              int64   `json:"limit"`
	Used               int64   `json:"used"`
	Remaining          int64   `json:"remaining"`
	Percentage         float64 `json:"percentage"`
	LimitFormatted     string  `json:"limitFormatted"`
	UsedFormatted      string  `json:"usedFormatted"`
	RemainingFormatted string  `json:"remainingFormatted"`
}

// FormatFileSize 格式化文件大小
func FormatFileSize(bytes int64) string {
	if bytes <= 0 {
		return "0 Bytes"
	}

	units := []string{"Bytes", "KB", "MB", "GB", "TB"}
	digitGroups := int(math.Log10(float64(bytes)) / math.Log10(1024))
	if digitGroups >= len(units) {
		digitGroups = len(units) - 1
	}

	return fmt.Sprintf("%.1f %s", float64(bytes)/math.Pow(1024, float64(digitGroups)), units[digitGroups])
}
