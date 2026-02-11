package main

import (
	"context"
	"hh_ai/agent/my_tool"
	"net/http"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type TimeInput struct {
	TimeZone string `json:"time_zone" jsonschema:"当地时区"`
}
type TimeOutput struct {
	Time string `json:"time" jsonschema:"当前时间"`
}

func GetCurrentTime(ctx context.Context, req *mcp.CallToolRequest, input TimeInput) (
	*mcp.CallToolResult, *TimeOutput, error) {
	if len(input.TimeZone) == 0 {
		input.TimeZone = "Asia/Shanghai"
	}
	loc, err := time.LoadLocation(input.TimeZone)
	if err != nil {
		return nil, nil, err
	}
	now := time.Now().In(loc)
	return nil, &TimeOutput{now.Format("2006-01-02 15:04:05")}, nil
}

func GetCurrentLocation(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (
	*mcp.CallToolResult, *my_tool.Location, error) {
	location, err := my_tool.GetMyLocation()
	return nil, location, err
}

func main() {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "hh_mcp",
		Version: "1.0.0",
	}, nil)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "time_tool",
		Description: "获取当前时间",
	}, GetCurrentTime)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "location_tool",
		Description: "获取当前的地理位置，包括省、城市（含城市名称和城市编码）",
	}, GetCurrentLocation)

	handler := mcp.NewStreamableHTTPHandler(
		func(*http.Request) *mcp.Server {
			return server
		}, &mcp.StreamableHTTPOptions{},
	)

	http.HandleFunc("/mcp", handler.ServeHTTP)
	if err := http.ListenAndServe("127.0.0.1:5678", nil); err != nil {
		panic(err)
	}
}
