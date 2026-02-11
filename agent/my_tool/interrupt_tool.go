package my_tool

import (
	"context"
	"log"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
)

const AskForClarificationToolName = "ask_for_clarification"

/*
支持中断的Tool。当用户未提供必要信息时，向用户问清楚。【固定写法】
*/
type askForClarificationOptions struct {
	NewInput *string
}

func WithNewInput(input string) tool.Option {
	return tool.WrapImplSpecificOptFn(func(t *askForClarificationOptions) {
		t.NewInput = &input
	})
}

type AskForClarificationInput struct {
	Question string `json:"question" jsonschema:"description=为了获得必要的缺失信息，你需要向用户询问的问题"`
}

// NewAskForClarificationTool 构造支持中断的Tool
func NewAskForClarificationTool() tool.InvokableTool {
	t, err := utils.InferOptionableTool(
		AskForClarificationToolName,
		"当用户的请求含糊不清或缺乏继续所需的信息时，调用此工具。在你能有效地使用其他工具之前，用它来问一个后续问题，以获得你需要的细节。",
		func(ctx context.Context, input *AskForClarificationInput, opts ...tool.Option) (output string, err error) {
			o := tool.GetImplSpecificOptions[askForClarificationOptions](nil, opts...)
			if o.NewInput == nil {
				return "", compose.Interrupt(ctx, input.Question)
			}
			return *o.NewInput, nil
		})
	if err != nil {
		log.Fatal(err)
	}
	return t
}
