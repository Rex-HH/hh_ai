package eino_graph

import (
	"context"
	"fmt"
	"hh_ai/agent"
	"log"
	"strings"

	"github.com/cloudwego/eino/compose"
)

type TA struct {
	Data string
}

type TB struct {
	Data int
}
type TC struct {
	Data float64
}
type TD struct {
	Data bool
}

// 构建图
func ComposeGraph(msg string) {
	ctx := context.Background()

	//每个node的输入、输出数据是泛型
	node1 := compose.InvokableLambda[*TA, *TB](func(ctx context.Context, input *TA) (output *TB, err error) {
		return &TB{len(input.Data)}, nil
	})
	node2 := compose.InvokableLambda[*TB, *TC](func(ctx context.Context, input *TB) (output *TC, err error) {
		return &TC{float64(input.Data)}, nil
	})
	node3 := compose.InvokableLambda[*TB, *TD](func(ctx context.Context, input *TB) (output *TD, err error) {
		return &TD{input.Data > 5}, nil
	})
	node4 := compose.InvokableLambda[*TC, *TD](func(ctx context.Context, input *TC) (output *TD, err error) {
		return &TD{input.Data <= 5}, nil
	})

	graph := compose.NewGraph[*TA, *TD]()
	// 添加node
	graph.AddLambdaNode("node1", node1)
	graph.AddLambdaNode("node2", node2)
	graph.AddLambdaNode("node3", node3)
	graph.AddLambdaNode("node4", node4)
	// 添加Edge
	graph.AddEdge(compose.START, "node1")
	graph.AddEdge("node2", "node4")
	graph.AddEdge("node3", compose.END)
	graph.AddEdge("node4", compose.END)
	// 添加Branch
	graph.AddBranch("node1", compose.NewGraphBranch[*TB](func(ctx context.Context, in *TB) (endNode string, err error) {
		if in.Data%2 == 0 {
			return "node2", nil
		} else {
			return "node3", nil
		}
	}, map[string]bool{"node2": true, "node3": true})) // 将所有可能的目标node添加到map中

	// 编译 Graph。 检查上下游节点的【input， output】 数据是否能衔接上
	runnable, err := graph.Compile(ctx)
	if err != nil {
		log.Fatal()
	}

	// 运行 Graph
	input := &TA{msg}
	result, err := runnable.Invoke(ctx, input)
	if err != nil {
		log.Fatal()
	}
	fmt.Println(result.Data)
}

// 构建图
func GraphWithCallBack(msg string) {
	ctx := context.Background()
	// callbacks.AppendGlobalHandlers(&llm.LoggerCallbacks{}) // 全局CallBack
	// callbacks.AppendGlobalHandlers(llm.GetStartCallback()) // 全局CallBack

	//每个node的输入、输出数据是泛型
	node1 := compose.InvokableLambda[*TA, *TB](func(ctx context.Context, input *TA) (output *TB, err error) {
		return &TB{len(input.Data)}, nil
	})
	node2 := compose.InvokableLambda[*TB, *TC](func(ctx context.Context, input *TB) (output *TC, err error) {
		return &TC{float64(input.Data)}, nil
	})
	node3 := compose.InvokableLambda[*TB, *TD](func(ctx context.Context, input *TB) (output *TD, err error) {
		return &TD{input.Data > 5}, nil
	})
	node4 := compose.InvokableLambda[*TC, *TD](func(ctx context.Context, input *TC) (output *TD, err error) {
		return &TD{input.Data <= 5}, nil
	})

	graph := compose.NewGraph[*TA, *TD]()
	// 添加node
	graph.AddLambdaNode("node1", node1, compose.WithNodeName("node1"))
	graph.AddLambdaNode("node2", node2, compose.WithNodeName("node2"))
	graph.AddLambdaNode("node3", node3, compose.WithNodeName("node3"))
	graph.AddLambdaNode("node4", node4, compose.WithNodeName("node4"))
	// 添加Edge
	graph.AddEdge(compose.START, "node1")
	graph.AddEdge("node2", "node4")
	graph.AddEdge("node3", compose.END)
	graph.AddEdge("node4", compose.END)
	// 添加Branch
	graph.AddBranch("node1", compose.NewGraphBranch[*TB](func(ctx context.Context, in *TB) (endNode string, err error) {
		if in.Data%2 == 0 {
			return "node2", nil
		} else {
			return "node3", nil
		}
	}, map[string]bool{"node2": true, "node3": true})) // 将所有可能的目标node添加到map中

	// 编译 Graph。 检查上下游节点的【input， output】 数据是否能衔接上
	runnable, err := graph.Compile(ctx)
	if err != nil {
		log.Fatal()
	}

	// 运行 Graph
	input := &TA{msg}
	result, err := runnable.Invoke(ctx, input, compose.WithCallbacks(agent.GetStartCallback())) // 仅针对本次运行使用的callback

	if err != nil {
		log.Fatal()
	}
	fmt.Println(result.Data)

	fmt.Println(strings.Repeat("~", 100))
	/*再次运行Graph*/
	result, err = runnable.Invoke(ctx, input)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.Data)
}

