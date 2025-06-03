package users

import (
	"errors"
	"log/slog"

	"hrms/common"
	"hrms/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	SaveUser(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	GetProfile(userID int) (*models.Profile, error)
	GetJobs(userID int) ([]*models.Job, error)
	CreateProfile(userID int, profile *models.Profile) (int, error)
	CreateJob(userID int, job *models.Job) (int, error)
	UpdateProfile(profile *models.Profile) error
	UpdateJob(job *models.Job) error
	DeleteJob(jobID int) error
}

type userRepository struct {
	db  *gorm.DB
	log *slog.Logger
}

func NewUserRepository(db *gorm.DB, log *slog.Logger) UserRepository {
	return &userRepository{db: db, log: log}
}

func (r *userRepository) SaveUser(user *models.User) error {
	tx := r.db.Begin()

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		r.log.Error("Database error when creating user", "error", err.Error())
		return common.ErrDatabase
	}
	candidate := models.Candidate{
		UserID: user.ID,
		Email:  user.Email,
	}
	if err := tx.Create(&candidate).Error; err != nil {
		tx.Rollback()
		r.log.Error("Database error when creating candidate", "error", err.Error())
		return err
	}

	return tx.Commit().Error
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		r.log.Info("User not found", "email", email)
		return nil, nil
	}
	if err != nil {
		r.log.Error("Database error when finding user", "error", err.Error())
		return nil, common.ErrDatabase
	}
	return &user, nil
}

func (r *userRepository) GetProfile(userID int) (*models.Profile, error) {
	var profile models.Profile
	err := r.db.Where("user_id = ?", userID).First(&profile).Error
	if err == nil || errors.Is(err, gorm.ErrRecordNotFound) {
		return &profile, nil
	} else {
		r.log.Error("Database error when getting profile", "error", err.Error())
		return nil, common.ErrDatabase
	}
}

func (r *userRepository) GetJobs(userID int) ([]*models.Job, error) {
	var jobs []*models.Job
	err := r.db.Where("user_id = ?", userID).Find(&jobs).Error
	if err != nil {
		r.log.Error("Database error when finding jobs", "error", err.Error())
		return nil, common.ErrDatabase
	}
	return jobs, nil
}

func (r *userRepository) CreateProfile(userID int, profile *models.Profile) (int, error) {
	profileFound := new(models.Profile)
	r.db.Where("user_id = ?", userID).First(&profileFound)
	if profileFound.ID != 0 {
		r.log.Error("Profile already exists", "userID", userID)
		return 0, ErrProfileExists
	}
	profile.UserID = userID
	err := r.db.Create(profile).Error
	if err != nil {
		r.log.Error("Database error when creating profile", "error", err.Error())
		return 0, common.ErrDatabase
	}
	return profile.ID, nil
}

func (r *userRepository) CreateJob(userID int, job *models.Job) (int, error) {
	job.UserID = userID
	err := r.db.Create(job).Error
	if err != nil {
		r.log.Error("Database error when creating job", "error", err.Error())
		return 0, common.ErrDatabase
	}
	return job.ID, nil
}

func (r *userRepository) UpdateProfile(profile *models.Profile) error {
	err := r.db.Save(profile).Error
	if err != nil {
		r.log.Error("Database error when updating profile", "error", err.Error())
		return common.ErrDatabase
	}
	return nil
}

func (r *userRepository) UpdateJob(job *models.Job) error {
	err := r.db.Save(job).Error
	if err != nil {
		r.log.Error("Database error when updating job", "error", err.Error())
		return common.ErrDatabase
	}
	return nil
}

func (r *userRepository) DeleteJob(jobID int) error {
	err := r.db.Delete(&models.Job{}, jobID).Error
	if err != nil {
		r.log.Error("Database error when deleting job", "error", err)
		return common.ErrDatabase
	}
	return nil
}
