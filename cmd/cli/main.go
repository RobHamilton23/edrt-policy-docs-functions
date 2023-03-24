package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	_ "embed"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"pantheon.io/edrt-policy-docs-functions/config"
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
	getHostnameCommand := &cobra.Command{
		Use:   "get-hostname site_id env hostname",
		Short: "Fetches a hostname from the normalized collection",
		Long:  "Fetches a hostname from the normalized collection",
		Args:  cobra.MinimumNArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			f := firestore.New(logger, firestoreProject)
			siteId := args[0]
			env := args[1]
			hostname := args[2]
			result, err := f.ReadHostname(context.Background(), siteId, env, hostname)
			if err != nil {
				logger.Fatalf("unable to fetch hostname %s", err)
			}

			hostnameJson, _ := json.Marshal(result)
			fmt.Println(string(hostnameJson))
		},
	}

	getHostnameMetadataCommand := &cobra.Command{
		Use:   "get-hostname-metadata site_id env hostname",
		Short: "Fetches hostname metadata from the normalized collection",
		Long:  "Fetches hostname metadata from the normalized collection",
		Args:  cobra.MinimumNArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			f := firestore.New(logger, firestoreProject)
			siteId := args[0]
			env := args[1]
			hostname := args[2]
			result, err := f.ReadHostnameMetadata(context.Background(), siteId, env, hostname)
			if err != nil {
				logger.Fatalf("unable to fetch hostname %s", err)
			}

			hostnameMetadataJson, _ := json.Marshal(result)
			fmt.Println(string(hostnameMetadataJson))
		},
	}

	getEdgeLogicCommand := &cobra.Command{
		Use:   "get-edge-logic site_id env hostname",
		Short: "Fetches edge logic from the normalized collection",
		Long:  "Fetches edge logic from the normalized collection",
		Args:  cobra.MinimumNArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			f := firestore.New(logger, firestoreProject)
			siteId := args[0]
			env := args[1]
			hostname := args[2]
			result, err := f.ReadEdgeLogic(context.Background(), siteId, env, hostname)
			if err != nil {
				logger.Fatalf("unable to fetch hostname %s", err)
			}

			edgeLogicJson, _ := json.Marshal(result)
			fmt.Println(string(edgeLogicJson))
		},
	}

	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(getHostnameCommand)
	rootCmd.AddCommand(getHostnameMetadataCommand)
	rootCmd.AddCommand(getEdgeLogicCommand)
	rootCmd.Execute()
}
