package mongodb

import (
	// Go Internal Packages
	"context"

	// Local Packages
	models "learn-go/models/students"

	// External Packages
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StudentsRepository struct {
	client     *mongo.Client
	collection string
}

func NewStudentsRepository(client *mongo.Client) *StudentsRepository {
	return &StudentsRepository{client: client, collection: "class"}
}

// GetAllStudents returns all students in collection
func (r *StudentsRepository) GetAllStudents(ctx context.Context) (*[]models.StudentModel, error) {
	collection := r.client.Database("mybase").Collection(r.collection)
	findOptions := options.Find().SetSort(bson.D{{Key: "Roll_No", Value: 1}})

	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}

	students := []models.StudentModel{}
	if err := cursor.All(ctx, &students); err != nil {
		return nil, err
	}

	return &students, nil
}

// GetOneStudent returns a student with given rollNo
func (r *StudentsRepository) GetOneStudent(ctx context.Context, rollNo string) (*models.StudentModel, error) {
	collection := r.client.Database("mybase").Collection(r.collection)
	filter := bson.M{"Roll_No": rollNo}

	var student models.StudentModel
	err := collection.FindOne(ctx, filter).Decode(&student)
	if err != nil {
		return nil, err
	}
	return &student, nil
}

// InsertStudent inserts a students to the collection
func (r *StudentsRepository) InsertStudent(ctx context.Context, student models.StudentModel) error {
	collection := r.client.Database("mybase").Collection(r.collection)
	_, err := collection.InsertOne(ctx, student)
	if err != nil {
		return err
	}
	return nil
}

// UpdateStudent updates the student details with given rollNo
func (r *StudentsRepository) UpdateStudent(ctx context.Context, rollNo string, updatedStudent models.StudentModel) error {
	collection := r.client.Database("mybase").Collection(r.collection)
	filter := bson.M{"Roll_No": rollNo}
	res, err := collection.ReplaceOne(ctx, filter, updatedStudent)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

// DeleteStudent deletes a student with given rollNo
func (r *StudentsRepository) DeleteStudent(ctx context.Context, rollNo string) error {
	collection := r.client.Database("mybase").Collection(r.collection)
	filter := bson.M{"Roll_No": rollNo}
	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
