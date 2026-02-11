package my_tool

import (
	"context"
	"fmt"
	"log"

	"github.com/bytedance/sonic"
	eino_mcp "github.com/cloudwego/eino-ext/components/tool/mcp"
	"github.com/cloudwego/eino/components/tool"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

var (
	trainCli *client.Client
)

const (
	TRAIN_MCP_URL  = "https://mcp.api-inference.modelscope.net/f214351e6b124c/mcp"
	TRAIN_MCP_NAME = "12306-mcp"
)

func CreateMcpClient(url, name string) *client.Client {
	// 创建 MCP Client
	cli, err := client.NewStreamableHttpClient(url)
	if err != nil {
		log.Fatal("call tool err", err)
	}
	initReqest := mcp.InitializeRequest{}
	initReqest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initReqest.Params.ClientInfo = mcp.Implementation{
		Name:    name,
		Version: "1.0.0",
	}
	_, err = cli.Initialize(context.Background(), initReqest)
	if err != nil {
		log.Fatal("call tool err", err)
	}
	return cli
}

func InitMcpClient() {
	trainCli = CreateMcpClient(TRAIN_MCP_URL, TRAIN_MCP_NAME)
}

func CreateMcp12306Tool() []tool.BaseTool {
	ctx := context.Background()
	tools, err := eino_mcp.GetTools(ctx, &eino_mcp.Config{
		Cli:          trainCli,
		ToolNameList: []string{},
	})
	if err != nil {
		log.Fatal("call tool err", err)
	}
	return tools
}

func List12306Tool() {
	ctx := context.Background()
	tools := CreateMcp12306Tool()
	for _, tool := range tools {
		info, _ := tool.Info(ctx)
		fmt.Println(info.Name, info.Desc)
	}
}

type GetTicketsRequest struct {
	Date        string `json:"date"`
	FromStation string `json:"fromStation"`
	ToStation   string `json:"toStation"`
}

func Use12306Tool() {
	toolName := "get-tickets"
	ctx := context.Background()
	tools, err := eino_mcp.GetTools(ctx, &eino_mcp.Config{
		Cli:          trainCli,
		ToolNameList: []string{toolName},
	})
	if err != nil {
		log.Fatal("call tool err", err)
	}
	ticketTool := tools[0].(tool.InvokableTool)
	getTicketsRequest := GetTicketsRequest{
		Date:        "2026-02-09",
		FromStation: "SJP",
		ToStation:   "VAB",
	}
	ticketRequestJson, err := sonic.Marshal(getTicketsRequest)
	if err != nil {
		log.Fatal("call tool err", err)
	}
	run, err := ticketTool.InvokableRun(ctx, string(ticketRequestJson))
	if err != nil {
		log.Fatal("call tool err", err)
	}
	fmt.Println(run)
}
