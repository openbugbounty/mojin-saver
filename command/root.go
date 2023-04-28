package command

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "mojin-saver",
		Short: "辅助工具(数据格式化)",
		Run: func(cmd *cobra.Command, args []string) {
			err := runDefaultCmd(cmd, args)
			if err != nil {
				logrus.WithError(err).Error("runDefaultCmd")
				os.Exit(1)
			}
		},
	}
)

func runDefaultCmd(_ *cobra.Command, args []string) error {
	return nil
}
