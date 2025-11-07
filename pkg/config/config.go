package config

import (
	"os"
	"strconv"
)

type Config struct {
	accessTokenSecret  []byte
	accessTokenExpiry  int
	refreshTokenSecret []byte
	refreshTokenExpiry int
}

func (c *Config) GetAccessTokenSecret() []byte {
	return c.accessTokenSecret
}
func (c *Config) GetRefreshTokenSecret() []byte {
	return c.refreshTokenSecret
}
func (c *Config) GetAccessTokenExpiry() int {
	return c.accessTokenExpiry
}
func (c *Config) GetRefreshTokenExpiry() int {
	return c.refreshTokenExpiry
}

var config *Config

func GetConfig() *Config {
	if config != nil {
		return config
	}
	accessTknExp, err := strconv.Atoi(os.Getenv("ACCESS_TKN_EXP"))
	if err != nil {
		accessTknExp = 5
	}
	refreshTknExp, err := strconv.Atoi(os.Getenv("ACCESS_TKN_EXP"))
	if err != nil {
		refreshTknExp = 72
	}

	config = &Config{
		accessTokenSecret:  []byte(os.Getenv("ACCESS_TKN_SECRET")),
		refreshTokenSecret: []byte(os.Getenv("REFRESH_TKN_SECRET")),
		accessTokenExpiry:  accessTknExp,
		refreshTokenExpiry: refreshTknExp,
	}
	return config
}
