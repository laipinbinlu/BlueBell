package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

var Conf = new(Config)

// 将配置文件读到结构体中   -- >特别注意tag的使用，防止匹配不上
type Config struct {
	Name      string `mapstructure:"name"`
	Mode      string `mapstructure:"mode"`
	Port      int    `mapstructure:"port"`
	Version   string `mapstructure:"version"`
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_id"`

	*LogConfig   `mapstructure:"log"`
	*MysqlConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}
type LogConfig struct {
	Level      string `mapstructure:"level"`
	FileName   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}
type MysqlConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DBName    string `mapstructure:"dbname"`
	MaxIdle   int    `mapstructure:"max_idle_conns"`
	MaxActive int    `mapstructure:"max_active_conns"`
}
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// 加载配置文件
func Init() (err error) {
	// 1. 直接指定配置文件路径  相对路径和绝对路径都可以
	viper.SetConfigFile("./conf/config.yaml")
	// 方法二：viper查找文件
	//viper.SetConfigName("config") // 查找文件的名称，不需要指示后缀   注意同名的情况
	//viper.AddConfigPath(".")      // 查找文件的路径,（相对路径）

	//viper.SetConfigFile(filepath)

	err = viper.ReadInConfig() // 查找并读取配置文件
	if err != nil {            // 处理读取配置文件的错误
		fmt.Printf("viper.ReadInConfig(),err: %v\n ", err)
		return err
	}

	if err = viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal(),err: %v\n ", err)
		return
	}

	// 监听配置文件
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件发生了变化")
		if err = viper.Unmarshal(Conf); err != nil { // 文件改变了所以任然需要反序化到结构体
			fmt.Printf("viper.Unmarshal(),err: %v\n ", err)
			return
		}
	})
	return
}
