package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigModel struct {
	Debug            bool   `yaml:"debug"`
	NapCatServerUrl  string `yaml:"napcatServerUrl"`
	LocalListenPort  uint   `yaml:"localListenPort"`
	Admin_id         uint   `yaml:"admin_id"`
	Test_group       uint   `yaml:"test_group"`
	ZanaoToken       string `yaml:"zanao_token"`
	SiliconflowToken string `yaml:"siliconflow_token"`
	TavilyToken      string `yaml:"tavilyToken"`
	Group            struct {
		Include []string `yaml:"include"`
		Exclude []string `yaml:"exclude"`
	} `yaml:"group"`
	Friend struct {
		Include []string `yaml:"include"`
		Exclude []string `yaml:"exclude"`
	} `yaml:"friend"`
	AI_Prompt string `yaml:"ai_prompt"`
	HelpWords struct {
		Group  string `yaml:"group"`
		Friend string `yaml:"friend"`
	} `yaml:"help_words"`
}

// 配置
var config ConfigModel

func GetConfig() *ConfigModel {
	return &config
}

func LoadConfig() {
	file, err := os.Open("./config.yaml")
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Fatalf("Failed to decode config: %v", err)
	}
}
