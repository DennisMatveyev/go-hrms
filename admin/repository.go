package admin

import (
	"errors"
	"hrms/common"
	"hrms/models"
	"log/slog"

	"gorm.io/gorm"
)

type AdminRepository interface {
	GetCandidates(status models.CandidateStatus) ([]*models.Candidate, error)
	GetCandidateByID(id int) (*CandidateResponse, error)
	UpdateCandidate(candidate *models.Candidate) (*models.Candidate, error)
}

type adminRepository struct {
	db  *gorm.DB
	log *slog.Logger
}

func NewAdminRepository(db *gorm.DB, log *slog.Logger) AdminRepository {
	return &adminRepository{db: db, log: log}
}

func (r *adminRepository) GetCandidates(status models.CandidateStatus) ([]*models.Candidate, error) {
	var candidates []*models.Candidate
	err := r.db.Where("status = ?", status).Order("updated_at DESC").Find(&candidates).Error
	if err != nil {
		r.log.Error("Database error when getting candidates", "error", err.Error())
		return nil, common.ErrDatabase
	}
	return candidates, nil
}

func (r *adminRepository) GetCandidateByID(id int) (*CandidateResponse, error) {
	var candidate models.Candidate
	err := r.db.First(&candidate, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Error("Candidate not found", "id", id)
			return nil, ErrCandidateNotFound
		}
		r.log.Error("Database error when getting candidate by ID", "error", err.Error())
		return nil, common.ErrDatabase
	}
	candidateResponse := &CandidateResponse{
		Candidate: candidate,
		Profile:   &models.Profile{},
		Jobs:      []*models.Job{},
	}
	var profile models.Profile
	if err := r.db.Where("user_id = ?", candidate.UserID).First(&profile).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Error("Database error when getting candidate profile", "error", err.Error())
		}
	} else {
		candidateResponse.Profile = &profile
	}
	if err := r.db.Where("user_id = ?", candidate.UserID).Find(&candidateResponse.Jobs).Error; err != nil {
		r.log.Error("Error loading jobs", "error", err.Error())
	}

	return candidateResponse, nil
}

func (r *adminRepository) UpdateCandidate(candidate *models.Candidate) (*models.Candidate, error) {
	err := r.db.Save(candidate).Error
	if err != nil {
		r.log.Error("Database error when updating candidate", "error", err.Error())
		return nil, common.ErrDatabase
	}
	return candidate, nil
}
