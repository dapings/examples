package global

import (
	"os"
	"path/filepath"
	"sync"
)

var (
	RootDir string
	once    = new(sync.Once)
)

func init() {
	Init()
}

func Init() {
	once.Do(func() {
		inferRootDir()
		initConfig()
	})
}

// 推断出项目根目录。
func inferRootDir() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	var infer func(d string) string
	infer = func(d string) string {
		// 确保项目根目录下存在template目录
		if exists(filepath.Join(d, "template")) {
			return d
		}

		return infer(filepath.Dir(d))
	}

	RootDir = infer(cwd)
}

func exists(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}
