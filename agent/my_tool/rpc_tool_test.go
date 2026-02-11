package my_tool

import (
	"encoding/json"
	"testing"

	"github.com/bytedance/sonic"
)

func TestUseLocationTool(t *testing.T) {
	UseLocationTool()
}

func TestLocationParsing(t *testing.T) {
	// 模拟高德返回的数据
	jsonData := []byte(`{
		"status": "1",
		"province": "河北省",
		"city": "石家庄市",
		"adcode": "130100",
		"rectangle": "114.344832,37.908511;114.734131,38.185158"
	}`)

	var loc Location
	err := sonic.Unmarshal(jsonData, &loc)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	if loc.Province != "河北省" {
		t.Errorf("Province解析错误: 期望='河北省', 实际='%s'", loc.Province)
	}
	if loc.City != "石家庄市" {
		t.Errorf("City解析错误: 期望='石家庄市', 实际='%s'", loc.City)
	}
	if loc.Adcode != "130100" {
		t.Errorf("Adcode解析错误: 期望='130100', 实际='%s'", loc.Adcode)
	}

	t.Logf("解析成功: Province='%s', City='%s', Adcode='%s'",
		loc.Province, loc.City, loc.Adcode)
}

func TestJsonLocationParsing(t *testing.T) {
	// 模拟高德返回的数据
	jsonData := []byte(`{
		"status": "1",
		"province": "河北省",
		"city": "石家庄市",
		"adcode": "130100",
		"rectangle": "114.344832,37.908511;114.734131,38.185158"
	}`)

	var loc Location
	err := json.Unmarshal(jsonData, &loc)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	if loc.Province != "河北省" {
		t.Errorf("Province解析错误: 期望='河北省', 实际='%s'", loc.Province)
	}
	if loc.City != "石家庄市" {
		t.Errorf("City解析错误: 期望='石家庄市', 实际='%s'", loc.City)
	}
	if loc.Adcode != "130100" {
		t.Errorf("Adcode解析错误: 期望='130100', 实际='%s'", loc.Adcode)
	}

	t.Logf("解析成功: Province='%s', City='%s', Adcode='%s'",
		loc.Province, loc.City, loc.Adcode)
}

func TestDifferentCases(t *testing.T) {
	tests := []struct {
		name     string
		jsonStr  string
		expected string
	}{
		{
			name:     "精确匹配",
			jsonStr:  `{"Province":"河北"}`,
			expected: "河北",
		},
		{
			name:     "小写匹配",
			jsonStr:  `{"province":"河北"}`,
			expected: "河北",
		},
		{
			name:     "首字母大写",
			jsonStr:  `{"Province":"河北"}`,
			expected: "河北",
		},
		{
			name:     "全部大写",
			jsonStr:  `{"PROVINCE":"河北"}`,
			expected: "河北",
		},
		{
			name:     "大小写混合",
			jsonStr:  `{"PrOvInCe":"河北"}`,
			expected: "河北",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var loc Location
			json.Unmarshal([]byte(tt.jsonStr), &loc)
			if loc.Province != tt.expected {
				t.Errorf("%s: 期望='%s', 实际='%s'", tt.name, tt.expected, loc.Province)
			}
		})
	}
}

func TestUseWeatherTool(t *testing.T) {
	UseWeatherTool(0)
	UseWeatherTool(1)
	UseWeatherTool(4)
	UseWeatherTool(5) // 数day非法

}
