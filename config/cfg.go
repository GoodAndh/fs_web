package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

var Env Config

type Config struct {
	TestOrMain   string
	PublicHost   string
	Port         string
	DBUser       string
	DBPassword   string
	DBAddress    string
	DBName       string
	DBTest       string
	JWTSecretKey string
}

func (c *Config) Main() Config {
	c.TestOrMain = "main"

	return *c

}

func (c *Config) Test() Config {
	c.TestOrMain = "test"
	return *c

}

func (c Config) InitConfig() Config {

	if c.TestOrMain == "" {
		log.Fatal("db unknown for product or testing")
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	envPath := ""

	if c.TestOrMain == "main" {
		envPath = filepath.Join(dir, ".", ".env")
		// envPath = filepath.Join(dir, "..", "..", "..", ".env")

	}

	if c.TestOrMain == "test" {
		// envPath = filepath.Join(dir, ".", ".env")
		envPath = filepath.Join(dir, "..", "..", "..", ".env")
	}

	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("env path for db %s error,message:%v", c.TestOrMain, err)
	}

	return Config{
		TestOrMain:   Env.TestOrMain,
		PublicHost:   getEnv("PUBLIC_HOST", "http://localhost:"),
		Port:         getEnv("PORT", "3000"),
		DBUser:       getEnv("DB_USER", "root"),
		DBPassword:   getEnv("DB_PASSWORD", "r23password"),
		DBAddress:    fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		DBName:       getEnv("DB_NAME", "YOUR_NAME"),
		DBTest:       getEnv("DB_NAME_TEST", "YOUR_NAME_TEST"),
		JWTSecretKey: getEnv("SECRET_KEY", "12IU3YJSGDKAH91Y28	HJASHhhjashdkj"),
	}

}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
