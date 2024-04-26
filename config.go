package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	BotUrl      string `json:"bot_url"`
	Key         string `json:"api_key"`
	PictureDir  string `json:"picture_dir"`
	FailbackDir string `json:"failback_picture_dir"`
	LogFile     string `json:"log_file"`
	PollTimeout int    `json:"poll_timeout"`
	BotPwd      string `json:"bot_pwd"`
	Debug       bool   `json:"debug"`
}

func (c *Config) readConfigFile(s string) error {

	data, err := ioutil.ReadFile(s)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, c); err != nil {
		return err
	}
	return nil
}
