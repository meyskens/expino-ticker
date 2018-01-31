package main

import (
	"encoding/json"
	"net/http"
	"time"

	client "github.com/influxdata/influxdb/client/v2"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Is the bitcoin up yet?")
	})
	e.GET("/diff/:setup/:metric/:interval", handleDataRequest)
	e.Logger.Fatal(e.Start(":8080"))
}

func handleDataRequest(c echo.Context) error {
	setup := c.Param("setup")
	metric := c.Param("metric")
	interval, err := time.ParseDuration(c.Param("interval"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"status": "error", "error": err.Error()})
	}

	old, err := getOldDataPoints(setup, metric, interval)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"status": "error", "error": err.Error()})
	}

	new, err := getLatestDataPoints(setup, metric)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"status": "error", "error": err.Error()})
	}

	oldAvg := getAverage(old)
	newAvg := getAverage(new)

	diff := oldAvg - newAvg
	avg := (oldAvg + newAvg) / 2
	percentage := (diff / avg) * 100

	return c.JSON(http.StatusOK, map[string]float64{"result": percentage})

}

func getAverage(data []client.Result) float64 {
	var total float64
	var count float64
	for _, item := range data {
		for _, row := range item.Series {
			for _, value := range row.Values {
				num, err := value[1].(json.Number).Float64()
				if err == nil {
					count += 1.0
					total += num
				}
			}
		}
	}

	return total / count
}