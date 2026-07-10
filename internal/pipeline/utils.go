package pipeline

import "fmt"

func FormatTime(seconds float64) string {
	totalSeconds := int(seconds)
	hours := totalSeconds / 3600
	remainder := totalSeconds % 3600
	minutes := remainder / 60
	secs := remainder % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
}
