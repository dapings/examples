package cmd

import (
	"log"

	"github.com/dapings/examples/go-programing-tour-2023/tour/internal/json2struct"
	"github.com/spf13/cobra"
)

var jsonCmd = &cobra.Command{
	Use:   "json",
	Short: "JSON转换和处理",
	Long:  "JSON转换和处理",
	Run:   func(cmd *cobra.Command, args []string) {},
}

var json2structCmd = &cobra.Command{
	Use:   "struct",
	Short: "json转换",
	Long:  "json转换",
	Run: func(cmd *cobra.Command, args []string) {
		parser, err := json2struct.NewParser(str)
		if err != nil {
			log.Fatalf("json2struct.NewParser err: %v", err)
		}
		content := parser.JSON2Struct()
		log.Printf("output result: %s", content)
	},
}

func init() {
	jsonCmd.AddCommand(json2structCmd)
	json2structCmd.Flags().StringVarP(&str, "str", "s", "", "")
}
