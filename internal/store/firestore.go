package firestore

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/sirupsen/logrus"
	"pantheon.io/edrt-policy-docs-functions/internal/types"
)

const hostnameCollectionName = "hostnames"
const hostnameMetadataCollectionName = "hostnameMetadata"
const edgeLogicCollectionName = "edgelogic"

type Firestore struct {
	logger          *logrus.Logger
	projectId       string
	firestoreClient *firestore.Client
}

// NewFirestoreClient initializes a Firestore instance used for operations on
// policy docs in firestore
func NewFirestoreClient(ctx context.Context, logger *logrus.Logger, projectId string) (*Firestore, error) {
	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("unable to create firestore client: %w", err)
	}

	return &Firestore{
		logger:          logger,
		projectId:       projectId,
		firestoreClient: client,
	}, nil
}

// GetNormalizedDocs reads the normalized Hostname, Hostname Metadata, and
// Edge Logic documents from firestore
func (f *Firestore) GetNormalizedDocs(
	ctx context.Context,
	siteId string,
	environment string,
	hostname string,
) (h *types.Hostname, m *types.HostnameMetadata, el *types.EdgeLogic, err error) {
	err = f.firestoreClient.RunTransaction(
		ctx,
		func(ctx context.Context, t *firestore.Transaction) error {
			h, err = f.readHostname(
				t,
				siteId,
				environment,
				hostname,
			)
			if err != nil {
				return fmt.Errorf("unable to fetch hostname in transaction: %w: ", err)
			}

			m, err = f.readHostnameMetadata(
				t,
				siteId,
				environment,
				hostname,
			)
			if err != nil {
				return fmt.Errorf("unable to fetch hostname metadata in transaction: %w", err)
			}

			el, err = f.readEdgeLogic(
				t,
				siteId,
				environment,
				hostname,
			)
			if err != nil {
				return fmt.Errorf("unable to fetch edge logic in transaction: %w", err)
			}

			return nil
		},
	)

	return
}

// WriteDenormalizedDocs writes the denormalized document to the provided paths
func (f *Firestore) WriteDenormalizedDocs(
	ctx context.Context,
	paths []string,
	denormalizedDoc *types.Denormalized,
) error {
	logger := f.logger.WithField("method", "WriteDenormalizedDocs")
	return f.firestoreClient.RunTransaction(
		ctx,
		func(ctx context.Context, t *firestore.Transaction) error {
			for _, path := range paths {
				collectionPath, docName, err := parseFirestorePath(path)
				if err != nil {
					return err
				}

				collection := f.firestoreClient.Collection(collectionPath)
				docRef := collection.Doc(docName)
				err = t.Set(docRef, denormalizedDoc)
				if err != nil {
					return fmt.Errorf("unable to write denormed policydoc to %s: %w", path, err)
				}

				logger.WithField("path", path).Info("Write Complete")
			}

			return nil
		},
	)
}

func (f *Firestore) readHostname(
	t *firestore.Transaction,
	siteId string,
	env string,
	hostname string,
) (*types.Hostname, error) {
	pathComponents := []string{hostnameCollectionName, siteId, env}
	hostnameDocumentRef := f.firestoreClient.Collection(
		strings.Join(pathComponents, "/"),
	).Doc(hostname)

	hostnameDocument, err := t.Get(hostnameDocumentRef)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch hostname document for %s: %s", hostname, err)
	}

	var h *types.Hostname
	err = hostnameDocument.DataTo(&h)
	if err != nil {
		return nil, fmt.Errorf("unable to deserialize hostname document: %w", err)
	}
	return h, nil
}

func (f *Firestore) readHostnameMetadata(
	t *firestore.Transaction,
	siteId string,
	env string,
	hostname string,
) (*types.HostnameMetadata, error) {
	pathComponents := []string{hostnameMetadataCollectionName, siteId, env}
	hostnameMetadataDocumentRef := f.firestoreClient.Collection(
		strings.Join(pathComponents, "/"),
	).Doc(hostname)

	hostnameMetadataDocument, err := t.Get(hostnameMetadataDocumentRef)
	if err != nil {
		return nil, fmt.Errorf("unable to get hostname document for %s: %s", hostname, err)
	}
	var h *types.HostnameMetadata
	err = hostnameMetadataDocument.DataTo(&h)
	if err != nil {
		return nil, fmt.Errorf("unable to deserialize hostname metadata document: %w", err)
	}

	return h, nil
}

func (f *Firestore) readEdgeLogic(
	t *firestore.Transaction,
	siteId string,
	env string,
	hostname string,
) (*types.EdgeLogic, error) {
	pathComponents := []string{edgeLogicCollectionName, siteId, env}
	edgeLogicDocumentRef := f.firestoreClient.Collection(
		strings.Join(pathComponents, "/"),
	).Doc(hostname)

	edgeLogicDocument, err := t.Get(edgeLogicDocumentRef)
	if err != nil {
		return nil, fmt.Errorf("unable to get hostname focument for %s: %s", hostname, err)
	}

	var h *types.EdgeLogic
	err = edgeLogicDocument.DataTo(&h)
	if err != nil {
		return nil, fmt.Errorf("unable to deserialize edgelogic document: %w", err)
	}

	return h, nil
}

// parseFirestorePath parsed the given path and returns the path portion, the
// document name at the end of the path, and an error if the path is not valid
func parseFirestorePath(path string) (collectionPath string, docName string, err error) {
	if len(path) == 0 {
		err = fmt.Errorf("parseFirestorePath path must not be empty")
		return
	}

	if path[0] == '/' {
		path = strings.TrimLeft(path, "/")
	}

	splitPath := strings.Split(path, "/")
	if len(splitPath)%2 != 0 {
		err = fmt.Errorf("path should have an even number of components: %s", path)
		return
	}

	lastPathSeparatorIndex := strings.LastIndex(path, "/")
	collectionPath = path[0:lastPathSeparatorIndex]
	docName = path[lastPathSeparatorIndex+1:]
	return
}
