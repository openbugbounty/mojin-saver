package command

import (
	"fmt"
	"os"

	"github.com/hudangwei/couchdb"
	"github.com/joho/godotenv"
	"github.com/openbugbounty/mojin-saver/config"
	"github.com/openbugbounty/mojin-saver/db"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var GlobalConfig config.Config
var envfile string
var program string

func init() {
	rootCmd.PersistentFlags().StringVarP(&envfile, "env-file", "", ".env", "Read in a file of environment variables")
	rootCmd.PersistentFlags().StringVarP(&program, "program", "", "default", "项目名") //TODO:@INFER 根据项目inscope自动识别项目

	cobra.OnInitialize(
		initConfig,
		initDB,
	)
}

func initConfig() {
	godotenv.Load(envfile)
	config, err := config.Environ()
	if err != nil {
		logger := logrus.WithError(err)
		logger.Fatalln("main: invalid configuration")
	}
	GlobalConfig = config

	initLogging(config)

	if logrus.IsLevelEnabled(logrus.TraceLevel) {
		fmt.Println(GlobalConfig.String())
	}
}

func initLogging(c config.Config) {
	if c.Logging.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if c.Logging.Trace {
		logrus.SetLevel(logrus.TraceLevel)
	}
	if c.Logging.Text {
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors:   c.Logging.Color,
			DisableColors: !c.Logging.Color,
		})
	} else {
		logrus.SetFormatter(&logrus.JSONFormatter{
			PrettyPrint: c.Logging.Pretty,
		})
	}
}

func initDB() {
	db.GlobalCouchDB = couchdb.NewClient("http://"+GlobalConfig.CouchDB.Host, GlobalConfig.CouchDB.User, GlobalConfig.CouchDB.Password)
	s := db.GlobalCouchDB.Use(GlobalConfig.CouchDB.DBName)
	if designDocs, err := db.GlobalCouchDB.Parse(GlobalConfig.CouchDB.DesignViewImportDir); err == nil && len(designDocs) > 0 {
		if err := s.Seed(designDocs); err != nil {
			logger := logrus.WithError(err)
			logger.Warnln("同步couchdb design view错误")
		}
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
