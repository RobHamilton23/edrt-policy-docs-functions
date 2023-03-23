package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/sirupsen/logrus"
)

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
