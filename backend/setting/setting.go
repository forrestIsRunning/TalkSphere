package setting

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf = new(Config)

type Config struct {
	*AppConfig       `mapstructure:"app"`
	*LogConfig       `mapstructure:"log"`
	*MysqlConfig     `mapstructure:"mysql"`
	*RedisConfig     `mapstructure:"redis"`
	*GinConfig       `mapstructure:"gin"`
	*SnowFlakeConfig `mapstructure:"snowflake"`
	*EncryptConfig   `mapstructure:"encrypt"`
	*AuthConfig      `mapstructure:"auth"`
	*OSSConfig       `mapstructure:"oss"`
}

type AppConfig struct {
	Name    string `mapstructure:"name"`
	Mode    string `mapstructure:"mode"`
	Version string `mapstructure:"version"`
	Port    int    `mapstructure:"port"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	FileName   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type MysqlConfig struct {
	Host              string `mapstructure:"host"`
	Port              int    `mapstructure:"port"`
	User              string `mapstructure:"user"`
	PassWord          string `mapstructure:"password"`
	DB                string `mapstructure:"db"`
	MaxOpenConnection int    `mapstructure:"max_open_connection"`
	MaxIdleConnection int    `mapstructure:"max_idle_connection"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	PassWord string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type EncryptConfig struct {
	SecretKey string `mapstructure:"secret_key"`
}

type GinConfig struct {
	Mode string `mapstructure:"mode"`
}

type SnowFlakeConfig struct {
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_id"`
}

type AuthConfig struct {
	JwtExpire int64 `mapstructure:"jwt_expire"`
}

type OSSConfig struct {
	BucketName string `mapstructure:"bucket_name"`
	Region     string `mapstructure:"region"`
	SecretID   string `mapstructure:"secret_id"`
	SecretKey  string `mapstructure:"secret_key"`
}

func Init() (err error) {
	viper.SetConfigName("config") // 指定配置文件名称（不需要带后缀）
	viper.SetConfigType("yaml")   // 指定配置文件类型
	//viper.AddConfigPath(".")      // 指定查找配置文件的路径（这里使用相对路径）
	viper.AddConfigPath("./conf/")
	err = viper.ReadInConfig() // 读取配置信息
	if err != nil {
		// 读取配置信息失败
		fmt.Printf("viper.ReadInConfig() failed, err: %v\n", err)
		return
	}
	// 把读取到的信息反序列化到 Conf 变量中
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal() failed, err: %v\n", err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("Configure file changed ...")
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal() failed, err: %v\n", err)
		}
	})
	return
}
