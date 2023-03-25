package firestore

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func New(ctx context.Context, logger *logrus.Logger, projectId string) (*Firestore, error) {
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
	collection := f.firestoreClient.Collection(hostnameCollectionName)

	hostnameDoc, err := f.getHostnameDocumentBySiteEnv(t, collection, siteId, env, hostname)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch hostname document for %s: %s", hostname, err)
	}

	var h *types.Hostname
	hostnameDoc.DataTo(&h)
	return h, nil
}

func (f *Firestore) readHostnameMetadata(
	t *firestore.Transaction,
	siteId string,
	env string,
	hostname string,
) (*types.HostnameMetadata, error) {
	collection := f.firestoreClient.Collection(hostnameMetadataCollectionName)
	hostnameDocument, err := f.getHostnameDocumentBySiteEnv(
		t,
		collection,
		siteId,
		env,
		hostname)
	if err != nil {
		return nil, fmt.Errorf("unable to get hostname document for %s: %s", hostname, err)
	}
	var h *types.HostnameMetadata
	hostnameDocument.DataTo(&h)
	return h, nil
}

func (f *Firestore) readEdgeLogic(
	t *firestore.Transaction,
	siteId string,
	env string,
	hostname string,
) (*types.EdgeLogic, error) {
	collection := f.firestoreClient.Collection(edgeLogicCollectionName)

	hostnameDocument, err := f.getHostnameDocumentBySiteEnv(
		t,
		collection,
		siteId,
		env,
		hostname,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to get hostname focument for %s: %s", hostname, err)
	}

	var h *types.EdgeLogic
	hostnameDocument.DataTo(&h)
	return h, nil
}

func (*Firestore) getHostnameDocumentBySiteEnv(
	t *firestore.Transaction,
	rootCollection *firestore.CollectionRef,
	siteId string,
	env string,
	hostname string) (*firestore.DocumentSnapshot, error) {
	siteDocRef := rootCollection.Doc(siteId)
	siteDoc, err := t.Get(siteDocRef)

	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, fmt.Errorf("document not found for site %s %w", siteId, err)
		}

		return nil, fmt.Errorf("unable to load site document for %s: %w", siteId, err)
	}

	envCollection := siteDoc.Ref.Collection(env)
	hostnameDoc := envCollection.Doc(hostname)
	hostnameDocument, err := t.Get(hostnameDoc)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, fmt.Errorf("document not found for hostname %s", hostname)
		}

		return nil, fmt.Errorf("unable to load hostname document for %s: %w", hostname, err)
	}
	return hostnameDocument, nil
}

// parseFirestorePath parsed the given path and returns the path portion, the
// document name at the end of the path, and an error if the path is not valid
func parseFirestorePath(path string) (collectionPath string, docName string, err error) {
	if len(path) == 0 {
		return "", "", fmt.Errorf("parseFirestorePath path must not be empty")
	}

	if path[0] == '/' {
		path = strings.TrimLeft(path, "/")
	}

	splitPath := strings.Split(path, "/")
	if len(splitPath)%2 != 0 {
		return "", "", fmt.Errorf("path should have an even number of components: %s", path)
	}

	lastPathSeparatorIndex := strings.LastIndex(path, "/")
	collectionPath = path[0:lastPathSeparatorIndex]
	docName = path[lastPathSeparatorIndex+1:]
	return
}
