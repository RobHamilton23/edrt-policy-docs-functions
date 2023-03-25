package policyDocsUpdate

import (
	"context"
	"strings"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"pantheon.io/edrt-policy-docs-functions/config"
	firestore "pantheon.io/edrt-policy-docs-functions/internal/store"
	updateFunction "pantheon.io/edrt-policy-docs-functions/internal/updateFunction"
)

var logger *logrus.Logger

func init() {
	configMap := config.GetConfig()
	viper.MergeConfigMap(configMap)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	fs, err := firestore.NewFirestoreClient(
		context.Background(),
		logger,
		viper.GetString("firestore-project"),
	)

	if err != nil {
		// Yes, we want the process to die here. If we can't create the
		// firestore client, this cloud function cannot run and should
		// restart.
		logger.Fatalf("Unable to create firestore client: %w", err)
	}

	updateHandler := updateFunction.NewUpdateHandler(logger, fs)
	functions.CloudEvent("PolicyDocUpdated", updateHandler.PolicyDocUpdated)
}
