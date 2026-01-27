package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

type CustomLogger struct {
	logger.Interface
}

func (l *CustomLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	fmt.Printf("[%s] [%s] [rows:%d] %s\n", time.Now().Format("2006-01-02 15:04:05"), formatDuration(elapsed), rows, sql)
}

func formatDuration(d time.Duration) string {
	if d >= time.Minute {
		return fmt.Sprintf("%.1fm", d.Minutes())
	}
	if d >= time.Second {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	return fmt.Sprintf("%.1fms", float64(d.Nanoseconds())/1e6)
}

func ConnectDatabase() {
	var err error

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_PORT"),
	)

	var gormLogger logger.Interface
	appEnv := os.Getenv("APP_ENV")

	if appEnv == "production" {
		gormLogger = logger.Default.LogMode(logger.Silent)
	} else {
		gormLogger = &CustomLogger{
			Interface: logger.Default.LogMode(logger.Silent),
		}
	}

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connected successfully!")
}

func GetDB() *gorm.DB {
	return DB
}
