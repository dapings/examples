package global

import (
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	SensitiveWords  []string
	MessageQueueLen = 1024
)

func initConfig() {
	viper.SetConfigName("chatroom")
	viper.AddConfigPath(filepath.Join(RootDir, "config"))

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	SensitiveWords = viper.GetStringSlice("sensitive")
	MessageQueueLen = viper.GetInt("message-queue")

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		err := viper.ReadInConfig()
		if err != nil {
			_ = viper.ReadInConfig()
		}

		SensitiveWords = viper.GetStringSlice("sensitive")
	})
}
