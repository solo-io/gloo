package utils

import (
	"fmt"
	"strings"
)

// ConsolidateLog is a helper function to dedup the log message.
func ConsolidateLog(logMessage string) string {
	logCountMap := make(map[string]int)
	stderrSlice := strings.Split(logMessage, "\n")
	for _, item := range stderrSlice {
		if item == "" {
			continue
		}
		_, exist := logCountMap[item]
		if exist {
			logCountMap[item]++
		} else {
			logCountMap[item] = 1
		}
	}
	var sb strings.Builder
	for _, item := range stderrSlice {
		if logCountMap[item] == 0 {
			continue
		}
		sb.WriteString(fmt.Sprintf("%s (repeated %v times)\n", item, logCountMap[item]))
		// reset seen log count
		logCountMap[item] = 0
	}
	return sb.String()
}

func Log(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args)
}
