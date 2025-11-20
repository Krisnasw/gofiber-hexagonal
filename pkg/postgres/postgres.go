package postgres

import (
	"fmt"
	"log"
	"os"
	"time"

	driverPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type psql struct {
	DBHost         string
	DBPort         int
	DBUserName     string
	DBPassword     string
	DBDatabaseName string
	DBTimezone     string

	printLog     bool
	logLevel     logger.LogLevel
	logThreshold time.Duration

	maxIdleConnection             int
	maxOpenConnection             int
	connectionMaxLifetimeInSecond int
	namingStrategy                schema.Namer
}
type pgsqlOption func(*psql)

func Connect(DBHost string, DBPort int, DBUserName string, DBPassword string, DBDatabaseName string, options ...pgsqlOption) (*gorm.DB, error) {
	db := &psql{
		DBHost:         DBHost,
		DBPort:         DBPort,
		DBUserName:     DBUserName,
		DBPassword:     DBPassword,
		DBDatabaseName: DBDatabaseName,
		DBTimezone:     "UTC",

		printLog:     false,
		logLevel:     logger.Silent,
		logThreshold: 200 * time.Millisecond,

		maxIdleConnection:             5,
		maxOpenConnection:             10,
		connectionMaxLifetimeInSecond: 60,
		namingStrategy:                nil,
	}

	for _, o := range options {
		o(db)
	}

	return connect(db)
}

func connect(param *psql) (*gorm.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&TimeZone=%s",
		param.DBUserName,
		param.DBPassword,
		param.DBHost,
		param.DBPort,
		param.DBDatabaseName,
		param.DBTimezone)

	cfg := &gorm.Config{}
	if param.printLog {
		cfg.Logger = logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold: param.logThreshold,
			LogLevel:      param.logLevel,
			Colorful:      true,
		})
	}
	if param.namingStrategy != nil {
		cfg.NamingStrategy = param.namingStrategy
	}

	db, err := gorm.Open(driverPostgres.Open(dsn), cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database handle: %w", err)
	}
	sqlDB.SetMaxOpenConns(param.maxOpenConnection)
	sqlDB.SetConnMaxLifetime(time.Duration(param.connectionMaxLifetimeInSecond) * time.Second)
	sqlDB.SetMaxIdleConns(param.maxIdleConnection)

	return db, nil
}

func SetMaxIdleConns(conns int) pgsqlOption {
	return func(c *psql) {
		if conns > 0 {
			c.maxIdleConnection = conns
		}
	}
}

func SetMaxOpenConns(conns int) pgsqlOption {
	return func(c *psql) {
		if conns > 0 {
			c.maxOpenConnection = conns
		}
	}
}

func SetConnMaxLifetime(seconds int) pgsqlOption {
	return func(c *psql) {
		if seconds > 0 {
			c.connectionMaxLifetimeInSecond = seconds
		}
	}
}

func SetNamingStrategy(namingStrategy schema.Namer) pgsqlOption {
	return func(c *psql) {
		c.namingStrategy = namingStrategy
	}
}

// SetPrintLog level: 1=silent, 2=Error, 3=Warn, 4=Info. latencyThreshold: suggestion 200ms.
func SetPrintLog(isEnable bool, level logger.LogLevel, latencyThreshold time.Duration) pgsqlOption {
	return func(c *psql) {
		if latencyThreshold > 0 {
			c.printLog = isEnable
			c.logLevel = level
			c.logThreshold = latencyThreshold
		}
	}
}

func SetTimezone(timezone string) pgsqlOption {
	return func(c *psql) {
		if timezone == "" {
			c.DBTimezone = "Etc/UTC"
		} else {
			c.DBTimezone = timezone
		}
	}
}

// SetTablePrefix sets the table prefix for all tables
func SetTablePrefix(prefix string) pgsqlOption {
	return func(c *psql) {
		c.namingStrategy = schema.NamingStrategy{
			TablePrefix: prefix,
		}
	}
}
