package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	colorBlue    = "\033[34m"
	colorCyan    = "\033[36m"
	colorGreen   = "\033[32m"
	colorMagenta = "\033[35m"
	colorYellow  = "\033[33m"
	colorReset   = "\033[0m"
)

type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Region  string `json:"region"`
		Country string `json:"country"`
		TzID    string `json:"tz_id"`
	} `json:"location"`
	Current struct {
		LastUpdatedEpoch int64   `json:"last_updated_epoch"`
		TempC            float64 `json:"temp_c"`
		TempF            float64 `json:"temp_f"`
		Condition        struct {
			Text string `json:"text"`
			Code int    `json:"code"`
		} `json:"condition"`
		WindMPH    float64 `json:"wind_mph"`
		WindKPH    float64 `json:"wind_kph"`
		WindDir    string  `json:"wind_dir"`
		Humidity   int     `json:"humidity"`
		FeelsLikeF float64 `json:"feelslike_f"`
		FeelsLikeC float64 `json:"feelslike_c"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			DateEpoch int64 `json:"date_epoch"`
			Day       struct {
				MaxTempC          float64 `json:"maxtemp_c"`
				MaxTempF          float64 `json:"maxtemp_f"`
				MinTempC          float64 `json:"mintemp_c"`
				MinTempF          float64 `json:"mintemp_f"`
				DailyChanceOfRain int     `json:"daily_chance_of_rain"`
				DayChanceOfSnow   int     `json:"daily_chance_of_snow"`
				Condition         struct {
					Text string `json:"text"`
					Code int    `json:"code"`
				} `json:"condition"`
			} `json:"day"`
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				TempF     float64 `json:"temp_f"`
				Condition struct {
					Text string `json:"text"`
					Code int    `json:"code"`
				} `json:"condition"`
				ChanceOfRain int64 `json:"chance_of_rain"`
				ChanceOfSnow int64 `json:"chance_of_snow"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func getWeatherEmoji(code int) string {
	emojiMap := map[int]string{
		1000: "☀️", 1003: "⛅", 1006: "☁️",
		1009: "☁️", 1030: "🌫️", 1063: "🌦️",
		1066: "🌨️", 1069: "🌨️", 1072: "🌧️",
		1087: "⛈️", 1114: "🌨️", 1135: "🌫️",
		1147: "🌫️", 1150: "🌦️", 1180: "🌦️",
		1183: "🌧️", 1192: "🌧️", 1210: "🌨️",
		1225: "❄️", 1240: "🌦️", 1273: "⛈️",
	}

	emoji, ok := emojiMap[code]
	if !ok {
		return "" // fallback if code not mapped
	}
	return emoji
}

func getHTTPRequest(city string) (Weather, error) {
	var data Weather
	city = url.QueryEscape(city)
	url := fmt.Sprintf("https://api.weatherapi.com/v1/forecast.json?key=bc3ccf0bcc444a0c8e9214338252709&q=%s&days=5&aqi=yes&alerts=no", city)
	res, err := http.Get(url)
	if err != nil {
		return Weather{}, fmt.Errorf("City not found in  or API unavailable: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return Weather{}, fmt.Errorf("City not found in  or API unavailable:(status %d)", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Weather{}, fmt.Errorf("City not found in  or API unavailable: %v", err)
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return Weather{}, fmt.Errorf("City not found in or API unavailable: %v", err)
	}

	return data, nil
}

func displayCurrent(data Weather) {
	city, region, country := data.Location.Name, data.Location.Region, data.Location.Country
	weatherCode, weatherDescription := data.Current.Condition.Code, data.Current.Condition.Text
	tempF, feelsLikeF, humidity, windSpeedMPH, windDir := data.Current.TempF, data.Current.FeelsLikeF, data.Current.Humidity, data.Current.WindMPH, data.Current.WindDir

	weatherEmoji := getWeatherEmoji(weatherCode)

	fmt.Printf("------------------------------------------\n")
	fmt.Printf("%s, %s, %s\n", city, region, country)
	fmt.Printf("------------------------------------------\n")
	fmt.Printf("%s %s\n", weatherEmoji, weatherDescription)
	fmt.Printf("Temp: %.0f\u00b0F (Feels like %.0f\u00b0F)\n", tempF, feelsLikeF)
	fmt.Printf("Humidity: %d%% | Wind: %.0f mph %s\n", humidity, windSpeedMPH, windDir)
	fmt.Printf("-----------------------------------------\n")
}

func displayDaily(data Weather) {
	city := data.Location.Name

	fmt.Printf("------------------------------------------\n")
	fmt.Printf("5 Day Forecast for %s", city)
	fmt.Printf("------------------------------------------\n")

	maxDays := 5
	totalDays := len(data.Forecast.Forecastday)

	if totalDays < maxDays {
		maxDays = totalDays
	}

	for i := 0; i < maxDays; i++ {
		day := data.Forecast.Forecastday[i]
		date := time.Unix(day.DateEpoch, 0)

		fmt.Printf(
			"%s: %s %s, High: %.0f\u00b0F, Low: %.0f\u00b0, Chance of Rain: %d%%\n",
			date.Format("Mon 01/02"),
			getWeatherEmoji(day.Day.Condition.Code),
			day.Day.Condition.Text,
			day.Day.MaxTempF,
			day.Day.MinTempF,
			day.Day.DailyChanceOfRain,
		)
	}
	fmt.Printf("------------------------------------------\n")
}

func displayHourly(data Weather) {
	city := data.Location.Name

	loc, err := time.LoadLocation(data.Location.TzID)
	if err != nil {
		loc = time.Local // fallback
	}

	fmt.Printf("------------------------------------------\n")
	fmt.Printf("Hourly Forecast for %s\n", city)
	fmt.Printf("------------------------------------------\n")

	now := time.Now().In(loc) // current time in location
	maxHours := 12
	hoursShown := 0

	for _, day := range data.Forecast.Forecastday {
		for _, h := range day.Hour {
			hourTime := time.Unix(h.TimeEpoch, 0).In(loc) // convert to local timezone

			// Skip current hours
			if !hourTime.After(now) {
				continue
			}

			fmt.Printf(
				"%s: %s %s, %.0f\u00b0F, Chance of Rain: %d%%\n",
				hourTime.Format("Mon 01/02 3 PM"),
				getWeatherEmoji(h.Condition.Code),
				h.Condition.Text,
				h.TempF,
				h.ChanceOfRain,
			)
			hoursShown++
			if hoursShown >= maxHours {
				fmt.Printf("------------------------------------------\n")
				return
			}
		}

	}
}

func askYesNo(prompt string) bool {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(prompt)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(response)
		response = strings.ToLower(response)
		switch response {
		case "y":
			return true
		case "n":
			return false
		default:
			fmt.Println("Invalid entry:")
		}
	}
}
func getUserChoice(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	choice, _ := reader.ReadString('\n')
	return strings.TrimSpace(choice)
}

func readUserInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter city name: ")
	city, err := reader.ReadString('\n')
	if err != nil {
		return city, fmt.Errorf("Enter a valid city name: %w\n", err)
	}
	city = strings.TrimSpace(city)
	return city, nil
}

func showMenu() string {
	for {
		fmt.Println("Would you like to see more details?")
		fmt.Println("1) Hourly forecast")
		fmt.Println("2) Daily forecast")
		fmt.Println("3) Enter a new city")
		fmt.Println("4) Exit")

		choice := getUserChoice("Enter choice (1-4): ")

		if choice == "1" || choice == "2" || choice == "3" || choice == "4" {
			return choice
		} else {
			fmt.Println("Invalid choice. Please enter 1-4. ")
		}
	}
}

func main() {
	for {
		city, err := readUserInput()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		data, err := getHTTPRequest(city)
		if err != nil {
			fmt.Println("Error:", err)

			// Ask user if they want to retry
			if !askYesNo("Do you want to try again? (y/n): ") {
				fmt.Println("Exiting program.")
				return
			}
			continue
		}

		displayCurrent(data)

		for {
			choice := showMenu()
			switch choice {
			case "1":
				displayHourly(data)
			case "2":
				displayDaily(data)
			case "3":
				break
			case "4":
				if askYesNo("Are you sure want exit program? (y/n): ") {
					fmt.Println("Exiting program")
					return
				}
			}
			if choice == "3" {
				break
			}
		}

	}
}
