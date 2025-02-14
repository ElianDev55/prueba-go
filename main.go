package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
)

func main() {
    var count int
    var cities []string

    fmt.Print("Ingresa la cantidad de ciudades a buscar: ")
    fmt.Scanln(&count)

    for i := 0; i < count; i++ {
        var city string
        fmt.Print("Ingresa la ciudad: ")
        fmt.Scanln(&city)
        cities = append(cities, city)
    }

    weatherData := make(map[string]float64)
    for _, city := range cities {
        url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=bf391ec9fafcb2e9c8bc7be4a3e39ada&units=metric", city)
        data, err := GetRequest(url)
        if err != nil {
            fmt.Printf("Error en la solicitud para %s\n", city)
            continue
        }

        var result map[string]interface{}
        if err := json.Unmarshal(data, &result); err != nil {
            fmt.Printf("Error al leer el JSON para %s\n", city)
            continue
        }

        if mainData, ok := result["main"].(map[string]interface{}); ok {
            if temp, exists := mainData["temp"]; exists {
                if tempValue, ok := temp.(float64); ok {
                    weatherData[city] = tempValue
                }
            }
        }
    }

    sort.Slice(cities, func(i, j int) bool {
        return weatherData[cities[i]] > weatherData[cities[j]]
    })

    fmt.Println("Ciudades ordenadas por temperatura mayor a menor :")
    for _, city := range cities {
        if temp, ok := weatherData[city]; ok {
            fmt.Printf("%s: %.2f", city, temp)
        }
    }
}

func GetRequest(url string) ([]byte, error) {
    response, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer response.Body.Close()

    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return nil, err
    }
    return body, nil
}
