package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type RequestBody struct {
	Model    string     `json:"model"`
	Messages []*Message `json:"messages"`
	Stream   bool       `json:"stream"`
}
type Choice struct {
	Index   int     `json:"index"`
	Message Message `json:"message"`
}

type ResponseBody struct {
	Choices []Choice `json:"choices"`
}
type StreamChoice struct {
	Index   int     `json:"index"`
	Message Message `json:"delta"`
}
type StreamResponseBody struct {
	Choices []StreamChoice `json:"choices"`
}

const (
	DeepseekUrl = "https://api.deepseek.com/chat/completions"
)

func ChatWithDeepseek(rb *RequestBody) {
	bs, err := json.Marshal(rb)
	if err != nil {
		log.Println("json 序列化失败： %s", err)
		return
	}
	request, _ := http.NewRequest("POST", DeepseekUrl, bytes.NewReader(bs))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "Bearer "+os.Getenv("DEEPSEEK_API_KEY"))
	client := http.Client{
		Timeout: time.Second * 10,
	}
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
	}
	if response.StatusCode != http.StatusOK {
		io.Copy(os.Stdout, response.Body)
	}

	// 打印非流式
	//bs, _ = io.ReadAll(response.Body)
	////fmt.Println(string(bs))
	//var respBody ResponseBody
	//err = json.Unmarshal(bs, &respBody)
	//if err == nil {
	//	fmt.Println(respBody.Choices[0].Message.Role, respBody.Choices[0].Message.Content)
	//} else {
	//	fmt.Println(string(bs))
	//}

	// 打印流式
	buffer := make([]byte, 2048)
	for {
		n, err := response.Body.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("read response body failed: %s", err)
			break
		}
		//fmt.Println(string(buffer[:n]))
		for _, segment := range bytes.Split(buffer[:n], []byte("data:")) {
			segment = bytes.TrimSpace(segment)
			if len(segment) == 0 {
				continue
			}
			if string(segment) == "[DONE]" {
				break
			}
			var body StreamResponseBody
			err := json.Unmarshal(segment, &body)
			if err == nil {
				fmt.Print(body.Choices[0].Message.Content)
			} else {
				log.Printf("json反序列化失败: %s. [%s]", err, string(segment))
				break
			}

		}
	}

}

func main() {
	rb := &RequestBody{
		Model: "deepseek-chat",
		Messages: []*Message{
			&Message{Role: "system", Content: ""},
			&Message{Role: "user", Content: "今天是" + time.Now().Format("2006年01月02日") + "离春节还有几天？"},
		},
		Stream: true,
	}
	ChatWithDeepseek(rb)

}
