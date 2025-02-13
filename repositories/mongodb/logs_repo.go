package mongodb

import (
	// Go Internal Packages
	"context"

	// Local Packages
	models "learn-go/models/logs"

	// External Packages
	"go.mongodb.org/mongo-driver/mongo"
)

type LogsRepository struct {
	client     *mongo.Client
	collection string
}

func NewLogsRepository(client *mongo.Client) *LogsRepository {
	return &LogsRepository{client: client, collection: "logs"}
}

func (r *LogsRepository) InsertLog(ctx context.Context, log models.LogModel) error {
	collection := r.client.Database("mybase").Collection(r.collection)
	_, err := collection.InsertOne(ctx, log)
	if err != nil {
		return err
	}
	return nil
}
