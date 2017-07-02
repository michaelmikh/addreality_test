package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/michaelmikh/addreality_test/email"
)

// AppConfig contains all data from config.json in current directory.
type AppConfig struct {
	Metric1        int `json:"metric1"`
	Metric2        int `json:"metric2"`
	Metric3        int `json:"metric3"`
	Metric4        int `json:"metric4"`
	Metric5        int `json:"metric5"`
	DatabaseConfig `json:"PostgreSQL"`
	RedisConfig    `json:"Redis"`
}

var config AppConfig

const (
	// emailRecipient is an e-mail address of alert recipient.
	emailRecipient = "admin@addreality.com"
	// emailRetryTimes is how many times system will try to send alert e-mail
	emailRetryTimes = 3
)

// init parses config.json into AppConfig struct
// and sets up PostgreSQL and Redis connections.
func init() {
	parseConfig()
	getDBConn()
	getRedisConn()
}

func main() {
	rows, err := getAllMetrics()
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range rows {
		processRow(row)
	}
}

// parseConfig reads config.json file and populates AppConfig data structure.
func parseConfig() {
	var err error

	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatal(err)
	}
}

// processRow processes single device_alert row,
// sends e-mail if metrics are out of bounds,
// sets alertText to Redis (key is DeviceID).
func processRow(row deviceMetricsRow) {
	if (row.Metric1 != config.Metric1) || (row.Metric2 != config.Metric2) ||
		(row.Metric3 != config.Metric3) || (row.Metric4 != config.Metric4) ||
		(row.Metric5 != config.Metric5) {
		var err error

		deviceID := fmt.Sprintf("%d", row.DeviceID)
		alertText := fmt.Sprintf("Device %d metric is out of bounds! Server time: %v", row.DeviceID, row.ServerTime)

		redisClient.Set(deviceID, alertText, 0)

		err = row.createAlert(alertText)
		if err != nil {
			log.Println(err)
		}

		for i := emailRetryTimes; i > 0; i-- {
			errMail := email.SendAlert(emailRecipient, alertText)
			if errMail != nil {
				log.Println(errMail)
			} else {
				break
			}
		}
	}
}
