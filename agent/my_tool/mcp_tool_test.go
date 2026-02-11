package my_tool

import "testing"

func TestList12306Tool(t *testing.T) {
	InitMcpClient()
	List12306Tool()
	Use12306Tool()
}

func TestUse12306Tool(t *testing.T) {
	InitMcpClient()
	Use12306Tool()
}
