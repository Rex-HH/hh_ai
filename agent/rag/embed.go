package rag

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"sync"
	"time"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
)

var (
	embedder     *ark.Embedder
	embedderOnce sync.Once
)

const (
	VECTOR_DIM = 2048
)

func GetEmbedder() *ark.Embedder {
	embedderOnce.Do(func() {
		ctx := context.Background()
		timeout := 3 * time.Second
		retryTimes := 3
		APITypeMultiModal := ark.APITypeMultiModal
		var err error
		embedder, err = ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
			Timeout:    &timeout,
			RetryTimes: &retryTimes,
			APIKey:     os.Getenv("ARK_API_KEY"),
			Model:      "doubao-embedding-vision-251215",
			APIType:    &APITypeMultiModal,
		})
		if err != nil {
			log.Fatal(err)
		}
	})
	return embedder
}

// Embedding 文本转向量
func Embedding(text []string) [][]float64 {
	ctx := context.Background()
	// text 不仅可以是一个词，还可以是一段文本
	embeddings, err := GetEmbedder().EmbedStrings(ctx, text)
	if err != nil {
		log.Fatal(err)
	}
	// 向量归一化
	for i, vec := range embeddings {
		embeddings[i] = NormVector(vec)
	}
	return embeddings
}

// 对向量进行归一化（即模长化为1）
func NormVector(vec []float64) []float64 {
	if len(vec) == 0 {
		return nil
	}
	sum := 0.
	for _, ele := range vec {
		sum += ele * ele
	}
	norm := math.Sqrt(sum)
	for i, ele := range vec {
		vec[i] = ele / norm
	}
	return vec
}

// 多个向量按位求平均
func AvgOfVector(vec [][]float64) ([]float64, error) {
	n := len(vec)
	if n == 0 {
		return nil, errors.New("empty vector")
	}
	if n == 1 {
		return vec[0], nil
	}
	l := len(vec[0])
	sum := make([]float64, l)
	for i := 0; i < n; i++ {
		if len(vec[i]) != l {
			return nil, fmt.Errorf("%dth vector dim not equal to first vector", i+1)
		}
		for j := 0; j < l; j++ {
			sum[j] += vec[i][j] // 按位求和
		}
	}
	for j := 0; j < l; j++ {
		sum[j] /= float64(n)
	}

	// 向量归一化
	return NormVector(sum), nil
}

func ToFloat32(vector []float64) []float32 {
	rect := make([]float32, len(vector))
	for i, ele := range vector {
		rect[i] = float32(ele)
	}
	return rect
}

// 求两个向量的内积。如果是两个归一化的向量，则内积就是余弦相似度
func InnerProduct(vec1, vec2 []float64) (float64, error) {
	if len(vec1) == 0 || len(vec1) != len(vec2) {
		return 0, fmt.Errorf("invalid vector length: %d %d", len(vec1), len(vec2))
	}
	sum := 0.
	for i, ele := range vec1 {
		sum += ele * vec2[i]
	}
	return sum, nil
}

func WrodSim() {
	words := []string{"跑步", "瑜伽", "北京", "上海"}
	embeddings := Embedding(words)
	for i := 0; i < len(words); i++ {
		for j := i + 1; j < len(words); j++ {
			dot, err := InnerProduct(embeddings[i], embeddings[j])
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("sim of %s and %s is %.4f\n", words[i], words[j], dot)
		}
	}
	fmt.Println()

	doc1, err := AvgOfVector([][]float64{embeddings[0], embeddings[1]})
	if err != nil {
		log.Fatal(err)
	}
	doc2, err := AvgOfVector([][]float64{embeddings[2], embeddings[3]})
	if err != nil {
		log.Fatal(err)
	}
	words = []string{"滑冰", "杭州"}
	vectors := Embedding(words)
	for i, vector := range vectors {
		dot, err := InnerProduct(vector, doc1)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("sim of word %s and doc1 is %.4f\n", words[i], dot)
		dot, err = InnerProduct(vector, doc2)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("sim of word %s and doc2 is %.4f\n", words[i], dot)
	}
}
