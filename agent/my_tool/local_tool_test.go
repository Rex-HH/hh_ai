package my_tool

import (
	"fmt"
	"testing"
)

func TestUseTimeTool(t *testing.T) {
	fmt.Println("纽约时间")
	UseTimeTool("America/New_York")
	fmt.Println("伦敦时间")
	UseTimeTool("Europe/London")
	fmt.Println("上海时间")
	UseTimeTool("")
}
