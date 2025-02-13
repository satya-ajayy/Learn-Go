package handlers

import (
	// Go Internal Packages
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	// Local Packages
	errors "learn-go/errors"
	models "learn-go/models/students"

	// External Packages
	"github.com/go-chi/chi/v5"
)

type StudentsService interface {
	GetOneStudent(context.Context, string) (*models.StudentModel, error)
	GetAllStudents(context.Context) (*[]models.StudentModel, error)
	InsertStudent(context.Context, models.StudentModel) error
	UpdateStudent(context.Context, string, models.StudentModel) error
	DeleteStudent(context.Context, string) error
}

type StudentsHandler struct {
	svc StudentsService
}

func NewStudentsHandler(svc StudentsService) *StudentsHandler {
	return &StudentsHandler{svc: svc}
}

func (a *StudentsHandler) GetAll(w http.ResponseWriter, r *http.Request) (response any, status int, err error) {
	students, err := a.svc.GetAllStudents(r.Context())
	if err == nil {
		return students, http.StatusOK, nil
	}
	return
}

func (a *StudentsHandler) GetOne(w http.ResponseWriter, r *http.Request) (response any, status int, err error) {
	rollNo := chi.URLParam(r, "rollNo")
	if rollNo == "" {
		return nil, http.StatusBadRequest, errors.EmptyParamErr("rollNo")
	}

	student, err := a.svc.GetOneStudent(r.Context(), rollNo)
	if err == nil {
		return student, http.StatusOK, nil
	}
	return
}

func (a *StudentsHandler) Insert(w http.ResponseWriter, r *http.Request) (response any, status int, err error) {
	var student models.StudentModel
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		return nil, http.StatusBadRequest, errors.InvalidBodyErr(err)
	}
	if err := student.Validate(); err != nil {
		return nil, http.StatusBadRequest, errors.ValidationFailedErr(err)
	}

	err = a.svc.InsertStudent(r.Context(), student)
	if err == nil {
		return student, http.StatusCreated, nil
	}
	return
}

func (a *StudentsHandler) Update(w http.ResponseWriter, r *http.Request) (response any, status int, err error) {
	rollNo := chi.URLParam(r, "rollNo")
	if rollNo == "" {
		return nil, http.StatusBadRequest, errors.EmptyParamErr("rollNo")
	}

	var updatedStudent models.StudentModel
	if err := json.NewDecoder(r.Body).Decode(&updatedStudent); err != nil {
		return nil, http.StatusBadRequest, errors.InvalidBodyErr(err)
	}

	if err := updatedStudent.Validate(); err != nil {
		return nil, http.StatusBadRequest, errors.ValidationFailedErr(err)
	}

	err = a.svc.UpdateStudent(r.Context(), rollNo, updatedStudent)
	if err == nil {
		return updatedStudent, http.StatusOK, nil
	}
	return
}

func (a *StudentsHandler) Delete(w http.ResponseWriter, r *http.Request) (response any, status int, err error) {
	rollNo := chi.URLParam(r, "rollNo")
	if rollNo == "" {
		return nil, http.StatusBadRequest, errors.EmptyParamErr("rollNo")
	}

	err = a.svc.DeleteStudent(r.Context(), rollNo)
	if err == nil {
		return map[string]string{"message": fmt.Sprintf("%s is deleted", rollNo)}, http.StatusOK, nil
	}
	return
}
