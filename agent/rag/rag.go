package rag

import (
	"context"
	"fmt"
	"hh_ai/agent"
	"log"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino-ext/components/document/loader/file"
	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	qdrant_retriever "github.com/cloudwego/eino-ext/components/retriever/qdrant"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/components/document/parser"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	"github.com/qdrant/go-client/qdrant"
)

type Doc struct {
	ID      string    `json:"id,omitempty"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Vector  []float64 `json:"vector,omitempty"`
}

const (
	QaCollection = "qa"
)

var (
	docMap = make(map[string]*Doc, 300) // 实际中，这种数据应该存储在关系型数据库中，通过DocID跟向量数据库进行关联
)

func IndexDocument(markdownFile string) {
	ctx := context.Background()

	// 加载文档
	textParser := parser.TextParser{}
	loader, _ := file.NewFileLoader(ctx, &file.FileLoaderConfig{
		Parser: textParser,
	})
	docs, err := loader.Load(ctx, document.Source{URI: markdownFile})
	if err != nil {
		log.Fatal(err)
	}
	// 切分文档

	splitTransformer, _ := markdown.NewHeaderSplitter(ctx, &markdown.HeaderConfig{
		Headers: map[string]string{
			"#":  "",
			"##": "",
		},
	})
	transformedDocs, _ := splitTransformer.Transform(ctx, docs)

	points := make([]*qdrant.PointStruct, 0, len(transformedDocs))

	// 计算每一个片段的向量
	for _, doc := range transformedDocs {
		DocId := uuid.NewString()
		var title, content string
		line := strings.SplitN(doc.Content, "\n", 2)
		if strings.HasPrefix(line[0], "## ") {
			title = line[0][3:]
			content = line[1]
		} else {
			continue
		}

		vectors := Embedding([]string{title, title, content})
		vector, _ := AvgOfVector(vectors) // 标题向量权重为2， 正文向量权重为1
		docMap[DocId] = &Doc{
			ID:      DocId,
			Title:   title,
			Content: content,
			Vector:  vector,
		}
		points = append(points, &qdrant.PointStruct{
			Id: &qdrant.PointId{
				PointIdOptions: &qdrant.PointId_Uuid{Uuid: DocId},
			},
			Vectors: &qdrant.Vectors{
				VectorsOptions: &qdrant.Vectors_Vector{
					Vector: &qdrant.Vector{
						Vector: &qdrant.Vector_Dense{
							Dense: &qdrant.DenseVector{
								Data: ToFloat32(vector),
							},
						},
					},
				},
			},
		})
	}

	GetQdrantClient().DeleteCollection(ctx, QaCollection) // 出于测试目的，清空之前的数据
	// 使用 Upsert 之前需要保证Collection 已经存在
	if exists, err := GetQdrantClient().CollectionExists(ctx, QaCollection); err != nil {
		log.Fatal(err)
	} else {
		if !exists {
			err = GetQdrantClient().CreateCollection(ctx, &qdrant.CreateCollection{
				CollectionName: QaCollection,
				VectorsConfig: &qdrant.VectorsConfig{
					Config: &qdrant.VectorsConfig_Params{
						Params: &qdrant.VectorParams{
							Size:     VECTOR_DIM,
							Distance: qdrant.Distance_Cosine,
						},
					},
				},
			})
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// 把文档片段写入向量数据库
	result, err := GetQdrantClient().Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: QaCollection,
		Points:         points,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("upsert status %d", result.Status)
}

func RetrieveDocument(query string, limit int) []*Doc {
	ctx := context.Background()

	//创建Retriever
	ScoreThreshold := 0.5
	retriever, _ := qdrant_retriever.NewRetriever(ctx, &qdrant_retriever.Config{
		Client:         GetQdrantClient(),
		Collection:     QaCollection,
		Embedding:      GetEmbedder(),
		ScoreThreshold: &ScoreThreshold,
		TopK:           limit,
	})
	//执行检索
	neighbors, _ := retriever.Retrieve(ctx, query)
	log.Printf("retrieve %d docs", len(neighbors))

	result := make([]*Doc, 0, limit)
	for _, doc := range neighbors {
		id := doc.ID
		fmt.Println(id, doc.Score())
		if v, exists := docMap[id]; exists {
			result = append(result, &Doc{
				Title:   v.Title,
				Content: v.Content,
			})
		}
	}
	return result
}

func ChatBot(question string) string {
	ctx := context.Background()

	docs := RetrieveDocument(question, 4)
	knowledge := ""
	for _, doc := range docs {

		s, _ := sonic.MarshalString(doc)
		knowledge += s + "\n"
	}
	fmt.Println("[资料]", knowledge)

	// 创建 ChatModel
	model := agent.CreateDeepseekModel()
	// 创建 Agent
	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "custormer_service",
		Description: "客服机器人",
		Instruction: "请根据我提供的资料，回答用户的问题",
		Model:       model,
	})
	if err != nil {
		log.Fatal(err)
	}

	// 运行 Agent
	runner := adk.NewRunner(ctx, adk.RunnerConfig{Agent: agent})
	iter := runner.Query(ctx, fmt.Sprintf("资料如下：%s\n 用户问题：%s", knowledge, question))

	var lastMsg adk.Message
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			log.Fatal(event.Err)
		}
		mes, err := event.Output.MessageOutput.GetMessage()
		if err != nil {
			log.Fatal(err)
		}
		lastMsg = mes
	}

	if lastMsg.Role == schema.Assistant && len(lastMsg.Content) > 0 {
		return lastMsg.Content
	} else {
		return "对不起，这个问题我还在学习中。。。"
	}
}
