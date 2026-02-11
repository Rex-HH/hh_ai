package dev_kit

import (
	"bytes"
	"context"
	"fmt"
	"hh_ai/agent"
	"hh_ai/agent/my_tool"
	"log"
	"time"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func InterruptResume() {
	ctx := context.Background()
	model := agent.CreateDeepseekModel()

	tools := []tool.BaseTool{
		//my_tool.CreateLocationTool(),
		my_tool.CreateWeatherTool(),
		my_tool.CreateTimeTool(),
	}

	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Model: model,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: tools,
			},
		},
		Name:        "activity_arrangement",
		Description: "活动计划大师",
		Instruction: "请根据未来的天气安排活动",
	})
	if err != nil {
		log.Fatal(err)
	}

	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent: agent,
	})

	messages := []adk.Message{
		&schema.Message{
			Role: schema.User, Content: time.Now().Add(48*time.Hour).Format("2006年01月02日") + "适合开运动会吗？",
		},
	}
	iter := runner.Run(ctx, messages)

	var lastMsg adk.Message
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			log.Println(event.Err)
		}
		msg, err := event.Output.MessageOutput.GetMessage()
		if err != nil {
			log.Println(err)
		}
		lastMsg = msg
		fmt.Println(lastMsg)
	}

	fmt.Println("[最终答案]")
	if lastMsg.Role == schema.Assistant && len(lastMsg.Content) > 0 {
		fmt.Println(lastMsg)
		messages = append(messages, lastMsg)
	}

	fmt.Println(bytes.Repeat([]byte("-"), 30))
	// 第二次询问
	//messages = append(messages, &schema.Message{Role: schema.User, Content: "运动会的前天晚上会结冰吗？"})
	messages = append(messages, &schema.Message{Role: schema.User, Content: "我在石家庄，城市编码是130100"})
	iter = runner.Run(ctx, messages)

	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			log.Println(event.Err)
		}
		msg, err := event.Output.MessageOutput.GetMessage()
		if err != nil {
			log.Println(err)
		}
		lastMsg = msg
		fmt.Println(lastMsg)
	}

	fmt.Println("[最终答案]")
	if lastMsg.Role == schema.Assistant && len(lastMsg.Content) > 0 {
		fmt.Println(lastMsg)
		messages = append(messages, lastMsg)
	}
}
