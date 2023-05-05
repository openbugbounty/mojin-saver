package command

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/hudangwei/couchdb"
	"github.com/hudangwei/mojin-saver/db"
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
		} else {
			fmt.Println(err)
		}
	}
	if len(results) == 0 {
		return errors.New("httpx result empty")
	}

	var docs []couchdb.CouchDoc
	for _, v := range results {
		info := &WebsiteInfo{
			Type:        "website",
			Program:     program,
			HttpxResult: v,
		}
		info.SetID(v.Url)
		docs = append(docs, info)
	}
	return SaveWebsiteInfo(docs)
}

/*
bbrf url add url statuscode contentlength -t title:x -t webserver:x -t contenttype:x -t contentlength:x -t statuscode:x -t serverresponse:x -p @INFER
*/
type HttpxResult struct {
	Timestamp string `json:"timestamp,omitempty"`
	Hash      *Hash  `json:"hash,omitempty"`
	Input     string `json:"input,omitempty"`
	Url       string `json:"url,omitempty"`
	Scheme    string `json:"scheme,omitempty"`
	Method    string `json:"method,omitempty"`
	Host      string `json:"host,omitempty"`
	Port      string `json:"port,omitempty"`
	Path      string `json:"path,omitempty"`

	Title         string   `json:"title,omitempty"`
	WebServer     string   `json:"webserver,omitempty"`
	ContentType   string   `json:"content_type,omitempty"`
	ContentLength int      `json:"content_length,omitempty"`
	StatusCode    int      `json:"status_code,omitempty"`
	Technologies  []string `json:"tech,omitempty"`

	Words  int  `json:"words,omitempty"`
	Lines  int  `json:"lines,omitempty"`
	Failed bool `json:"failed,omitempty"`

	A        []string `json:"a,omitempty"`
	CNames   []string `json:"cnames,omitempty"`
	FinalUrl string   `json:"final-url,omitempty"`
	TLSData  *TLS     `json:"tls-grab,omitempty"`
	Favicon  string   `json:"favicon-mmh3,omitempty"`
}

type Hash struct {
	BodyMd5       string `json:"body_md5"`
	BodyMmh3      string `json:"body_mmh3"`
	BodySha256    string `json:"body_sha256"`
	BodySimhash   string `json:"body_simhash"`
	HeaderMd5     string `json:"header_md5"`
	HeaderMmh3    string `json:"header_mmh3"`
	HeaderSha256  string `json:"header_sha256"`
	HeaderSimhash string `json:"header_simhash"`
}
type TLS struct {
	DNSName            []string `json:"dns_names,omitempty"`
	CommonName         []string `json:"common_name,omitempty"`
	Organization       []string `json:"organization,omitempty"`
	IssuerCommonName   []string `json:"issuer_common_name,omitempty"`
	IssuerOrganization []string `json:"issuer_organization,omitempty"`
}

type WebsiteInfo struct {
	couchdb.Document
	Type    string `json:"type"`
	Program string `json:"program"`
	HttpxResult
}

func SaveWebsiteInfo(docs []couchdb.CouchDoc) error {
	if db.GlobalCouchDB == nil {
		return errors.New("couch db instance is nil")
	}

	s := db.GlobalCouchDB.Use(GlobalConfig.CouchDB.DBName)
	return s.MultiStore(docs)
}
