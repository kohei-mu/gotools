package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configfile string

type WeatherJson struct {
	Lat      float64   `json: "lat"`
	Lon      float64   `json: "lon"`
	Status   string    `json: "status"`
	Temp     float64   `json: "temp"`
	Pressure float64   `json: "pressure"`
	Humidity float64   `json: "humidity"`
	Wind     float64   `json: "wind"`
	Time     time.Time `json: "time"`
	Country  string    `json: "country"`
	City     string    `json: "city"`
}

var weatherCmd = &cobra.Command{
	Use:   "weather [COUNTRY]",
	Args:  cobra.MinimumNArgs(1),
	Short: "weather command",
	Long:  "weather command",
	Run: func(cmd *cobra.Command, args []string) {
		for _, country := range args {
			w := weatherReq(country)
			switch globalFlags.language {
			case "ja":
				fmt.Println("##### JSON解析結果 #####")
				fmt.Println("時刻", w.Time)
				fmt.Println("国", w.Country)
				fmt.Println("都市", w.City)
				fmt.Println("緯度", w.Lat)
				fmt.Println("経度", w.Lon)
				fmt.Println("天候", w.Status)
				fmt.Println("温度", w.Temp)
				fmt.Println("気圧", w.Pressure)
				fmt.Println("湿度", w.Humidity)
				fmt.Println("風速", w.Wind)
			case "en":
				fmt.Println("##### JSON PARSE RESULT #####")
				fmt.Println("Time", w.Time)
				fmt.Println("Country", w.Country)
				fmt.Println("City", w.City)
				fmt.Println("Latitude", w.Lat)
				fmt.Println("Longitude", w.Lon)
				fmt.Println("Weather", w.Status)
				fmt.Println("Temperature", w.Temp)
				fmt.Println("Pressure", w.Pressure)
				fmt.Println("Humidity", w.Humidity)
				fmt.Println("Wind", w.Wind)
			default:
				fmt.Println("please specify the language (en/ja).")
				os.Exit(1)
			}
		}
	},
}

func weatherReq(country string) *WeatherJson {
	var appid string = getAppID()
	var urlweather string = fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&appid=%s", country, appid)
	req, err := http.NewRequest(http.MethodGet, urlweather, nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	weather := weatherParse(result)
	weatherjson, err := json.Marshal(weather)
	if err != nil {
		log.Fatal(err)
	}
	if globalFlags.language == "ja" {
		fmt.Println("##### JSON 取得結果 #####")
		fmt.Println(string(weatherjson))
		fmt.Printf("\n")
	} else if globalFlags.language == "en" {
		fmt.Println("##### JSON RESULT #####")
		fmt.Println(string(weatherjson))
		fmt.Printf("\n")
	}
	return weather
}

func weatherParse(result map[string]interface{}) *WeatherJson {
	coord := result["coord"].(map[string]interface{})
	main := result["main"].(map[string]interface{})
	weather := new(WeatherJson)
	weather.Lat = coord["lat"].(float64)
	weather.Lon = coord["lon"].(float64)
	weather.Status = result["weather"].([]interface{})[0].(map[string]interface{})["main"].(string)
	weather.Temp = main["temp"].(float64)
	weather.Pressure = main["pressure"].(float64)
	weather.Humidity = main["humidity"].(float64)
	weather.Wind = result["wind"].(map[string]interface{})["speed"].(float64)
	weather.Time = time.Unix(int64(result["dt"].(float64)), 0)
	weather.Country = result["sys"].(map[string]interface{})["country"].(string)
	weather.City = result["name"].(string)
	return weather
}

func getAppID() string {
	viper.SetConfigFile(configfile)
	viper.AddConfigPath("./")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return config.AppID
}

func init() {
	rootCmd.AddCommand(weatherCmd)
	weatherCmd.PersistentFlags().StringVarP(&configfile, "config", "c", "", "config file name")
}
