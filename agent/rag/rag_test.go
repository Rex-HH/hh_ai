package rag

import (
	"fmt"
	"strings"
	"testing"
)

func TestParseDocument(t *testing.T) {
	ParseDocument("../../data/qa.md")
	fmt.Println(strings.Repeat("~", 100))
	ParseDocument("../../data/qa.html")
	fmt.Println(strings.Repeat("~", 100))
	ParseDocument("../../data/qa.pdf")
	fmt.Println(strings.Repeat("~", 100))
	ParseDocument("../../data/df.xlsx")
	fmt.Println(strings.Repeat("~", 100))
	ParseDocument("../../data/观沧海.docx")
	fmt.Println(strings.Repeat("~", 100))
}
func TestLoadDocument(t *testing.T) {
	LoadDocument("../../data/qa.md") // 相对路径是相对于本_test.go文件的路径
	fmt.Println(strings.Repeat("~", 100))
	LoadDocument("../../data/qa.html")
	fmt.Println(strings.Repeat("~", 100))
	LoadDocument("../../data/qa.pdf")
	fmt.Println(strings.Repeat("~", 100))
	LoadDocument("../../data/df.xlsx")
	fmt.Println(strings.Repeat("~", 100))
	LoadDocument("../../data/观沧海.docx")
	fmt.Println(strings.Repeat("~", 100))
}

func TestTransformDocument(t *testing.T) {
	TransformDocument("../../data/qa.md")
}

func TestWrodSim(t *testing.T) {
	WrodSim()
}
