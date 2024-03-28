package helpers

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// synchronous timer aggregator
var components []string
var lock sync.Mutex

func ListTimers() {
	fmt.Println("Timer Stats:")

	for _, component := range components {
		fmt.Println(component)
	}
}

func TimerFunc(component string) (stop func()) {
	start := time.Now()

	return func() {
		lock.Lock()
		defer lock.Unlock()
		components = append(components, fmt.Sprintf("  Step: %s took time: %s", strings.ToUpper(component), time.Since(start).Truncate(time.Second).String()))
	}
}
