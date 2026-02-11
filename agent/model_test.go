package agent

import (
	"fmt"
	"strings"
	"testing"
)

func TestRunChatModel(t *testing.T) {
	RunChatModel(false)
	fmt.Println(strings.Repeat("-", 100))
	RunChatModel(true)
}

func TestRunChatModelAgent(t *testing.T) {
	RunChatModelAgent()
}

// go test -v ./agent -run=^TestRunChatModel$ -count=1
