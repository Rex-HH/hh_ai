package dev_kit

import (
	"context"
	"fmt"
	"hh_ai/agent"
	"hh_ai/agent/my_tool"
	"log"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func SearchAgent() {
	ctx := context.Background()

	model := agent.CreateDeepseekModel()

	tools := []tool.BaseTool{
		my_tool.CreateCrawlTool(),
		my_tool.CreateDuckDuckSearchTool(),
	}

	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "internet_infomation_collector",
		Description: "互联网信息整合大师",
		Instruction: "请你积极地使用搜索引擎和爬虫获取互联网上的最新信息，来回答我的问题。" +
			"每个问题爬取的网页数不要超过3个", // SystemMessage，支持FString渲染
		Model: model,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: tools,
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// 运行 Agent
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		EnableStreaming: false,
		Agent:           agent,
	})
	//iter := runner.Query(ctx, "海淀高中前三强")
	iter := runner.Query(ctx, "黑龙江省计算机专业最强的三所大学")

	var lastmsg adk.Message
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			log.Println(event.Err)
		}
		message, err := event.Output.MessageOutput.GetMessage()
		if err != nil {
			log.Println(err)
		}
		lastmsg = message
		log.Println(message)
	}

	// 打印最终的结果
	fmt.Println("[最终结果]")
	if lastmsg.Role == schema.Assistant && len(lastmsg.Content) > 0 {
		fmt.Println(lastmsg.Content)
	}
}
