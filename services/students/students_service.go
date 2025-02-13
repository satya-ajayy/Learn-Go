package students

import (
	// Go Internal Packages
	"context"
	"fmt"

	// Local Packages
	errors "learn-go/errors"
	models "learn-go/models/students"

	// External Packages
	"go.mongodb.org/mongo-driver/mongo"
)

type StudentsRepository interface {
	GetOneStudent(ctx context.Context, rollNo string) (*models.StudentModel, error)
	GetAllStudents(ctx context.Context) (*[]models.StudentModel, error)
	InsertStudent(ctx context.Context, student models.StudentModel) error
	UpdateStudent(ctx context.Context, rollNo string, updatedStudent models.StudentModel) error
	DeleteStudent(ctx context.Context, rollNo string) error
}

type StudentsService struct {
	studentsRepository StudentsRepository
}

func NewService(studentsRepository StudentsRepository) *StudentsService {
	return &StudentsService{studentsRepository: studentsRepository}
}

// GetAllStudents returns all the students details
func (s *StudentsService) GetAllStudents(ctx context.Context) (*[]models.StudentModel, error) {
	students, err := s.studentsRepository.GetAllStudents(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all students details due to :: %w", err)
	}
	return students, nil
}

// GetOneStudent returns the students details for the given rollNo
func (s *StudentsService) GetOneStudent(ctx context.Context, rollNo string) (*models.StudentModel, error) {
	student, err := s.studentsRepository.GetOneStudent(ctx, rollNo)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) || student == nil {
			return nil, errors.E(errors.NotFound, "student details not found")
		}
		return nil, fmt.Errorf("failed to get student details for rollNo :: %s due to :: %w", rollNo, err)
	}
	return student, nil
}

// InsertStudent inserts a new student into the database
func (s *StudentsService) InsertStudent(ctx context.Context, student models.StudentModel) error {
	err := s.studentsRepository.InsertStudent(ctx, student)
	if err != nil {
		return fmt.Errorf("failed to insert student due to :: %w", err)
	}
	return nil
}

// UpdateStudent updates the student details for the given rollNo
func (s *StudentsService) UpdateStudent(ctx context.Context, rollNo string, updatedStudent models.StudentModel) error {
	err := s.studentsRepository.UpdateStudent(ctx, rollNo, updatedStudent)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.E(errors.NotFound, "student details not found")
		}
		return fmt.Errorf("failed to update student details for rollNo :: %s due to :: %w", rollNo, err)
	}
	return nil
}

// DeleteStudent deletes the student details for the given rollNo
func (s *StudentsService) DeleteStudent(ctx context.Context, rollNo string) error {
	err := s.studentsRepository.DeleteStudent(ctx, rollNo)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.E(errors.NotFound, "student details not found")
		}
		return fmt.Errorf("failed to delete student details for rollNo :: %s due to :: %w", rollNo, err)
	}
	return nil
}
