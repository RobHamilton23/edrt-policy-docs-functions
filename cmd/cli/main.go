package main

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"

	_ "embed"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"pantheon.io/edrt-policy-docs-functions/config"
	firestore "pantheon.io/edrt-policy-docs-functions/internal/store"
)

func init() {
	configMap := config.GetConfig()
	viper.MergeConfigMap(configMap)
	fmt.Println(viper.GetString("FIRESTORE_PROJECT"))
}

func main() {
	logger := logrus.New()
	filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		logger.Info(path)
		return nil
	})

	fmt.Println("This will be a CLI allowing interactive testing of logic in this repo during development.")
	fmt.Println("Will probably use a library like cobra to set up the CLI.")
	f := firestore.New(logger, viper.GetString("FIRESTORE_PROJECT"))
	result := f.Read(context.Background())
	fmt.Println(result)
}
