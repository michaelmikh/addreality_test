package pgsql

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"github.com/jmoiron/sqlx"
	// pq lib to initialize postgresql driver.
	_ "github.com/lib/pq"
)

// config contains all data related to establishing PostgreSQL connection.
type config struct {
	DBUsername string `json:"db_username"`
	DBPassword string `json:"db_password"`
	Database   string `json:"db_name"`
}

// PostgreSQL configuration and connection.
var (
	conf config
	sql  *sqlx.DB
)

func init() {
	getConnection()
}

// lastID stores last fetched device_metrics ID,
// so fetching will go on as table is modified.
var lastID = 0

// DeviceMetricsRow represents a single row from device_metrics table.
type DeviceMetricsRow struct {
	ID         uint      `db:"id"`
	DeviceID   uint      `db:"device_id"`
	Metric1    int       `db:"metric_1"`
	Metric2    int       `db:"metric_2"`
	Metric3    int       `db:"metric_3"`
	Metric4    int       `db:"metric_4"`
	Metric5    int       `db:"metric_5"`
	ServerTime time.Time `db:"server_time"`
}

// PollDB provides polling mechanism and launches getAllMetrics every 5 seconds.
func PollDB(rowsCh chan<- []DeviceMetricsRow, errCh chan<- error) {
	for {
		<-time.After(5 * time.Second)
		go getAllMetrics(rowsCh, errCh)
	}
}

// getAllMetrics reads device_metrics table and sends retrieved data over the channel,
// the first time function fetches all the rows of device_metrics,
// then it retrieves only the rows that was added later.
func getAllMetrics(rowsCh chan<- []DeviceMetricsRow, errCh chan<- error) {
	var err error
	rows := []DeviceMetricsRow{}

	if lastID == 0 {
		err = sql.Select(rows,
			"SELECT id, device_id, metric_1, metric_2, metric_3, metric_4, metric_5, server_time FROM device_metrics ORDER BY id ASC")

		if err != nil {
			errCh <- err
			return
		}

		lastID += len(rows)

		rowsCh <- rows
	} else {
		err = sql.Select(rows,
			"SELECT id, device_id, metric_1, metric_2, metric_3, metric_4, metric_5, server_time FROM device_metrics WHERE id>? ORDER BY id ASC",
			lastID)

		if err != nil {
			errCh <- err
			return
		}

		lastID += len(rows)

		rowsCh <- rows
	}
}

// CreateAlert inserts alert for given device in device_alerts.
func (r DeviceMetricsRow) CreateAlert(message string) error {
	var err error

	_, err = sql.Exec("INSERT INTO device_alerts (device_id, message) VALUES (?, ?)", r.DeviceID, message)

	return err
}

// getConnection establishes PostgreSQL connection.
func getConnection() {
	var err error

	file, err := filepath.Abs("pgsql/config.json")
	if err != nil {
		log.Fatal(err)
	}

	fileContent, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(fileContent, &conf)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to PostgreSQL.
	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		conf.DBUsername, conf.DBPassword, conf.Database)

	if sql, err = sqlx.Connect("postgres", dbInfo); err != nil {
		log.Fatal(err)
	}
}
