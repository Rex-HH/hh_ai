package my_tool

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

var (
	GAODE_KEY = os.Getenv("GAODE_KEY")
)

func HttpGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http get error: %W", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http response StatusCode: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body error: %W", err)
	}
	return data, nil
}

func GetOutboundIP() string {
	data, err := HttpGet("https://httpbin.org/ip")
	if err != nil {
		log.Printf("获取对外ip失败：%s", err)
		return ""
	}
	mp := make(map[string]string, 1)
	sonic.Unmarshal(data, &mp)
	return mp["origin"]
}

type Location struct {
	Status    string `json:"status"`
	Info      string `json:"info"`
	Infocode  string `json:"infocode"`
	Province  string `json:"province"`
	City      string `json:"city"`
	Adcode    string `json:"adcode"`
	Rectangle string `json:"rectangle"`
}

func GetMyLocation() (*Location, error) {
	ip := GetOutboundIP()
	url := "https://restapi.amap.com/v3/ip?key=" + GAODE_KEY + "&ip=" + ip
	data, err := HttpGet(url)
	if err != nil {
		log.Print("http get error: %W", err)
		return nil, err
	}
	var location Location
	err = sonic.Unmarshal(data, &location)
	if err != nil {
		log.Print("unmarshal json error: %W", err)
		return nil, err
	}
	if location.Status != "1" {
		log.Printf("GetMyLocation status: %s", location.Status)
		return nil, fmt.Errorf("   %s", location.Status)
	}
	return &location, nil
}

