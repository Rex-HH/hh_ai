package agent

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino-ext/components/model/deepseek"
)

var (
	deepseekModel *deepseek.ChatModel
	arkModel      *ark.ChatModel
	aOnce         sync.Once
)

func CreateDeepseekModel() *deepseek.ChatModel {
	aOnce.Do(func() {
		ctx := context.Background()
		var err error
		deepseekModel, err = deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
			APIKey:  os.Getenv("DEEPSEEK_API_KEY"),
			BaseURL: "https://api.deepseek.com/beta",
			Model:   "deepseek-chat",

			// 以下是可选项
			//Timeout:     30 * time.Second,
			//MaxTokens:   3000,
			//Temperature: 1.0,
			//Stop:        []string{"的", "了"},
		})
		if err != nil {
			log.Fatal(err)
		}
	})
	return deepseekModel
}

func CreateArkModel() *ark.ChatModel {
	aOnce.Do(func() {
		ctx := context.Background()
		var err error
		//Timeout, MaxTokens, Temperature := 30*time.Second, 3000, float32(1.0)
		arkModel, err = ark.NewChatModel(ctx, &ark.ChatModelConfig{
			APIKey: os.Getenv("ARK_API_KEY"),
			Model:  "doubao-seed-1-8-251228",

			// 以下是可选项
			//Timeout:     &Timeout,
			//MaxTokens:   &MaxTokens,
			//Temperature: &Temperature,
			//Stop:        []string{"的", "了"},
		})
		if err != nil {
			log.Fatal(err)
		}
	})
	return arkModel
}
