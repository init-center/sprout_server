package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

// Conf global variable, used to save all config info of the app
var Conf = new(AppConfig)

type AppConfig struct {
	Name            string `mapstructure:"name"`
	Mode            string `mapstructure:"mode"`
	Version         string `mapstructure:"version"`
	Port            int    `mapstructure:"port"`
	*LogConfig      `mapstructure:"log"`
	*MySQLConfig    `mapstructure:"mysql"`
	*RedisConfig    `mapstructure:"redis"`
	*SmtpConfig     `mapstructure:"smtp"`
	*SundriesConfig `mapstructure:"sundries"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DbName       string `mapstructure:"dbname"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type SmtpConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	UserName string `mapstructure:"userName"`
}

type SundriesConfig struct {
	DefaultAvatar   string `mapstructure:"default_avatar"`
	SaltPrefix      string `mapstructure:"salt_prefix"`
	JwtSecretPrefix string `mapstructure:"jwt_secret_prefix"`
}

// init settings
func Init() (err error) {
	// set config file info
	viper.SetConfigFile("./config.yaml")
	// viper.setConfigName("config")
	// viper.setConfigType(".")	// Dedicated to remote configuration center to read configuration files
	// viper.AddConfigPath(".")

	// start to read config
	err = viper.ReadInConfig()

	if err != nil {
		fmt.Printf("viper.ReadInConfig() failed, err: %v\n", err)
		return

	}

	// unmarshal the config info into the Conf variable
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
	}

	// watch config change
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("The config file has been modified")
		// when the config file changed, unmarshal the config info into the Conf variable too
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
		}
	})
	return
}
