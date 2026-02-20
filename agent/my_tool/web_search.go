package my_tool

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/karust/openserp/core"
	"github.com/karust/openserp/duckduckgo" // duckduckgo 官方并没有维护任何search API
)

var (
	browser    *core.Browser
	browseOnce sync.Once
)

func GetBrowser() *core.Browser {
	browseOnce.Do(func() {
		var err error
		// 打开电脑上的默认浏览器
		browser, err = core.NewBrowser(core.BrowserOpts{})
		if err != nil {
			panic(err)
		}
	})
	return browser
}

func CloseBrowser() {
	// 关闭浏览器时增加 panic 保护，避免底层 rod 在重复关闭或部分初始化失败时触发 nil pointer panic
	if browser != nil {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("recover from browser.Close panic: %v", r)
			}
		}()
		browser.Close()
	}
}

type SearchRequest struct {
	Query string `json:"query" jsonschema:"required,description=搜索词"`
	Count int    `json:"count" jsonschema:"description=最多返回搜索结果的条数，默认值为3"`
}

// 使用 DuckDuckGo 不需要 apikey 鉴权，而其他搜索引擎都需要
func DuckDuckSearch(ctx context.Context, request *SearchRequest) (string, error) {
	if request.Query == "" {
		return "", errors.New("搜索query为空")
	}
	if request.Count == 0 {
		request.Count = 3
	}
	if request.Count > 3 {
		request.Count = 3
	}
	query := core.Query{
		Text:     request.Query, // 搜索词
		LangCode: "zh-CN",       // 中文
		Limit:    request.Count, // 限制结果的条数
		Site:     "",            // 限定网站，如github.com, stackoverflow.com
		/*
			# SOCKS5 proxy
			./openserp serve --proxy socks5://127.0.0.1:1080
			# HTTP proxy with authentication
			./openserp search bing "query" --proxy http://user:pass@127.0.0.1:8080
		*/
		ProxyURL: "",
		Insecure: false, // 是否允许使用TLS连接
	}
	ddg := duckduckgo.New(*GetBrowser(), core.SearchEngineOptions{})
	results, err := ddg.Search(query)
	if err != nil {
		return "", err
	}
	s, err := sonic.MarshalString(results)
	if err != nil {
		return "", err
	}
	return s, nil
}

func CreateDuckDuckSearchTool() tool.InvokableTool {
	// 使用 InferTools 创建工具
	tool, err := utils.InferTool("web_search_tool", "调用搜索引擎，获取搜索结果", DuckDuckSearch)
	if err != nil {
		log.Fatal(err)
	}
	return tool
}

func useDuckDuckSearch(query string) {
	ctx := context.Background()
	// 创建工具
	searchTool := CreateDuckDuckSearchTool()

	// 构造参数
	request := SearchRequest{
		Query: query,
		Count: 2,
	}
	jsonReq, err := sonic.MarshalString(request)
	if err != nil {
		log.Fatalf("Marshal of search request failed: %s", err)
	}

	// 调用工具
	resp, err := searchTool.InvokableRun(ctx, jsonReq)
	if err != nil {
		log.Fatal(err)
	}

	// 输出结果
	fmt.Println(resp)
}

type Url struct {
	Url string `json:"url" jsonschema:"required,description=被爬取的网页的url"`
}

func Crawl(ctx context.Context, url *Url) (string, error) {
	if url.Url == "" {
		return "", errors.New("没有指定爬取的URL")
	}
	page, err := GetBrowser().Navigate(url.Url)
	if err != nil {
		return "", err
	}
	element, _ := page.Element("title")
	title, _ := element.Text()
	element, _ = page.Element("body")
	body, _ := element.Text()
	result := "网页标题:" + title + "\n网页内容：" + body
	return result, nil
}

func CreateCrawlTool() tool.InvokableTool {
	tool, err := utils.InferTool("web_crawler_tool", "网页爬虫", Crawl)
	if err != nil {
		log.Fatal(err)
	}
	return tool
}
func UseCrawlTool(url string) {
	ctx := context.Background()
	CrawlTool := CreateCrawlTool()
	request := Url{
		Url: url,
	}
	marshalString, err := sonic.MarshalString(request)
	if err != nil {
		log.Fatalf("Marshal of crawl request failed: %s", err)
	}
	resp, err := CrawlTool.InvokableRun(ctx, marshalString)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
}
