package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App struct {
	Secret             string
	RefreshTokenSecret string
}

type DB struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
	Schema   string
}

type Stripe struct {
	Secret        string
	WebhookSecret string
}

type Config struct {
	App    App
	DB     DB
	Stripe Stripe
}

func InitConfig() *Config {
	c := Config{}

	c.App.Secret = os.Getenv("APP_SECRET")
	c.App.RefreshTokenSecret = os.Getenv("APP_REFRESH_TOKEN_SECRET")
	c.DB.Host = os.Getenv("DB_HOST")
	c.DB.Port, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	c.DB.Database = os.Getenv("DB_NAME")
	c.DB.User = os.Getenv("DB_USERNAME")
	c.DB.Password = os.Getenv("DB_PASSWORD")
	c.DB.Schema = os.Getenv("DB_SCHEMA")

	c.Stripe.Secret = os.Getenv("STRIPE_SECRET_KEY")
	c.Stripe.WebhookSecret = os.Getenv("STRIPE_WEBHOOK_SECRET")

	return &c
}

func PostgresDSN() string {
	c := Config{}

	c.DB.Host = os.Getenv("DB_HOST")
	c.DB.Port, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	c.DB.Database = os.Getenv("DB_NAME")
	c.DB.User = os.Getenv("DB_USERNAME")
	c.DB.Password = os.Getenv("DB_PASSWORD")
	c.DB.Schema = os.Getenv("DB_SCHEMA")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", c.DB.User, c.DB.Password, c.DB.Host, c.DB.Port, c.DB.Database)

	return dsn
}

func (c *Config) NewDatabase() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.DB.Host, c.DB.Port, c.DB.User, c.DB.Password, c.DB.Database)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatalln("failed to connect database", err)

		return nil, err
	}

	return db, err
}

func (c *Config) CloseDatabase(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		logrus.Fatalln("failed to get connection from database", err)

		return
	}

	if err = dbSQL.Close(); err != nil {
		logrus.Fatalln("failed to close connection from database", err)
	}
}
