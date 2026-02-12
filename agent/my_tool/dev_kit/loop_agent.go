package dev_kit

import (
	"context"
	"fmt"
	"hh_ai/agent"
	"log"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

/*
LoopAgent 的执行遵循以下设定：

循环执行：重复执行 SubAgents 序列 （序列是“有序的”）
History 累积：每次迭代的结果都会累积到 History 中，后续迭代可以访问所有历史信息
条件退出：达到最大迭代次数来终止循环，配置 MaxIterations=0 时表示无限循环
*/

// 代码分析 Agent
func NewCodeAnalyzerAgent() adk.Agent {
	a, err := adk.NewChatModelAgent(context.Background(), &adk.ChatModelAgentConfig{
		Name:        "CodeAnalyzer",
		Description: "分析代码质量和性能问题",
		Instruction: `你是一个代码分析专家。请分析提供的代码（如果之前给出过优化后的代码，请直接分析最后一次优化后的代码），是否存在以下5种问题：
1. 性能瓶颈
2. 代码重复
3. 可读性问题
4. 潜在的 bug
5. 不符合最佳实践的地方

如果存在问题，请指明具体是什么问题，不需要给出优化建议；否则请直接调用exit tool。注意，如果代码仍然存在需要优化的问题，不要调用exit tool`, // Instruction是SystemMessage
		Model: agent.CreateDeepseekModel(),
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []tool.BaseTool{adk.ExitTool{}}, // ExitTool内部会创建一个ExitAction
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	return a
}

// 代码优化 Agent
func NewCodeOptimizerAgent() adk.Agent {
	a, err := adk.NewChatModelAgent(context.Background(), &adk.ChatModelAgentConfig{
		Name:        "CodeOptimizer",
		Description: "根据分析结果优化代码",
		Instruction: `基于前面的代码分析结果，对代码进行优化改进：
1. 修复识别出的性能问题
2. 消除代码重复
3. 提高代码可读性
4. 修复潜在 bug
5. 应用最佳实践

请提供优化后的完整代码。`,
		Model: agent.CreateDeepseekModel(),
	})
	if err != nil {
		log.Fatal(err)
	}
	return a
}

func LoopAgent(code string) {
	ctx := context.Background()

	//创建 LoopAgent，最多执行 3 轮优化
	loopAgent, err := adk.NewLoopAgent(
		ctx, &adk.LoopAgentConfig{
			Name:        "CodeOptimizationLoop",
			Description: "代码优化循环：分析问题 -> 代码优化",
			SubAgents: []adk.Agent{
				NewCodeAnalyzerAgent(),
				NewCodeOptimizerAgent(),
			},
			MaxIterations: 3,
		})
	if err != nil {
		log.Fatal(err)
	}

	//创建Runner
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent: loopAgent,
	})

	fmt.Println("开始代码优化循环...")
	iter := runner.Query(ctx, "请优化以下 GO 代码: \n"+code)

	iteration := 1
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			log.Fatal(event.Err)
		}

		if event.Output != nil && event.Output.MessageOutput != nil {
			msg, _ := event.Output.MessageOutput.GetMessage()
			if msg.Role == schema.Assistant {
				fmt.Printf("\n================== 第 %d 轮 SubAgent: %s ==================\n", iteration, event.AgentName)
				fmt.Println(msg.Role, msg.Content)
			}

			if event.AgentName == "CodeAnalyzer" {
				if event.Action != nil && event.Action.Exit {
					fmt.Println("\n优化循环提前结束！")
					break
				}
			}

			if event.AgentName == "CodeOptimizer" {
				if msg.Role == schema.Assistant {
					iteration++
				}
			}
		}
	}

	fmt.Println("\n代码优化循环执行完成！")
}
