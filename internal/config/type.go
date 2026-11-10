package config

import "time"

// mapstructure tags are used by viper to map config fields (you can't use yaml tags here,
// even though the config file is in yaml format)
type ConfigApp struct {
	Port       int              `mapstructure:"port"`
	Jwt        JwtConfig        `mapstructure:"jwt"`
	App        AppConfig        `mapstructure:"app"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Ai         AiConfig         `mapstructure:"ai"`
	GoogleAuth GoogleAuthConfig `mapstructure:"google-auth"`
	OpenAi     OpenAiConfig     `mapstructure:"openai"`
	AssemblyAi AssemblyAiConfig `mapstructure:"assemblyai"`
	LogLevel   string 			`mapstructure:"log_level"`
}

type AiCulturalContext struct {
	UsEnglish            string `mapstructure:"us_english"`
	UkEnglish            string `mapstructure:"uk_english"`
	InternationalEnglish string `mapstructure:"international_english"`
}

type AiConfig struct {
	BasicConversationMainObjective string            `mapstructure:"basic_conversation_main_objective"`
	BasicConversationMainRule      string            `mapstructure:"basic_conversation_main_rule"`
	RolePlayModeMainObjective      string            `mapstructure:"roleplay_mode_main_objection"`
	GenerateQuestionPrompt         string            `mapstructure:"generate_question_prompt"`
	EnglishEvaluationPrompt        string            `mapstructure:"english_evaluation_prompt"`
	ScenarioModeMainRule           string            `mapstructure:"scenario_mode_main_rule"`
	OneSentenceEnglishEvaluation   string            `mapstructure:"one_sentence_english_evaluation"`
	AiCulturalContext              AiCulturalContext `mapstructure:"cultural_context"`
}

type GoogleAuthConfig struct {
	ClientId string `mapstructure:"client_id"`
}

type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	DatabaseName string `mapstructure:"database_name"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
}

type OpenAiConfig struct {
	ApiKey      string  `mapstructure:"api_key"`
	Model       string  `mapstructure:"model"`
	Temperature float32 `mapstructure:"temperature"`
	MaxTokens   int64     `mapstructure:"max_tokens"`
}

type AssemblyAiConfig struct {
	ApiKey  string `mapstructure:"api_key"`
	BaseUrl string `mapstructure:"base_url"`
}

type JwtConfig struct {
	PrivateKeyPath  string        `mapstructure:"private_key_path"`
	PublicKeyPath   string        `mapstructure:"public_key_path"`
	Issuer          string        `mapstructure:"issuer"`
	Audience        string        `mapstructure:"audience"`
	AccessTokenTTL  time.Duration `mapstructure:"access_token_ttl"`
	RefreshTokenTTL time.Duration `mapstructure:"refresh_token_ttl"`
}

type AppConfig struct {
	PackageName string `mapstructure:"package_name"`
}
