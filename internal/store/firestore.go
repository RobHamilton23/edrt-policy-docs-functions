package firestore

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/sirupsen/logrus"
	"pantheon.io/edrt-policy-docs-functions/internal/types"
)

const hostnameCollectionName = "hostnames"
const hostnameMetadataCollectionName = "hostnameMetadata"
const edgeLogicCollectionName = "edgelogic"

type Firestore struct {
	logger    *logrus.Logger
	projectId string
}

func New(logger *logrus.Logger, projectId string) Firestore {
	return Firestore{
		logger:    logger,
		projectId: projectId,
	}
}

func (f *Firestore) Read(ctx context.Context) interface{} {
	client, err := firestore.NewClient(ctx, f.projectId)
	if err != nil {
		f.logger.Fatalf("Unable to create firestore client %s", err)
	}

	collection := client.Collection("foo")
	doc := collection.Doc("bar")
	data, err := doc.Get(ctx)
	if err != nil {
		f.logger.Fatalf("Unable to get document %s", err)
	}

	theData := data.Data()["baz"]
	switch stringData := theData.(type) {
	case string:
		return stringData
	default:
		f.logger.Fatal("Unable to read data")
	}

	return nil
}

func (f *Firestore) ReadHostname(ctx context.Context, siteId string, env string, hostname string) (*types.Hostname, error) {
	logger := f.logger.WithField("method", "ReadHostname")
	collection, err := f.getCollection(ctx, hostnameCollectionName, logger)
	if err != nil {
		return nil, fmt.Errorf("unable to get collection %s: %s", hostnameCollectionName, err)
	}

	hostnameDoc, err := f.getHostnameDocumentBySiteEnv(ctx, collection, siteId, env, hostname)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch hostname document for %s: %s", hostname, err)
	}

	var h *types.Hostname
	hostnameDoc.DataTo(&h)
	return h, nil
}

func (f *Firestore) ReadHostnameMetadata(ctx context.Context, siteId string, env string, hostname string) (*types.HostnameMetadata, error) {
	logger := f.logger.WithField("method", "ReadHostnameMetadata")

	collection, err := f.getCollection(ctx, hostnameMetadataCollectionName, logger)
	if err != nil {
		return nil, fmt.Errorf("unable to get collection %s: %s", hostnameMetadataCollectionName, err)
	}
	hostnameDocument, err := f.getHostnameDocumentBySiteEnv(
		ctx,
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

func (f *Firestore) ReadEdgeLogic(ctx context.Context, siteId string, env string, hostname string) (*types.EdgeLogic, error) {
	logger := f.logger.WithField("method", "ReadEdgeLogic")

	collection, err := f.getCollection(ctx, edgeLogicCollectionName, logger)
	if err != nil {
		return nil, fmt.Errorf("unable to get collection %s: %s", edgeLogicCollectionName, err)
	}

	hostnameDocument, err := f.getHostnameDocumentBySiteEnv(
		ctx,
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
	ctx context.Context,
	rootCollection *firestore.CollectionRef,
	siteId string,
	env string,
	hostname string) (*firestore.DocumentSnapshot, error) {
	siteDocRef := rootCollection.Doc(siteId)
	siteDoc, err := siteDocRef.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to load site document for %s: %s", siteId, err)
	}

	envCollection := siteDoc.Ref.Collection(env)
	hostnameDoc := envCollection.Doc(hostname)
	hostnameDocument, err := hostnameDoc.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to load hostname document for %s: %s", hostname, err)
	}
	return hostnameDocument, nil
}

func (f *Firestore) getCollection(ctx context.Context, collectionName string, logger *logrus.Entry) (*firestore.CollectionRef, error) {
	client, err := firestore.NewClient(ctx, f.projectId)
	if err != nil {
		return nil, fmt.Errorf("unable to create firestore client %s", err)
	}

	collection := client.Collection(collectionName)
	return collection, nil
}
