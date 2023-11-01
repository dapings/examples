package cmd

import (
	"log"
	"strings"

	"github.com/dapings/examples/go-programing-tour-2023/tour/internal/word"
	"github.com/spf13/cobra"
)

var usage = strings.Join([]string{
	"支持的单词格式转换模式如下：",
	"1：全部转大写",
	"2：全部转小写",
	"3：下划线转大写驼峰",
	"4：下划线转小写驼峰",
	"5：驼峰转下划线",
}, "\n")

var wordCmd = &cobra.Command{
	Use:   "word",
	Short: "单词格式转换",
	Long:  usage,
	Run: func(cmd *cobra.Command, args []string) {
		var content string
		switch mode {
		case AllUpperMode:
			content = word.ToUpper(str)
		case AllLowerMode:
			content = word.ToLower(str)
		case UnderscoreToUpperCamelCaseMode:
			content = word.UnderscoreToUpperCamelCase(str)
		case UnderscoreToLowerCamelCaseMode:
			content = word.UnderscoreToLowerCamelCase(str)
		case CamelCaseToUnderscoreMode:
			content = word.CamelCaseToUnderscore(str)
		default:
			log.Fatalf("不支持该转换模式，请执行 help word 查看帮助文档")
		}
		log.Printf("output result: %s", content)
	},
}

const (
	AllUpperMode                   = iota + 1 // 全部转大写
	AllLowerMode                              // 全部转小写
	UnderscoreToUpperCamelCaseMode            // 下划线转大写驼峰
	UnderscoreToLowerCamelCaseMode            // 下划线转小写驼峰
	CamelCaseToUnderscoreMode                 // 驼峰转下划线
)

var (
	mode uint8
)

func init() {
	wordCmd.Flags().Uint8VarP(&mode, "mode", "m", 0, "请输入单词转换的模式")
	wordCmd.Flags().StringVarP(&str, "str", "s", "", "请输入单词内容")
}