func CreateLocationTool() tool.InvokableTool {
	locationTool, err := utils.InferTool("location_tool", "获取当前的地理位置，"+
		"包括省、城市（含城市名称和城市编码）", func(ctx context.Context, input struct{}) (string, error) {
		location, err := GetMyLocation()

		if location == nil {
			return "", err
		}
		return fmt.Sprintf("当前ip位置为：%s,%s，位置编码是%s",
			location.Province, location.City, location.Adcode), nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return locationTool
}

func UseLocationTool() {
	ctx := context.Background()
	tool := CreateLocationTool()
	resp, err := tool.InvokableRun(ctx, "{}")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
}

type LiveWeather struct {
	Province         string `json:"province"`
	City             string `json:"city"`
	Adcode           string `json:"adcode"`
	Weather          string `json:"weather"`
	Temperature      string `json:"temperature"`
	WindDirection    string `json:"winddirection"`
	WindPower        string `json:"windpower"`
	Humidity         string `json:"humidity"`
	Reporttime       string `json:"reporttime"`
	TemperatureFloat string `json:"temperature_float"`
	HumidityFloat    string `json:"humidity_float"`
}

type WeatherCast struct {
	Date string // 日期
	Week string // 周几
	// 白天
	DayWeather       string // 天气
	DayTemperature   string `json:"daytemp"`  // 温度
	DayWindDirection string `json:"daywind"`  // 风向
	DayWindPowerr    string `json:"daypower"` // 风力
	// 晚上
	NightWeather       string
	NightTemperature   string `json:"nighttemp"`
	NightWindDirection string `json:"nightwind"`
	NightWindPowerr    string `json:"nightpower"`
}
type WeatherInfo struct {
	Status   string        `json:"status"`
	Count    string        `json:"count"`
	Info     string        `json:"info"`
	Infocode string        `json:"infocode"`
	Lives    []LiveWeather `json:"lives"`
	Forcasts []struct {
		Casts []WeatherCast `json:"casts"`
	} `json:"forecasts"`
}

func WeatherForcast(GaodeCityCode string, live bool) *WeatherInfo {
	url := "https://restapi.amap.com/v3/weather/weatherInfo?key=" + GAODE_KEY +
		"&city=" + GaodeCityCode
	if !live {
		url += "&extensions=all"
	}
	resp, err := HttpGet(url)
	if err != nil {
		log.Print("http get error: %W", err)
		return nil
	}
	var weatherInfo WeatherInfo
	err = sonic.Unmarshal(resp, &weatherInfo)
	if err != nil {
		log.Print("unmarshal json error: %W", err)
		return nil
	}
	if weatherInfo.Status != "1" {
		log.Printf("WeatherForcast status: %s", weatherInfo.Status)
		return nil
	}
	return &weatherInfo
}

type WeatherParams struct {
	CityCode string `json:"city_code" jsonschema:"required,description=城市的编码"`
	Day      string `json:"day" jsonschema:"required,enum=0,enum=1,enum=2,enum=3,enum=4" 
	jsonschema_description:获取现在和未来的气象数据。
	0表示当前实时的气象数据，1表示今天的气象数据，2表示明天的气象数据，3表示后天的气象数据，4表示大后天的气象数据，
	最大只能到4"` // description中如果包含,还需要转义（,在tag中是分隔的意思），此时建议用jsonschema_description
}

func CreateWeatherTool() tool.InvokableTool {
	// 使用InferTool()创建工具。虽然创建Tool也可以定义一个struct，
	// 然后实现接口BaseTool、InvokableTool，但还是只推荐InferTool()这一种方式
	weatherTool, err := utils.InferTool("weather_tool",
		"获取气象数据，包含天气、温度、风力、风向等，既可以获得当前的实时气象数据，也可以获取未来3天的气象数据",
		func(ctx context.Context, params WeatherParams) (output string, err error) {
			day, err := strconv.Atoi(params.Day)
			if err != nil {
				return "", errors.New("参数day非法")
			}
			if day == 0 {
				// 2. 处理业务逻辑
				weather := GetCurrentWeather(params.CityCode)
				if weather != nil {
					// 3. 把结果序列化为 string 并返回
					return fmt.Sprintf("现在的天气情况是: %s，温度 %s℃，湿度 %s%%，风力 %s级，风向 %s",
						weather.Weather, weather.Temperature, weather.Humidity, weather.WindPower,
						weather.WindDirection), nil
				} else {
					return "", errors.New("获取不到天气数据")
				}
			} else if day >= 1 && day <= 4 {
				weather := GetFutureWeather(params.CityCode, day-1)
				if weather != nil {
					return fmt.Sprintf("%s 白天的天气情况是：%s，温度 %s℃，风力 %s级，风向 %s；"+
						"晚上的天气情况是：%s，温度 %s℃，风力 %s级，风向 %s",
						weather.Date, weather.DayWeather, weather.DayTemperature, weather.DayWindPowerr,
						weather.DayWindDirection, weather.NightWeather, weather.NightTemperature,
						weather.NightWindPowerr, weather.NightWindDirection), nil
				} else {
					return "", errors.New("获取不到天气数据")
				}
			} else {
				return "", errors.New("参数day非法")
			}
		})
	if err != nil {
		log.Fatal(err)
	}
	return weatherTool
}

func GetFutureWeather(code string, i int) *WeatherCast {
	weatherInfo := WeatherForcast(code, false)
	if weatherInfo == nil || weatherInfo.Forcasts == nil ||
		len(weatherInfo.Forcasts) == 0 || weatherInfo.Forcasts[0].Casts == nil ||
		len(weatherInfo.Forcasts[0].Casts) <= i {
		return nil
	}
	return &weatherInfo.Forcasts[0].Casts[i]
}

func GetCurrentWeather(code string) *LiveWeather {
	weatherInfo := WeatherForcast(code, true)
	if weatherInfo == nil || len(weatherInfo.Lives) == 0 {
		return nil
	}
	return &weatherInfo.Lives[0]
}

func UseWeatherTool(day int) {
	ctx := context.Background()
	tool := CreateWeatherTool()
	params, err := sonic.Marshal(WeatherParams{
		CityCode: "130100",
		Day:      strconv.Itoa(day),
	})
	if err != nil {
		log.Fatal(err)
	}
	resp, err := tool.InvokableRun(ctx, string(params))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
}
