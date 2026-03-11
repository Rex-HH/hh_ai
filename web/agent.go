package main

import (
	"context"
	"hh_ai/agent"
	"hh_ai/agent/my_tool"
	"log"
	"sync"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
)

var (
	runner     *adk.Runner
	runnerOnce sync.Once
)

func GetRunner() *adk.Runner {
	runnerOnce.Do(func() {
		ctx := context.Background()
		callbacks.AppendGlobalHandlers(agent.GetChatModelInputCallback())

		// 创建 ChatModel
		model := agent.CreateDeepseekModel()
		// 创建工具集
		tools := []tool.BaseTool{
			my_tool.CreateTimeTool(),
			my_tool.CreateLocationTool(),
			my_tool.CreateWeatherTool(),
			my_tool.CreateCrawlTool(),
			my_tool.CreateDuckDuckSearchTool(),
		}
		// 创建 Agent
		agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
			Model: model,
			ToolsConfig: adk.ToolsConfig{
				ToolsNodeConfig: compose.ToolsNodeConfig{
					Tools: tools,
				},
			},
			Name:        "my_agent",
			Description: "我的万能AI助手", // 一个Agent可能是另一个Agent的SubAgent，所以Agent的功能需要描绘清楚
			Instruction: "",         // SystemMessage，支持FString渲染
		})
		if err != nil {
			log.Fatal(err)
		}
		// 创建 runner
		runner = adk.NewRunner(ctx, adk.RunnerConfig{
			Agent:           agent,
			EnableStreaming: true, //对于 Agent 内部能够流式输出的组件（如 ChatModel 调用），应以流的形式逐步返回结果。如果某个组件天然不支持流式（比如Tool），它仍然可以按其原有的非流式方式工作。
		})
		if err != nil {
			log.Fatal(err)
		}
	})
	return runner
}
