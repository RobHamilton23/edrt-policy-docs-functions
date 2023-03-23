package policyDocsUpdate

import (
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/sirupsen/logrus"
	update "pantheon.io/edrt-policy-docs-functions/internal/updateFunction"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	updateHandler := update.NewUpdateHandler(logger)
	functions.CloudEvent("PolicyDocUpdated", updateHandler.PolicyDocUpdated)
}
