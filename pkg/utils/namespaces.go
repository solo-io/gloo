package utils

import (
	"os"
	"strconv"
	"strings"
	"sync"
)

func AllNamespaces(watchNamespaces []string) bool {

	if len(watchNamespaces) == 0 {
		return true
	}
	if len(watchNamespaces) == 1 && watchNamespaces[0] == "" {
		return true
	}
	return false
}

func ProcessWatchNamespaces(watchNamespaces []string, writeNamespace string) []string {
	if AllNamespaces(watchNamespaces) {
		return watchNamespaces
	}

	var writeNamespaceProvided bool
	for _, ns := range watchNamespaces {
		if ns == writeNamespace {
			writeNamespaceProvided = true
			break
		}
	}
	if !writeNamespaceProvided {
		watchNamespaces = append(watchNamespaces, writeNamespace)
	}

	return watchNamespaces
}

func GetPodNamespace() string {
	if podNamespace := os.Getenv("POD_NAMESPACE"); podNamespace != "" {
		return podNamespace
	}
	return "gloo-system"
}

var labelsOnce sync.Once

func GetPodLabels() map[string]string {
	var labels map[string]string
	labelsOnce.Do(func() {
		data, err := os.ReadFile("/etc/gloo/labels")
		if err != nil {
			return
		}

		m := map[string]string{}
		lines := strings.Split(string(data), "\n")
		for _, l := range lines {
			l = strings.TrimSpace(l)
			parts := strings.SplitN(l, "=", 2)
			if len(parts) != 2 {
				continue
			}

			key := parts[0]
			value, err := strconv.Unquote(parts[1])
			if err != nil {
				continue
			}
			m[key] = value
		}
		labels = m
	})
	return labels
}
