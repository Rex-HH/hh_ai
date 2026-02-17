package dev_kit

import (
	"context"
	"fmt"
	"hh_ai/agent"
	"log"

	"github.com/cloudwego/eino/adk"
)

/*
ParallelAgent 的执行遵循以下设定：
并发执行：所有子 Agent 同时启动，在独立的 goroutine 中并行执行
共享输入：所有子 Agent 接受相同的处事输入和上下文
等待与结果聚合：内部使用 sync.WaitGroup 等待所有子 Agent 执行完成，收集所有子 Agent 的执行结果并按接受顺序输出
另外 Parallel 内部默认包含异常处理机制

Panic恢复：每个 goroutine 都有独立的 panic 恢复机制
错误隔离：单个子 Agent 的错误不会影响其他子 Agent 的执行
中断处理：支持子 Agent 的中断和恢复机制
*/

// 技术分析 Agent
func NewTechnicalAnalystAgent() adk.Agent {
	a, err := adk.NewChatModelAgent(context.Background(), &adk.ChatModelAgentConfig{
		Name:        "technicalAnalyst",
		Description: "从技术角度分析内容",
		Instruction: `你是一个技术专家。请从技术实现，架构设计、性能优化等技术角度分析提供的内容。
重点关注：
1.技术可行性
2.架构合理性
3.性能考量
4.技术风险
5.实现复杂度`,
		Model: agent.CreateDeepseekModel(),
	})
	if err != nil {
		log.Fatal(err)
	}
	return a
}

// 商业分析 Agent
func NewBusinessAnalystAgent() adk.Agent {
	a, err := adk.NewChatModelAgent(context.Background(), &adk.ChatModelAgentConfig{
		Name:        "businessAnalyst",
		Model:       agent.CreateDeepseekModel(),
		Description: "从商业角度分析内容",
		Instruction: `你是一个商业分析专家。请从商业价值、市场前景、成本效益等商业角度分析提供的内容。
重点关注：
1. 商业价值 
2. 市场需求
3. 竞争优势
4. 成本分析
5. 盈利模式`,
	})
	if err != nil {
		log.Fatal(err)
	}
	return a
}

// 用户体验分析 Agent
func NewUXAnalystAgent() adk.Agent {
	a, err := adk.NewChatModelAgent(context.Background(), &adk.ChatModelAgentConfig{
		Name:        "uxAnalyst",
		Description: "从用户体验角度分析内容",
		Instruction: `你是一个用户体验专家。请从用户体验、易用性、用户满意度等角度分析提供的内容。
重点关注：
1. 用户友好性
2. 操作便利性
3. 学习成本
4. 用户满意度
5. 可访问性`,
		Model: agent.CreateDeepseekModel(),
	})
	if err != nil {
		log.Fatal(err)
	}
	return a
}

// 安全分析 Agent
func NewSecurityAnalystAgent() adk.Agent {
	a, err := adk.NewChatModelAgent(context.Background(), &adk.ChatModelAgentConfig{
		Name:        "SecurityAnalyst",
		Description: "从安全角度分析内容",
		Instruction: `你是一个安全专家。请从信息安全、数据保护、隐私合规等安全角度分析提供的内容。
重点关注：
1. 数据安全
2. 隐私保护
3. 访问控制
4. 安全漏洞
5. 合规要求`,
		Model: agent.CreateDeepseekModel(),
	})
	if err != nil {
		log.Fatal(err)
	}
	return a
}

func ParallelAgent() {
	ctx := context.Background()
	techAnalyst := NewTechnicalAnalystAgent()
	bizAnalyst := NewBusinessAnalystAgent()
	uxAnalyst := NewUXAnalystAgent()
	secAnalyst := NewSecurityAnalystAgent()

	parallelAgent, err := adk.NewParallelAgent(ctx, &adk.ParallelAgentConfig{
		Name:        "MultiperspectiveAnalyzer",
		Description: "多角度并行分析：技术 + 商业 + 用户体验 + 安全",
		SubAgents:   []adk.Agent{techAnalyst, bizAnalyst, uxAnalyst, secAnalyst},
	})
	if err != nil {
		log.Fatal(err)
	}

	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent: parallelAgent,
	})
	// 要分析的产品方案
	productProposal := `
产品方案：智能客服系统

概述：开发一个基于大语言模型的智能客服系统，能够自动回答用户问题，处理常见业务咨询，并在必要时转接人工客服。

主要功能：
1. 自然语言理解和回复
2. 多轮对话管理
3. 知识库集成
4. 情感分析
5. 人工客服转接
6. 对话历史记录
7. 多渠道接入（网页、微信、APP）

技术架构：
- 前端：React + TypeScript
- 后端：Go + Gin 框架
- 数据库：PostgreSQL + Redis
- AI模型：GPT-4 API
- 部署：Docker + Kubernetes
`

	fmt.Println("开始多角度并行分析...")
	iter := runner.Query(ctx, "请分析以下产品方案：\n"+productProposal)

	results := make(map[string]string)

	for {
		event, ok := iter.Next()
		if !ok {
			break
		}

		if event.Err != nil {
			log.Printf("分析过程中出现错误：%v", event.Err)
			continue
		}

		if event.Output != nil && event.Output.MessageOutput != nil {
			results[event.AgentName] = event.Output.MessageOutput.Message.Content
			fmt.Printf("\n=== %s 分析完成 ===\n", event.AgentName)
		}

	}
	// 输出所有分析结果
	fmt.Println("\n" + "============================================================")
	fmt.Println("多角度分析结果汇总")
	fmt.Println("============================================================")
	analysisOrder := []string{"TechnicalAnalyst", "BusinessAnalyst", "UXAnalyst", "SecurityAnalyst"}
	analysisNames := map[string]string{
		"TechnicalAnalyst": "技术分析",
		"BusinessAnalyst":  "商业分析",
		"UXAnalyst":        "用户体验分析",
		"SecurityAnalyst":  "安全分析",
	}

	for _, agentName := range analysisOrder {
		if result, exists := results[agentName]; exists {
			fmt.Printf("\n【%s】\n", analysisNames[agentName])
			fmt.Printf("%s\n", result)
			fmt.Println("----------------------------------------")
		}
	}

	fmt.Println("\n多角度并行分析完成！")
	fmt.Printf("共收到 %d 个分析结果\n", len(results))
}
