package command

import (
	"encoding/json"
	"errors"
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
	var resultFile string
	if len(args) > 0 {
		resultFile = args[0]
	}
	lines, err := util.ReadLine(resultFile)
	if err != nil {
		return err
	}
	if len(lines) == 0 {
		return errors.New("httpx result empty")
	}
	var results []HttpxResult
	for _, v := range lines {
		var result HttpxResult
		if err := json.Unmarshal([]byte(v), &result); err == nil {
			results = append(results, result)
		}
	}
	if len(results) == 0 {
		return errors.New("httpx result empty")
	}

	return nil
}

/*
bbrf url add url statuscode contentlength -t title:x -t webserver:x -t contenttype:x -t contentlength:x -t statuscode:x -t serverresponse:x -p @INFER
*/
type HttpxResult struct {
	A             []string `json:"a,omitempty"`
	CNames        []string `json:"cnames,omitempty"`
	Url           string   `json:"url,omitempty"`
	Host          string   `json:"host,omitempty"`
	Title         string   `json:"title,omitempty"`
	WebServer     string   `json:"webserver,omitempty"`
	ContentType   string   `json:"content_type,omitempty"`
	ContentLength string   `json:"content_length,omitempty"`
	StatusCode    int      `json:"status_code,omitempty"`
	FinalUrl      string   `json:"final-url,omitempty"`
	TLSData       *TLS     `json:"tls-grab,omitempty"`
	Technologies  []string `json:"tech,omitempty"`
	Favicon       string   `json:"favicon-mmh3,omitempty"`
}

type TLS struct {
	DNSName            []string `json:"dns_names,omitempty"`
	CommonName         []string `json:"common_name,omitempty"`
	Organization       []string `json:"organization,omitempty"`
	IssuerCommonName   []string `json:"issuer_common_name,omitempty"`
	IssuerOrganization []string `json:"issuer_organization,omitempty"`
}
