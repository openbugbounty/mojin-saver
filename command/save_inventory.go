package command

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/hudangwei/couchdb"
	"github.com/hudangwei/mojin-saver/db"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var saveInventoryCmd = &cobra.Command{
	Use:   "save_inventory",
	Short: "保存Inventory资产信息",
	Run: func(cmd *cobra.Command, args []string) {
		err := runSaveInventoryCmd(cmd, args)
		if err != nil {
			logrus.WithError(err).Error("runSaveInventoryCmd")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(saveInventoryCmd)
}

func runSaveInventoryCmd(_ *cobra.Command, args []string) error {
	return downloadInventory()
}

func downloadInventory() error {
	tempdir, err := os.MkdirTemp("", "bbp-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempdir)

	_, err = git.PlainClone(tempdir, false, &git.CloneOptions{
		URL:           "https://github.com/trickest/inventory",
		Progress:      os.Stdout,
		Depth:         1,
		SingleBranch:  true,
		ReferenceName: plumbing.Main,
	})
	if err != nil {
		return err
	}
	p := filepath.Join(tempdir, "targets.json")
	f, err := os.Open(p)
	if err != nil {
		return err
	}
	defer f.Close()
	var data BugBountyData
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return err
	}

	var hasServerReportCnt int
	var noServerReportCnt int
	for _, item := range data.Targets {
		if err := SaveProgram(item.Name); err != nil {
			fmt.Println(err)
		}
		if err := saveServerReportCSV(item.Name, filepath.Join(tempdir, item.Name)); err != nil {
			fmt.Println(err)
			noServerReportCnt += 1
		} else {
			hasServerReportCnt += 1
		}
	}

	fmt.Println("hasServerReportCnt:", hasServerReportCnt)
	fmt.Println("noServerReportCnt:", noServerReportCnt)
	return nil
}

func getCSVFieldValue(program string, record []string, fieldMap map[string]int, field string) string {
	if i, ok := fieldMap[field]; ok {
		return record[i]
	}
	return ""
}

func saveServerReportCSV(program, dir string) error {
	file, err := os.Open(filepath.Join(dir, "server-report.csv"))
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}
	fieldMap := make(map[string]int)
	var docs []couchdb.CouchDoc
	for k, record := range records {
		if k == 0 {
			for i, field := range record {
				fieldMap[field] = i
			}
			continue
		} else {
			input := getCSVFieldValue(program, record, fieldMap, "input")
			url := getCSVFieldValue(program, record, fieldMap, "url")
			scheme := getCSVFieldValue(program, record, fieldMap, "scheme")
			method := getCSVFieldValue(program, record, fieldMap, "method")
			host := getCSVFieldValue(program, record, fieldMap, "host")
			port := getCSVFieldValue(program, record, fieldMap, "port")
			path := getCSVFieldValue(program, record, fieldMap, "path")
			title := getCSVFieldValue(program, record, fieldMap, "title")
			webserver := getCSVFieldValue(program, record, fieldMap, "webserver")
			contenttype := getCSVFieldValue(program, record, fieldMap, "content_type")
			contentlength := getCSVFieldValue(program, record, fieldMap, "content_length")
			statuscode := getCSVFieldValue(program, record, fieldMap, "status_code")
			tech := getCSVFieldValue(program, record, fieldMap, "tech_")
			finalurl := getCSVFieldValue(program, record, fieldMap, "final_url")
			contentLen, _ := strconv.Atoi(contentlength)
			statusCode, _ := strconv.Atoi(statuscode)

			httpxResult := HttpxResult{
				Input:         input,
				Url:           url,
				Scheme:        scheme,
				Method:        method,
				Host:          host,
				Port:          port,
				Path:          path,
				Title:         title,
				WebServer:     webserver,
				ContentType:   contenttype,
				ContentLength: contentLen,
				StatusCode:    statusCode,
				FinalUrl:      finalurl,
			}
			if len(tech) > 0 {
				httpxResult.Technologies = []string{tech}
			}
			info := &WebsiteInfo{
				Type:        "website",
				Program:     program,
				Tags:        make(Tags),
				HttpxResult: httpxResult,
			}
			info.Tags.AddTag("title", title)
			info.Tags.AddTag("webserver", webserver)
			info.Tags.AddTag("contenttype", contenttype)
			info.Tags.AddTag("contentlength", contentlength)
			info.Tags.AddTag("statuscode", statuscode)
			if len(tech) > 0 {
				info.Tags.AddTag("tech", []string{tech})
			}
			info.SetID(url)
			docs = append(docs, info)
		}
	}
	if len(docs) > 0 {
		return SaveWebsiteInfo(docs)
	}
	return nil
}

type BugBountyData struct {
	Targets []Target `json:"targets"`
}

type Target struct {
	Name    string   `json:"name"`
	Domains []string `json:"domains"`
	URL     string   `json:"url"`
}

type Program struct {
	couchdb.Document
	Type string `json:"type"`
	Tags Tags   `json:"tags"`
}

func SaveProgram(programName string) error {
	if db.GlobalCouchDB == nil {
		return errors.New("couch db instance is nil")
	}

	s := db.GlobalCouchDB.Use(GlobalConfig.CouchDB.DBName)
	doc := &Program{
		Type: "program",
		Tags: make(Tags),
	}
	doc.SetID(programName)
	doc.Tags.AddTag("programtype", "bugbounty")
	_, err := s.Store(doc)
	return err
}
