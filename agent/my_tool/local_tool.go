package my_tool

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

type TimeRequest struct {
	TimeZone string `json:"time_zone" jsonschema:"description=当地时区"`
}

func GetTime(ctx context.Context, timeZone *TimeRequest) (string, error) {
	if len(timeZone.TimeZone) == 0 {
		timeZone.TimeZone = "Asia/Shanghai"
	}
	loc, err := time.LoadLocation(timeZone.TimeZone)
	if err != nil {
		return "", err
	}
	now := time.Now().In(loc)
	return now.Format("2006-01-02 15:04:05"), nil
}

func CreateTimeTool() tool.InvokableTool {
	timeTool, err := utils.InferTool("time_tool", "获取当前时间", GetTime)
	if err != nil {
		log.Fatal(err)
	}
	return timeTool
}

func UseTimeTool(TimeZone string) {
	ctx := context.Background()
	timeTool := CreateTimeTool()
	timeRequest := TimeRequest{
		TimeZone: TimeZone,
	}
	jsonReq, err := json.Marshal(timeRequest)
	if err != nil {
		log.Fatalf("Marshal of time request failed: %s", err)
	}

	resp, err := timeTool.InvokableRun(ctx, string(jsonReq))
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
