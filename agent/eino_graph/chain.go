package eino_graph

import (
	"context"
	"fmt"
	"hh_ai/agent"
	"hh_ai/agent/my_tool"
	"io"
	"log"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

type HistoryMessage struct {
	History       []*schema.Message
	SystemMessage string
	UserMessage   string
}

// 用Graph来表示Chain。（chain的弊端：ChatModel和Tool的调用次数都是写死的）
func GraphOfChain() {
	ctx := context.Background()
	//callbacks.AppendGlobalHandlers(agent.GetChatModelInputCallback())

	//创建ChatModel
	chatModel := agent.CreateDeepseekModel()
	//创建工具集
	tools := []tool.BaseTool{
		my_tool.CreateWeatherTool(),
		my_tool.CreateLocationTool(),
	}
	//创建ToolsNode（该Node包含多个Tool）
	ToolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{Tools: tools})
	if err != nil {
		log.Fatal(err)
	}
	// 把Tool信息绑定到ChatModel（让Model理解每一个Tool的功能及使用方式）
	toolInfos := make([]*schema.ToolInfo, 0, len(tools))
	for _, tool := range tools {
		if info, err := tool.Info(ctx); err == nil {
			toolInfos = append(toolInfos, info)
		}
	}
	// BindTools 模型会自主决定直接给答案，还是调用一个/多个工具
	// BindForcedTools 模型必须调用一个/多个工具
	chatModel.BindTools(toolInfos)

	addHistory := compose.InvokableLambda[[]*schema.Message, []*schema.Message](
		func(ctx context.Context, input []*schema.Message) (output []*schema.Message, err error) {
			return
		},
	)

	// 构建Graph（其实compose可以直接NewChain）
	graph := compose.NewGraph[[]*schema.Message, *schema.Message](
		// 指明图是带状态的，并且初始化状态
		compose.WithGenLocalState(func(ctx context.Context) (state *HistoryMessage) {
			return &HistoryMessage{
				History:       make([]*schema.Message, 0, 4),
				UserMessage:   "",
				SystemMessage: "",
			}
		}),
	) // graph的输入[]*schema.Message，输出*schema.Message
	// 添加Node
	graph.AddChatModelNode("model1", chatModel,
		compose.WithStatePreHandler(func(ctx context.Context, input []*schema.Message, state *HistoryMessage) ([]*schema.Message, error) {
			// 给图状态里的SystemMessage和UserMessage赋值
			for _, msg := range input {
				if msg.Role == schema.System {
					state.SystemMessage = msg.Content
				} else if msg.Role == schema.User {
					state.UserMessage = msg.Content
				}
			}
			return input, nil // input 原封不动地返回（当然你也可以修改input）
		}),
		// 把ChatModel的输出加入HistoryMessage
		compose.WithStatePostHandler(func(ctx context.Context, output *schema.Message, state *HistoryMessage) (*schema.Message, error) {
			state.History = append(state.History, output)
			return output, nil // output 原封不动地返回（当然你也可以修改output）
		}),
	) // ChatModelNode输入[]*schema.Message，输出*schema.Message
	graph.AddToolsNode("tools", ToolsNode,
		// 把Tools的输出加入HistoryMessage
		compose.WithStatePostHandler(func(ctx context.Context, out []*schema.Message, state *HistoryMessage) ([]*schema.Message, error) {
			state.History = append(state.History, out...) //模型可能会同时调用toolsNode里的多个工具，所以toolsNode的输出是[]*schema.Message
			return out, nil
		}),
	)
	graph.AddLambdaNode("history", addHistory,
		compose.WithStatePostHandler(func(ctx context.Context, out []*schema.Message, state *HistoryMessage) ([]*schema.Message, error) {
			result := []*schema.Message{&schema.Message{Role: schema.System, Content: state.SystemMessage}}
			result = append(result, state.History...)
			result = append(result, &schema.Message{Role: schema.User, Content: state.UserMessage})
			return result, nil
		},
		))
	graph.AddChatModelNode("model2", chatModel) // ChatModelNode输入[]*schema.Message，输出*schema.Message
	// 添加Edge
	graph.AddEdge(compose.START, "model1")
	graph.AddEdge("model1", "tools")
	graph.AddEdge("tools", "history")
	graph.AddEdge("history", "model2")
	graph.AddEdge("model2", compose.END)

	/*编译Graph。检查上下游节点的[Input,Output]数据类型是否能衔接上*/
	runnable, err := graph.Compile(ctx)
	if err != nil {
		log.Fatal(err)
	}

	/*运行Graph*/
	// 非流式（第一次运行）
	input := []*schema.Message{
		&schema.Message{Role: schema.User, Content: "我在河北吗？"},
	}
	msg, err := runnable.Invoke(ctx, input)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("【非流式输出】")
	fmt.Println(msg.Content)

	// 流式（第二次运行）
	fmt.Println("【流式输出】")
	sr, err := runnable.Stream(ctx, input)
	if err != nil {
		log.Fatal(err)
	}
	defer sr.Close()
	for {
		msg, err := sr.Recv()
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			break
		}
		fmt.Print(msg.Content)
	}
}
