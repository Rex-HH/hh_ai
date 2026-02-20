package dev_kit

import (
	"hh_ai/agent/my_tool"
	"testing"
)

func TestWeatherAgent(t *testing.T) {
	WeatherAgent()
}

func TestLoopAgent(t *testing.T) {
	// 待优化的代码示例
	code := `
func processData(data []int) []int {
    result := []int{}
    for i := 0; i < len(data); i++ {
        for j := 0; j < len(data); j++ {
            if data[i] > data[j] {
                result = append(result, data[i])
                break
            }
        }
    }
    return result
}`
	LoopAgent(code)
}

func TestParallelAgent(t *testing.T) {
	ParallelAgent()
}

func TestSearchAgent(t *testing.T) {
	defer my_tool.CloseBrowser()
	SearchAgent()
}
