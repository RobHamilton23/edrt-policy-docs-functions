package service

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	firestore "pantheon.io/edrt-policy-docs-functions/internal/store"
	"pantheon.io/edrt-policy-docs-functions/internal/types"
)

// DocumentTransformation provides functionality for transforming policy docs
type DocumentTransformation struct {
	firestore *firestore.Firestore
	logger    *logrus.Logger
}

// NewDocumentTransformation initializes an instance of hte DocumentTransformation type
// with the provided firestore and logger instances.
func NewDocumentTransformation(firestore *firestore.Firestore, logger *logrus.Logger) DocumentTransformation {
	return DocumentTransformation{
		firestore: firestore,
		logger:    logger,
	}
}

// Denormalize loads the normalized policy doc for the given site, environment, and hostname,
// denormalizes the data, and writes it to firestore.
func (d *DocumentTransformation) Denormalize(
	ctx context.Context,
	siteId string,
	environment string,
	hostname string,
) error {
	// Read normalized document from firestore
	_, hm, el, err := d.firestore.GetNormalizedDocs(
		ctx,
		siteId,
		environment,
		hostname,
	)
	if err != nil {
		return fmt.Errorf("unable to load normalized policy docs: %w", err)
	}

	// Denormalize
	denormed := populateDenormalizedDocument(hm, el)

	// Write denormalized documents to firestore
	paths := []string{
		fmt.Sprintf("denormed/policydoc/%s/policydoc", denormed.Hostname),
	}
	err = d.firestore.WriteDenormalizedDocs(ctx, paths, denormed)
	if err != nil {
		return fmt.Errorf("unable to write denormalized docs: %w", err)
	}

	return nil
}

// populateDenormalizedDocument creates a Denormalized instance and populates it
// with the provided policy doc data.
func populateDenormalizedDocument(hm *types.HostnameMetadata, el *types.EdgeLogic) *types.Denormalized {
	denormed := &types.Denormalized{}
	denormed.Hostname = hm.Hostname
	denormed.Zone = hm.Zone
	denormed.RedirectTo = el.RedirectTo
	denormed.EnforceHttps = el.EnforceHTTPS
	denormed.Backend = el.Backend
	denormed.BuildId = el.BuildId
	denormed.Jurisdiction = el.Jurisdiction
	denormed.SiteId = hm.SiteId
	denormed.SiteEnv = hm.SiteEnv
	return denormed
}
