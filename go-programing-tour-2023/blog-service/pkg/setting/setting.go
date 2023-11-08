package setting

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Setting struct {
	vp *viper.Viper
}

func NewSetting(configs ...string) (*Setting, error) {
	vp := viper.New()
	vp.SetConfigName("config")
	vp.SetConfigType("yaml")
	for _, config := range configs {
		if config != "" {
			vp.AddConfigPath(config)
		}
	}
	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}

	s := &Setting{
		vp: vp,
	}
	// 热更新：监听配置的变更
	s.WatchSettingChange()
	return s, nil
}

func (s *Setting) WatchSettingChange() {
	go func() {
		// viper 实现了对文件的监听和热更新
		s.vp.WatchConfig()
		s.vp.OnConfigChange(func(in fsnotify.Event) {
			// 热更新的文件监听事件回调
			_ = s.ReloadAllSection()
		})
	}()
}
