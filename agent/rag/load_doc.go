package rag

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	file2 "github.com/cloudwego/eino-ext/components/document/loader/file"
	"github.com/cloudwego/eino-ext/components/document/parser/docx"
	"github.com/cloudwego/eino-ext/components/document/parser/html"
	"github.com/cloudwego/eino-ext/components/document/parser/pdf"
	"github.com/cloudwego/eino-ext/components/document/parser/xlsx"
	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/components/document/parser"
	"github.com/cloudwego/eino/schema"
)

func ParseDocument(filepath string) {
	ctx := context.Background()
	textParser := parser.TextParser{}
	selector := "body"
	htmlParser, _ := html.NewParser(ctx, &html.Config{Selector: &selector})
	pdfParser, _ := pdf.NewPDFParser(ctx, &pdf.Config{})
	xlsxParser, _ := xlsx.NewXlsxParser(ctx, &xlsx.Config{})
	wordParser, _ := docx.NewDocxParser(ctx, &docx.Config{
		IncludeHeaders: false,
		IncludeFooters: false,
		IncludeTables:  true,
	})

	// 创建扩展解析器
	extParser, _ := parser.NewExtParser(ctx, &parser.ExtParserConfig{
		Parsers: map[string]parser.Parser{
			".html": htmlParser,
			".xlsx": xlsxParser,
			".pdf":  pdfParser,
			".docx": wordParser,
		},
		FallbackParser: textParser,
	})

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("文件%s打开失败：%s", filepath, err)
	}
	defer file.Close()

	docs, err := extParser.Parse(ctx, file, parser.WithURI(filepath))
	if err != nil {
		log.Fatal(err)
	}
	for _, doc := range docs {
		fmt.Println(doc.Content)
	}
}

func LoadDocument(filepath string) {
	ctx := context.Background()
	textParser := parser.TextParser{}
	selector := "body"
	htmlParser, _ := html.NewParser(ctx, &html.Config{Selector: &selector})
	pdfParser, _ := pdf.NewPDFParser(ctx, &pdf.Config{})
	xlsxParser, _ := xlsx.NewXlsxParser(ctx, &xlsx.Config{})
	wordParser, _ := docx.NewDocxParser(ctx, &docx.Config{
		IncludeHeaders: false,
		IncludeFooters: false,
		IncludeTables:  true,
	})

	// 创建扩展解析器
	extParser, _ := parser.NewExtParser(ctx, &parser.ExtParserConfig{
		Parsers: map[string]parser.Parser{
			".html": htmlParser,
			".xlsx": xlsxParser,
			".pdf":  pdfParser,
			".docx": wordParser,
		},
		FallbackParser: textParser,
	})

	loader, _ := file2.NewFileLoader(ctx, &file2.FileLoaderConfig{
		UseNameAsID: true,
		Parser:      extParser,
	})

	docs, err := loader.Load(ctx, document.Source{URI: filepath})
	if err != nil {
		log.Fatal(err)
	}
	for _, doc := range docs {
		fmt.Println(doc.Content)
	}

}

func TransformDocument(markdownFile string) {
	ctx := context.Background()
	file, err := os.Open(markdownFile)
	if err != nil {
		log.Fatalf("文件%s打开失败：%s", markdownFile, err)
	}
	defer file.Close()
	content, _ := io.ReadAll(file)
	doc := &schema.Document{
		Content: string(content),
	}

	splitTransformer, _ := markdown.NewHeaderSplitter(ctx, &markdown.HeaderConfig{
		Headers: map[string]string{
			"##": "",
		},
		TrimHeaders: false,
	})

	transformedDoc, _ := splitTransformer.Transform(ctx, []*schema.Document{doc})
	for idx, seg := range transformedDoc {
		fmt.Printf("segment %d, content %s\n", idx, seg.Content)
	}
}
