package main

import (
	"fmt"
	"os"
	"time"

	client "github.com/influxdata/influxdb/client/v2"
)

var influxDBURL = os.Getenv("INFLUXURL")

const influxDBDB = "kiosk"

func getOldDataPoints(setup, metric string, interval time.Duration) ([]client.Result, error) {
	t := time.Now().Add(-1 * interval)
	end := t.Add(-1 * time.Minute)

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: influxDBURL,
	})
	if err != nil {
		return nil, err
	}
	return queryDB(c, fmt.Sprintf("SELECT %s FROM %s WHERE time <= %dms AND time >= %dms", metric, setup, t.Unix()*1000, end.Unix()*1000))
}

func getLatestDataPoints(setup, metric string) ([]client.Result, error) {
	t := time.Now()
	end := t.Add(-1 * time.Minute)

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: influxDBURL,
	})
	if err != nil {
		return nil, err
	}
	return queryDB(c, fmt.Sprintf("SELECT %s FROM %s WHERE time <= %dms AND time >= %dms", metric, setup, t.Unix()*1000, end.Unix()*1000))
}

// queryDB convenience function to query the database
func queryDB(clnt client.Client, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: "kiosk",
	}
	if response, err := clnt.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}
