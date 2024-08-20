package main

import (
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"
)

type StoreType int

const (
	StoreTypeFile StoreType = iota
	StoreTypeDynamodb
)

type Config struct {
	SystemCode  string
	Host        string
	Port        string
	BackendAddr string
	StoreType   StoreType
	NoCacheList []*regexp.Regexp
}

func (c Config) GetAddress() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func NewConfigFromEnv() Config {
	c := Config{
		SystemCode: "blitz-proxy",
		Host:       "127.0.0.1",
		Port:       "3000",
		StoreType:  StoreTypeFile,
	}
	if v := os.Getenv("SYSTEM_CODE"); v != "" {
		c.SystemCode = v
	}
	if v := os.Getenv("HOST"); v != "" {
		c.Host = v
	}
	if v := os.Getenv("PORT"); v != "" {
		c.Port = v
	}
	if v := os.Getenv("BACKEND"); v != "" {
		c.BackendAddr = v
	}
	if v := os.Getenv("STORE_TYPE"); v != "" {
		if v == "DYNAMODB" {
			c.StoreType = StoreTypeDynamodb
		}
	}
	if v := os.Getenv("NO_CACHE_LIST"); v != "" {
		list := strings.Split(v, ",")
		for _, val := range list {
			c.NoCacheList = append(c.NoCacheList, regexp.MustCompile(strings.TrimSpace(val)))
		}
	}
	return c
}

func GetLogLevel() slog.Level {
	ll := os.Getenv("LOG_LEVEL")
	switch ll {
	case "DEBUG":
		return slog.LevelDebug
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
