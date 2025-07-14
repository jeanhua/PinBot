package config

type ConfigModel struct {
	Debug            bool   `yaml:"debug"`
	Admin_id         int    `yaml:"admin_id"`
	Test_group       int    `yaml:"test_group"`
	ZanaoToken       string `yaml:"zanao_token"`
	SiliconflowToken string `yaml:"siliconflow_token"`
	BochaToken       string `yaml:"bochaSearchToken"`
	Group            struct {
		Include []string `yaml:"include"`
		Exclude []string `yaml:"exclude"`
	} `yaml:"group"`
	Friend struct {
		Include []string `yaml:"include"`
		Exclude []string `yaml:"exclude"`
	} `yaml:"friend"`
	AI_Prompt        string `yaml:"ai_prompt"`
	FunctionCallMaxC int    `yaml:"functionCall_maxC"`

	HelpWords struct {
		Group  string `yaml:"group"`
		Friend string `yaml:"friend"`
	} `yaml:"help_words"`
}

// 配置
var ConfigInstance ConfigModel
