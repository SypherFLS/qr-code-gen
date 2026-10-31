package postgres

import (
    "fmt"
    "time"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func NewPostgres(dsn string) (*gorm.DB, error) {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("connect to postgres: %w", err)
    }

    sqlDB, err := db.DB()
    if err != nil {
        return nil, fmt.Errorf("get sql DB: %w", err)
    }

    sqlDB.SetMaxOpenConns(20)
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetConnMaxLifetime(30 * time.Minute)
    sqlDB.SetConnMaxIdleTime(5 * time.Minute)

    if err := sqlDB.Ping(); err != nil {
        sqlDB.Close()
        return nil, fmt.Errorf("ping postgres: %w", err)
    }

    return db, nil
}

func GetStats(db *gorm.DB) error{
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	stats := sqlDB.Stats()

	fmt.Println("Open:", stats.OpenConnections) // открыто соединений
	fmt.Println("InUse:", stats.InUse) // используется
	fmt.Println("Idle:", stats.Idle) // идловые значения
	fmt.Println("WaitCount:", stats.WaitCount) // сколько запросам пришлось ждать соединения
	fmt.Println("WaitDuration:", stats.WaitDuration) // время ожидания суммарно

	return nil
}