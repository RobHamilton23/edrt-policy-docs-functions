package policyDocsUpdate

import (
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/sirupsen/logrus"
	"pantheon.io/edrt-policy-docs-functions/config"
	firestore "pantheon.io/edrt-policy-docs-functions/internal/store"
	update "pantheon.io/edrt-policy-docs-functions/internal/updateFunction"
)

var logger *logrus.Logger

func init() {
	configMap := config.GetConfig()

	logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	fs := firestore.New(logger, configMap["FIRESTORE_PROJECT"].(string))

	updateHandler := update.NewUpdateHandler(logger, &fs)
	functions.CloudEvent("PolicyDocUpdated", updateHandler.PolicyDocUpdated)
}
