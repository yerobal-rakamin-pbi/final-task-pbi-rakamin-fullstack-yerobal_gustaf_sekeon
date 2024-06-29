package configreader

import (
	"github.com/spf13/viper"
)

type Interface interface {
	ReadConfig(configVar interface{})
}

type configreader struct {
	configFile string
}

func Init(configFile string) Interface {
	return &configreader{
		configFile: configFile,
	}
}

func (c *configreader) ReadConfig(configVar interface{}) {
	viper.SetConfigFile(c.configFile)

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&configVar); err != nil {
		panic(err)
	}
}
