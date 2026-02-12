package dev_kit

import (
	"bufio"
	"context"
	"fmt"
	"hh_ai/agent"
	checkpointstore "hh_ai/agent/checkpoint_store"
	"hh_ai/agent/my_tool"
	"log"
	"os"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/rs/xid"
)

func InterruptResume() {
	ctx := context.Background()
	model := agent.CreateDeepseekModel()

	tools := []tool.BaseTool{
		//my_tool.CreateLocationTool(),
		my_tool.CreateWeatherTool(),
		my_tool.CreateTimeTool(),
		my_tool.NewAskForClarificationTool(),
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
		Agent:           agent,
		CheckPointStore: checkpointstore.NewInMemoryStore(),
	})

	messages := []adk.Message{
		&schema.Message{
			//Role: schema.User, Content: time.Now().Add(48*time.Hour).Format("2006年01月02日") + "适合开运动会吗？",
			Role: schema.User, Content: time.Now().Add(12*time.Hour).
				Format("2006年01月02日03时04分") + "适合骑电动车带80岁老人，去银行取钱吗？",
		},
	}
	checkPointId := xid.New().String()
	iter := runner.Run(ctx, messages,
		adk.WithCheckPointID(checkPointId),
	)

	var lastMsg adk.Message
LB:
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

		if len(msg.ToolCalls) > 0 {
			for _, toolCall := range msg.ToolCalls {
				if toolCall.Function.Name == my_tool.AskForClarificationToolName {
					var req my_tool.AskForClarificationInput
					if err := sonic.UnmarshalString(toolCall.Function.Arguments, &req); err == nil {
						fmt.Println(req.Question)
						scanner := bufio.NewScanner(os.Stdin)
						fmt.Print("请在此输入答案: ")
						scanner.Scan()
						fmt.Println()
						input := scanner.Text()

						iter, err = runner.Resume(ctx, checkPointId,
							adk.WithToolOptions([]tool.Option{my_tool.WithNewInput(input)}))
						if err != nil {
							log.Fatal(err)
						}
						continue LB
					}
				}
			}
		}

		lastMsg = msg
		fmt.Println(lastMsg.Role, lastMsg)
	}

	fmt.Println("[最终答案]")
	if lastMsg.Role == schema.Assistant && len(lastMsg.Content) > 0 {
		fmt.Println(lastMsg)
		messages = append(messages, lastMsg)
	}
}
