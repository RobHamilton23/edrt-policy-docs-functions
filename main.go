package policyDocsUpdate

import (
	"strings"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"pantheon.io/edrt-policy-docs-functions/config"
	firestore "pantheon.io/edrt-policy-docs-functions/internal/store"
	update "pantheon.io/edrt-policy-docs-functions/internal/updateFunction"
)

var logger *logrus.Logger

func init() {
	configMap := config.GetConfig()
	viper.MergeConfigMap(configMap)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	fs := firestore.New(logger, viper.GetString("firestore-project"))

	updateHandler := update.NewUpdateHandler(logger, &fs)
	functions.CloudEvent("PolicyDocUpdated", updateHandler.PolicyDocUpdated)
}
