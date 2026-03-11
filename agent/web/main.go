package main

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
	"github.com/rs/xid"
)

var (
	session sync.Map
)

func ParseUrlParams(rawQuery string) map[string]string {
	params := make(map[string]string, 10)
	args := strings.Split(rawQuery, "&")
	for _, ele := range args {
		arr := strings.Split(ele, "=")
		if len(arr) == 2 {
			key, _ := url.QueryUnescape(arr[0]) //url参数反转义
			value, _ := url.QueryUnescape(arr[1])
			params[key] = value
		}
	}
	return params
}

func Chat(w http.ResponseWriter, r *http.Request) {
	log.Println("call Chat")
	ctx := context.Background()
	w.Header().Add("Content-Type", "text/event-stream; charset=utf-8") //标识响应为事件流。charset=utf-8是为了解决中文乱码
	w.Header().Add("Cache-Control", "no-cache")                        //防止浏览器缓存响应，确保实时性
	w.Header().Add("Connection", "keep-alive")                         //保持连接开放，支持持续流式传输

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	params := ParseUrlParams(r.URL.RawQuery)
	messages := make([]adk.Message, 0, 10)
	sid := params["session"]
	if v, exists := session.Load(sid); exists {
		history := v.([]adk.Message)
		messages = append(messages, history...)
	}
	messages = append(messages, &schema.Message{Role: schema.User, Content: params["msg"]})
	iter := GetRunner().Run(ctx, messages)
	answer := strings.Builder{}
	for {
		event, ok := iter.Next()
		if !ok {
			fmt.Fprint(w, "data: [DONE]\n\n") // 结束标志
			flusher.Flush()
			//log.Println("AI助手输出完毕")
			messages = append(messages, &schema.Message{Role: schema.Assistant, Content: answer.String()}) // 把LLM的输出加追加到历史对话里去
			session.Store(sid, messages)
			break
		}
		if event.Err != nil {
			log.Printf("read llm response failed: %s", event.Err)
			break
		}

		// 非流方式
		//msg, err := event.Output.MessageOutput.GetMessage() // GetMessage() 会把流式输出的结果全部合并到一起再返回。这样就没有流式的效果了
		//if err != nil {
		//	log.Printf("get output message failed: %s", event.Err)
		//} else {
		//	if msg != nil && msg.Role == schema.Assistant {
		//		log.Print(msg.Content)
		//		answer.WriteString(msg.Content)
		//		fmt.Fprintf(w, "data: %s\n\n", strings.ReplaceAll(msg.Content, "\n", "<br>")) // SSE 协议要求 数据内部不能包含换行符。此处把\n替换为<br>，在前端代码里还需要把<br>再替换回\n
		//		flusher.Flush()                                                               // 强制数据立刻发给对方
		//	}
		//}

		// 流式
		s := event.Output.MessageOutput.MessageStream
		for {
			msg, err := s.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Printf("get output message failed: %s", event.Err)
				break
			}
			if msg != nil { // 无法从一个片断中取到 msg.Role
				fmt.Print(msg.Content)
				answer.WriteString(msg.Content)
				fmt.Fprintf(w, "data: %s\n\n", strings.ReplaceAll(msg.Content, "\n", "<br>")) // SSE 协议要求 数据内部不能包含换行符。此处把\n替换为<br>，在前端代码里还需要把<br>再替换回\n
				flusher.Flush()                                                               // 强制数据立刻发给对方
			}
		}
	}
}
func main() {
	mux := http.NewServeMux()
	// 访问根路径时，返回静态页面
	mux.HandleFunc("GET /", func(writer http.ResponseWriter, request *http.Request) {
		tmpl, err := template.ParseFiles("./web/chat.html") //相对于执行go run的路径
		if err != nil {
			fmt.Println("create template failed:", err)
			return
		}
		sid := xid.New().String()
		err = tmpl.Execute(writer, map[string]string{"url": "http://127.0.0.1:5678/chat", "session": sid})
		if err != nil {
			log.Fatal(err)
		}
	})
	// 设置浏览器标签页图标
	mux.HandleFunc("GET /favicon.ico", func(writer http.ResponseWriter, request *http.Request) {
		file, err := os.Open("./web/hh.png")
		if err != nil {
			return
		}
		defer file.Close()
		io.Copy(writer, file)
	})
	// 访问 /chat 接口时，返回大模型的 completion
	mux.HandleFunc("GET /chat", Chat) // SSE不支持发起POST请求
	log.Println("going to start web server on port 5678")
	if err := http.ListenAndServe("127.0.0.1:5678", mux); err != nil {
		panic(err)
	}
}
