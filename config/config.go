package config

import "sync"

type ConfigModel struct {
	Debug      bool   `yaml:"Debug"`
	Admin_id   int    `yaml:"admin_id"`
	Test_group int    `yaml:"test_group"`
	ZanaoToken string `yaml:"zanao_token"`
	Group      struct {
		Include []string `yaml:"include"`
		Exclude []string `yaml:"exclude"`
	} `yaml:"group"`
	Friend struct {
		Include []string `yaml:"include"`
		Exclude []string `yaml:"exclude"`
	} `yaml:"friend"`
	AI_Prompt string `yaml:"ai_prompt"`
}

// 配置
var Config_mu sync.RWMutex
var Config ConfigModel
