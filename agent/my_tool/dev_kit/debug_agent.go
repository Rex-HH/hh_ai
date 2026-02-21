package dev_kit

import (
	"context"
	"fmt"
	"hh_ai/agent"
	"hh_ai/agent/my_tool"
	"log"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func TripAgent() {
	ctx := context.Background()
	//callbacks.AppendGlobalHandlers(new(agent.LoggerCallbacks))
	//callbacks.AppendGlobalHandlers(agent.GetStartCallback())
	//callbacks.AppendGlobalHandlers(agent.GetEndCallback())
	callbacks.AppendGlobalHandlers(agent.GetChatModelInputCallback())
	callbacks.AppendGlobalHandlers(agent.GetToolInputCallback())
	model := agent.CreateDeepseekModel()

	tools := []tool.BaseTool{
		my_tool.CreateTimeTool(),
		my_tool.CreateLocationTool(),
		my_tool.CreateWeatherTool(),
	}

	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "trip_plan",
		Description: "出行规划",
		Instruction: "请帮我完成旅行规划，包括什么时间去什么景点，" +
			"并列出详细的高铁车次和出发时间", // SystemMessage，支持FString渲染
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

	iter := runner.Query(ctx, "我想去北京旅游3天(从明天开始)，请帮我做一份旅游攻略")

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
