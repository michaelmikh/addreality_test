package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// DatabaseConfig contains all data related to establishing PostgreSQL connection.
type DatabaseConfig struct {
	DBUsername string `json:"db_username"`
	DBPassword string `json:"db_password"`
	Database   string `json:"db_name"`
}

// RedisConfig contains all data related to establishing Redis connection.
type RedisConfig struct {
	Address  string `json:"redis_address"`
	Password string `json:"redis_password"`
	Database int    `json:"redis_db"`
}

// PostgreSQL and Redis connections.
var (
	sql         *sqlx.DB
	redisClient *redis.Client
)

// deviceMetricsRow represents a single row from device_metrics table.
type deviceMetricsRow struct {
	ID         uint      `db:"id"`
	DeviceID   uint      `db:"device_id"`
	Metric1    int       `db:"metric_1"`
	Metric2    int       `db:"metric_2"`
	Metric3    int       `db:"metric_3"`
	Metric4    int       `db:"metric_4"`
	Metric5    int       `db:"metric_5"`
	ServerTime time.Time `db:"server_time"`
}

// getAllMetrics reads device_metrics table and stores all data in slice of deviceMetricsRow.
func getAllMetrics() ([]deviceMetricsRow, error) {
	var err error
	rows := []deviceMetricsRow{}

	err = sql.Select(rows,
		"SELECT id, device_id, metric_1, metric_2, metric_3, metric_4, metric_5, server_time FROM device_metrics ORDER BY id ASC")

	return rows, err
}

// createAlert inserts alert for given device in device_alerts.
func (r deviceMetricsRow) createAlert(message string) error {
	var err error

	_, err = sql.Exec("INSERT INTO device_alerts (device_id, message) VALUES (?, ?)", r.DeviceID, message)

	return err
}

// getDBConn establishes PostgreSQL connection.
func getDBConn() {
	var err error

	// Connect to PostgreSQL.
	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		config.DatabaseConfig.DBUsername, config.DatabaseConfig.DBPassword, config.DatabaseConfig.Database)

	if sql, err = sqlx.Connect("postgres", dbInfo); err != nil {
		log.Fatal(err)
	}
}

// getRedisConn establishes Redis connection.
func getRedisConn() {
	var err error

	redisClient = redis.NewClient(&redis.Options{
		Addr:     config.RedisConfig.Address,
		Password: config.RedisConfig.Password,
		DB:       config.RedisConfig.Database,
	})

	_, err = redisClient.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}
}
