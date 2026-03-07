package eino_graph

import "testing"

func TestComposeGraph(t *testing.T) {
	ComposeGraph("hello, how are you")
	ComposeGraph("hello, how are you!")
}
func TestGraphWithCallBack(t *testing.T) {
	GraphWithCallBack("hello, how are you")
	GraphWithCallBack("hello, how are you!")
}

func TestGraphWithState(t *testing.T) {
	GraphWithState("hello, how are you")
	GraphWithState("hello, how are you!")
}