type GraphStat struct {
	Length int
}

// 构建图
func GraphWithState(msg string) {
	ctx := context.Background()

	//每个node的输入、输出数据是泛型
	node1 := compose.InvokableLambda[*TA, *TB](func(ctx context.Context, input *TA) (output *TB, err error) {
		return &TB{len(input.Data)}, nil
	})
	node2 := compose.InvokableLambda[*TB, *TC](func(ctx context.Context, input *TB) (output *TC, err error) {
		return &TC{float64(input.Data)}, nil
	})
	node3 := compose.InvokableLambda[*TB, *TD](func(ctx context.Context, input *TB) (output *TD, err error) {
		return &TD{input.Data > 5}, nil
	})
	node4 := compose.InvokableLambda[*TC, *TD](func(ctx context.Context, input *TC) (output *TD, err error) {
		return &TD{input.Data <= 5}, nil
	})

	graph := compose.NewGraph[*TA, *TD](
		// 指明图是带状态的，并且初始化状态
		compose.WithGenLocalState(func(ctx context.Context) *GraphStat {
			return &GraphStat{
				Length: 0,
			}
		}))
	// 添加node
	graph.AddLambdaNode("node1", node1,
		compose.WithStatePreHandler(NodeStatePreHandler[*TA]),
		compose.WithStatePostHandler(NodeStatePostHandler[*TB]),
	)
	graph.AddLambdaNode("node2", node2,
		compose.WithStatePreHandler(NodeStatePreHandler[*TB]),
		compose.WithStatePostHandler(NodeStatePostHandler[*TC]),
	)
	graph.AddLambdaNode("node3", node3,
		compose.WithStatePreHandler(NodeStatePreHandler[*TB]),
		compose.WithStatePostHandler(NodeStatePostHandler[*TD]),
	)
	graph.AddLambdaNode("node4", node4,
		compose.WithStatePreHandler(NodeStatePreHandler[*TC]),
		compose.WithStatePostHandler(NodeStatePostHandler[*TD]),
	)
	// 添加Edge
	graph.AddEdge(compose.START, "node1")
	// graph.AddEdge("node1","node2")
	// graph.AddEdge("node1","node3")
	graph.AddEdge("node2", "node4")
	graph.AddEdge("node3", compose.END)
	graph.AddEdge("node4", compose.END)
	// 添加Branch
	graph.AddBranch("node1", compose.NewGraphBranch[*TB](func(ctx context.Context, in *TB) (endNode string, err error) { // 根据不同的条件，node1可能会走向node2，也可能会走向node3
		if in.Data%2 == 0 {
			return "node2", nil
		} else {
			return "node3", nil
		}
	}, map[string]bool{"node2": true, "node3": true})) // 把所有可能的目标node添加到map里

	/*编译Graph。检查上下游节点的[Input,Output]数据类型是否能衔接上*/
	runnable, err := graph.Compile(ctx)
	if err != nil {
		log.Fatal(err)
	}

	/*运行Graph*/
	input := &TA{msg}
	result, err := runnable.Invoke(ctx, input)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.Data)
}

func NodeStatePreHandler[I any](ctx context.Context, in I, state *GraphStat) (I, error) {
	state.Length += 1 // 读取图状态，修改图状态
	return in, nil
}

func NodeStatePostHandler[O any](ctx context.Context, out O, state *GraphStat) (O, error) {
	log.Printf("grpath length %d", state.Length) // 读取图状态
	return out, nil
}
