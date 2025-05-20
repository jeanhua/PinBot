package model

type Config struct {
	Debug      bool `yaml:"Debug"`
	Admin_id   int  `yaml:"admin_id"`
	Test_group int  `yaml:"test_group"`
	Group      struct {
		Include []string `yaml:"include"`
		Exclude []string `yaml:"exclude"`
	} `yaml:"group"`
	Friend struct {
		Include []string `yaml:"include"`
		Exclude []string `yaml:"exclude"`
	} `yaml:"friend"`
}
