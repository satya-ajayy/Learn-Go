package students

import "learn-go/errors"

type StudentModel struct {
	RollNo string `json:"roll_no" bson:"Roll_No"`
	Name   string `json:"name" bson:"Student_Name"`
	Gender string `json:"gender" bson:"Gender"`
	MailID string `json:"mail_id" bson:"Mail_Id"`
}

func (s *StudentModel) Validate() error {
	ve := errors.ValidationErrs()
	if s.RollNo == "" {
		ve.Add("rollNo", "cannot be empty")
	}
	if s.Name == "" {
		ve.Add("name", "cannot be empty")
	}
	if s.Gender == "" {
		ve.Add("gender", "cannot be empty")
	}
	if s.MailID == "" {
		ve.Add("mail_id", "cannot be empty")
	}
	return ve.Err()
}
