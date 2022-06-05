package config

type ServiceConfig struct {
	Name  string
	Types []TypeConfig `yaml:"types"`
}

type TypeConfig struct {
	Name        string          `yaml:"name"`
	ListAPI     ListAPIConfig   `yaml:"listApi"`
	GetTagsAPI  GetTagAPIConfig `yaml:"getTagsApi"`
	Description string          `yaml:"description"`
	Global      bool            `yaml:"global"`
	Tags        TagField        `yaml:"tags"`
}

type TagField struct {
	Field string `yaml:"field"`
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

func (t TagField) IsZero() bool {
	return t.Field == ""
}

type ListAPIConfig struct {
	Call       string   `yaml:"call"`
	Pagination bool     `yaml:"pagination"`
	OutputKey  []string `yaml:"outputKey"`
	IDField    string   `yaml:"id"`
}

type GetTagAPIConfig struct {
	Call         string `yaml:"call"`
	InputIDField string `yaml:"inputIDField"`
	OutputKey    string `yaml:"outputKey"`
}

type RootConfig struct {
	Services []string `yaml:"services"`
}

type Config struct {
	Services []ServiceConfig
}

func newConfig() *Config {
	return &Config{
		// Services: map[string]ServiceConfig{},
	}
}