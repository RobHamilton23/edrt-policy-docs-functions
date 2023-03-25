package main

import (
	"context"
	"log"
	"strings"

	_ "embed"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"pantheon.io/edrt-policy-docs-functions/config"
	"pantheon.io/edrt-policy-docs-functions/internal/service"
	firestore "pantheon.io/edrt-policy-docs-functions/internal/store"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	log.SetOutput(logger.Writer())

	viper.MergeConfigMap(config.GetConfig())
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}

func main() {
	firestoreProject := viper.GetString("firestore-project")
	denormalizeCommand := &cobra.Command{
		Use:  "denormalize site_id env hostname",
		Args: cobra.MinimumNArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			f := getFirestore(firestoreProject)
			siteId := args[0]
			env := args[1]
			hostname := args[2]

			dts := service.NewDocumentTransformation(f, logger)
			err := dts.Denormalize(
				context.Background(),
				siteId,
				env,
				hostname,
			)
			if err != nil {
				logger.Fatalf("Unable to denormalize: %w", err)
			}
		},
	}

	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(denormalizeCommand)
	rootCmd.Execute()
}

func getFirestore(firestoreProject string) *firestore.Firestore {
	f, err := firestore.NewFirestoreClient(context.Background(), logger, firestoreProject)
	if err != nil {
		logger.Fatalf("Unable to create firestore instance: %w", err)
	}
	return f
}
