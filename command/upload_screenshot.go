package command

import (
	"errors"
	"io/ioutil"
	"mime"
	"os"
	"path"

	"github.com/hudangwei/couchdb"
	"github.com/openbugbounty/mojin-saver/db"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// DefaultIgnorePatterns contains the default list of glob patterns
// that are ignored when building a document from a directory.
var DefaultIgnorePatterns = []string{
	"*~", // editor swap files
	".*", // hidden files
	"_*", // CouchDB system fields
}

var uploadScreenshotCmd = &cobra.Command{
	Use:   "upload_screenshot",
	Short: "上传截图",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		err := runUploadScreenshotCmd(cmd, args)
		if err != nil {
			logrus.WithError(err).Error("runUploadScreenshotCmd")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(uploadScreenshotCmd)
}

func runUploadScreenshotCmd(_ *cobra.Command, args []string) error {
	return StoreIndependAttachments(args[0], args[1], nil)
}

func StoreIndependAttachments(
	docid, dir string,
	ignores []string,
) error {
	if db.GlobalCouchDB == nil {
		return errors.New("couch db instance is nil")
	}
	s := db.GlobalCouchDB.Use(GlobalConfig.CouchDB.DBName)
	var newrev string
	rev, err := s.Rev(docid)
	if err == nil {
		newrev = rev
	}
	err = walk(dir, ignores, func(p string) error {
		att := &couchdb.IndependAttachment{
			Name: path.Base(p),
			Type: mime.TypeByExtension(path.Ext(p)),
		}
		if att.Body, err = os.Open(p); err != nil {
			return err
		}
		resp, err := s.PutIndependAttachment(docid, att, newrev)
		if err != nil {
			return err
		}
		if resp.Ok {
			newrev = resp.Rev
		}
		return nil
	})
	return err
}

type walkFunc func(path string) error

func walk(dir string, ignores []string, callback walkFunc) error {
	if ignores == nil {
		ignores = DefaultIgnorePatterns
	}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, info := range files {
		isDir := info.IsDir()
		if isDir {
			continue
		}
		subpath := path.Join(dir, info.Name())
		// skip ignored files
		for _, pat := range ignores {
			if ign, err := path.Match(pat, info.Name()); err != nil {
				return err
			} else if ign {
				goto next
			}
		}

		if err := callback(subpath); err != nil {
			return err
		}
	next:
	}
	return nil
}
