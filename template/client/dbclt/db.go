package dbclt

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	Engine *gorm.DB
)

func InitDb(settings conf.DbSettings) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(10*time.Second))
	go func(ctx context.Context) {
		connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&charset=utf8mb4,utf8",
			settings.Username, settings.Password, settings.Hostname, settings.Dbname)
		var err1 error

		Engine, err1 = gorm.Open(mysql.Open(connStr), &gorm.Config{})
		if err1 != nil {
			panic("Database connect error," + err1.Error())
		}
		sqlDB, err := Engine.DB()
		if err != nil {
			panic("Database error")
		}
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(10000)
		sqlDB.SetConnMaxLifetime(time.Second * 3)
		cancel()
	}(ctx)

	select {
	case <-ctx.Done():
		switch ctx.Err() {
		case context.DeadlineExceeded:
			fmt.Println("context timeout exceeded")
			panic("Timeout when initialize database connection")
		case context.Canceled:
			fmt.Println("context cancelled by force. whole process is complete")
		}
	}
}
