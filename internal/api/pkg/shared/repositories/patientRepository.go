package repositories

import (
	E "yenexpress/internal/api/database/entities"
)

type PatientRepository struct {
	BUserRepository[E.Patient]
}

func NewPatientRepository() *PatientRepository {
	return &PatientRepository{}
}
