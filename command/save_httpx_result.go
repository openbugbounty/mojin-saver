package command

import (
	"fmt"
	"os"

	"github.com/hudangwei/mojin-saver/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var saveHttpxResultCmd = &cobra.Command{
	Use:   "save_httpx_result",
	Short: "保存httpx扫描结果",
	Args:  validateSaveHttpxResultArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := runSaveHttpxResultCmd(cmd, args)
		if err != nil {
			logrus.WithError(err).Error("runSaveHttpxResultCmd")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(saveHttpxResultCmd)
}

func validateSaveHttpxResultArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 && !util.IsPipedInput() {
		return fmt.Errorf("httpx结果文件为空")
	}

	return cobra.MaximumNArgs(1)(cmd, args)
}

func runSaveHttpxResultCmd(_ *cobra.Command, args []string) error {

	return nil
}
