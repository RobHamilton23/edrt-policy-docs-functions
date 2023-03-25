package service

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	firestore "pantheon.io/edrt-policy-docs-functions/internal/store"
	"pantheon.io/edrt-policy-docs-functions/internal/types"
)

type DocumentTransformation struct {
	firestore *firestore.Firestore
	logger    *logrus.Logger
}

func NewDocumentTransformation(firestore *firestore.Firestore, logger *logrus.Logger) DocumentTransformation {
	return DocumentTransformation{
		firestore: firestore,
		logger:    logger,
	}
}

func (d *DocumentTransformation) Denormalize(
	ctx context.Context,
	siteId string,
	environment string,
	hostname string,
) error {
	// Read normalized document from firestore
	h, hm, el, err := d.firestore.GetNormalizedDocs(
		ctx,
		siteId,
		environment,
		hostname,
	)

	fmt.Println(h)
	fmt.Println(hm)
	fmt.Println(el)
	fmt.Println(err)

	// Denormalize
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

	// Write denormalizes documents to firestore
	paths := []string{
		fmt.Sprintf("denormed/policydoc/%s", denormed.Hostname),
	}
	err = d.firestore.WriteDenormalizedDocs(ctx, paths, denormed)
	if err != nil {
		return fmt.Errorf("unable to write denormalized docs: %w", err)
	}

	return nil
}
