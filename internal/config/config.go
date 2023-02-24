package config

import (
	"CompanySystemsMonitoring/internal/domain/common"
	"github.com/spf13/viper"
	"strings"
	"time"
)

const (
	defaultHttpPort               = "8080"
	defaultHttpRWTimeout          = 10 * time.Second
	defaultHttpMaxHeaderMegabytes = 1
)

type (
	Config struct {
		HTTP            HTTPConfig
		FileStorageApp  FileStorageAppConfig
		ServerAPI       ServerAPIConfig
		FilesStorageAPI FilesStorageAPIConfig
	}

	HTTPConfig struct {
		Host               string        `mapstructure:"host"`
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
	}

	FileStorageAppConfig struct {
		RootPath string `mapstructure:"rootPath"`
		Alpha    string `mapstructure:"alphaFile"`
	}

	ServerAPIConfig struct {
		Address           string `mapstructure:"address"`
		Port              string `mapstructure:"port"`
		UrlMMSSystem      string `mapstructure:"urlMMSSystem"`
		UrlSupportSystem  string `mapstructure:"urlSupportSystem"`
		UrlIncidentSystem string `mapstructure:"urlIncidentSystem"`
	}

	FilesStorageAPIConfig struct {
		RootPath    string `mapstructure:"rootPath"`
		SmsFile     string `mapstructure:"smsFile"`
		VoiceFile   string `mapstructure:"voiceFile"`
		EmailFile   string `mapstructure:"emailFile"`
		BillingFile string `mapstructure:"billingFile"`
	}
)

func Init(path string) (*Config, error) {
	populateDefaults()

	if err := parseConfigFile(path); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}
	initGlobalVariables(&cfg)
	return &cfg, nil
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("filesStorageApp", &cfg.FileStorageApp); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("serverAPI", &cfg.ServerAPI); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("filesStorageAPI", &cfg.FilesStorageAPI); err != nil {
		return err
	}
	return nil
}

func parseConfigFile(filepath string) error {
	path := strings.Split(filepath, "/")
	viper.AddConfigPath(path[0])
	viper.SetConfigName(path[1])

	return viper.ReadInConfig()
}

func populateDefaults() {
	viper.SetDefault("http.port", defaultHttpPort)
	viper.SetDefault("http.max_header_megabytes", defaultHttpMaxHeaderMegabytes)
	viper.SetDefault("http.timeouts.read", defaultHttpRWTimeout)
	viper.SetDefault("http.timeouts.write", defaultHttpRWTimeout)
}

func initGlobalVariables(cfg *Config) {
	common.UrlApi = cfg.ServerAPI.Address + ":" + cfg.ServerAPI.Port
	common.UrlMMSSystem = common.UrlApi + cfg.ServerAPI.UrlMMSSystem
	common.UrlSupportSystem = common.UrlApi + cfg.ServerAPI.UrlSupportSystem
	common.UrlIncidentSystem = common.UrlApi + cfg.ServerAPI.UrlIncidentSystem
}
