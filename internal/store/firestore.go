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
			h, err = f.ReadHostname(
				ctx,
				t,
				siteId,
				environment,
				hostname,
			)
			if err != nil {
				return fmt.Errorf("unable to fetch hostname in transaction: %w: ", err)
			}

			m, err = f.ReadHostnameMetadata(
				ctx,
				t,
				siteId,
				environment,
				hostname,
			)
			if err != nil {
				return fmt.Errorf("unable to fetch hostname metadata in transaction: %w", err)
			}

			el, err = f.ReadEdgeLogic(
				ctx,
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
				if path[0] == '/' {
					path = strings.TrimLeft(path, "/")
				}
				err := assertValidFirestoreDocPath(path)
				if err != nil {
					return err
				}

				lastPathSeparatorIndex := strings.LastIndex(path, "/")
				collectionPath := path[0:lastPathSeparatorIndex]
				docName := path[lastPathSeparatorIndex+1:]
				collection := f.getCollection(ctx, collectionPath, logger)
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

func (f *Firestore) ReadHostname(
	ctx context.Context,
	t *firestore.Transaction,
	siteId string,
	env string,
	hostname string,
) (*types.Hostname, error) {
	logger := f.logger.WithField("method", "ReadHostname")
	collection := f.getCollection(ctx, hostnameCollectionName, logger)

	hostnameDoc, err := f.getHostnameDocumentBySiteEnv(t, collection, siteId, env, hostname)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch hostname document for %s: %s", hostname, err)
	}

	var h *types.Hostname
	hostnameDoc.DataTo(&h)
	return h, nil
}

func (f *Firestore) ReadHostnameMetadata(
	ctx context.Context,
	t *firestore.Transaction,
	siteId string,
	env string,
	hostname string,
) (*types.HostnameMetadata, error) {
	logger := f.logger.WithField("method", "ReadHostnameMetadata")

	collection := f.getCollection(ctx, hostnameMetadataCollectionName, logger)
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

func (f *Firestore) ReadEdgeLogic(
	ctx context.Context,
	t *firestore.Transaction,
	siteId string,
	env string,
	hostname string,
) (*types.EdgeLogic, error) {
	logger := f.logger.WithField("method", "ReadEdgeLogic")

	collection := f.getCollection(ctx, edgeLogicCollectionName, logger)

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

func (f *Firestore) getCollection(
	ctx context.Context,
	collectionName string,
	logger *logrus.Entry,
) *firestore.CollectionRef {
	collection := f.firestoreClient.Collection(collectionName)
	return collection
}

func assertValidFirestoreDocPath(path string) error {
	splitPath := strings.Split(path, "/")
	if len(splitPath)%2 != 0 {
		return fmt.Errorf("path should have an even number of components: %s", path)
	}

	return nil
}
