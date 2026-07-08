package logger

import (
	"fmt"
	"sync"
)

var (
	nodeID int
	mu     sync.Mutex
)

func Init(id int) {
	mu.Lock()
	defer mu.Unlock()
	nodeID = id
}

func Info(format string, args ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	prefix := fmt.Sprintf("[nó %d] ", nodeID)
	fmt.Printf(prefix+format+"\n", args...)
}