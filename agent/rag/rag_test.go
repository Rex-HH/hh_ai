package rag

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bytedance/sonic"
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

func TestSearchByVector(t *testing.T) {
	SearchByVector()
}

func TestRetrieveDocument(t *testing.T) {
	IndexDocument("../../data/qa.md")
	docs := RetrieveDocument("在美国，一个企业雇一个员工，医疗和社会保障方面的支出有多少？", 4)
	for _, doc := range docs {
		s, _ := sonic.MarshalString(doc)
		fmt.Println(s)
	}
	GetQdrantClient().Close()
}

func TestChatBot(t *testing.T) {
	IndexDocument("../../data/qa.md")
	answer := ChatBot("在美国，一个企业雇一个员工，医疗和社会保障方面的支出有多少？")
	fmt.Println(answer)
	GetQdrantClient().Close()
}
