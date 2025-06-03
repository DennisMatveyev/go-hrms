package models

import "time"

type User struct {
	ID        int       `json:"id,omitzero" gorm:"primaryKey"`
	Email     string    `json:"email" validate:"required,email" gorm:"unique;uniqueIndex;type:varchar(255)"`
	Password  string    `json:"-" validate:"-" gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `json:"created_at,omitzero" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at,omitzero" gorm:"autoUpdateTime"`

	Candidate *Candidate `json:"-" gorm:"foreignKey:UserID"`
	Profile   *Profile   `json:"-" gorm:"foreignKey:UserID"`
	Jobs      []Job      `json:"-" gorm:"foreignKey:UserID"`
}

type Profile struct {
	ID          int    `json:"id,omitempty" gorm:"primaryKey"`
	UserID      int    `json:"user_id"`
	FirstName   string `json:"first_name" validate:"required,min=2,max=32"`
	LastName    string `json:"last_name" validate:"required,min=2,max=32"`
	PhoneNumber string `json:"phone_number" validate:"required,min=9,max=15"`
	PhotoPath   string `json:"-"`
	Header      string `json:"header" validate:"required,min=10,max=200"`
}

type Job struct {
	ID          int        `json:"id,omitzero" gorm:"primaryKey"`
	UserID      int        `json:"user_id"`
	CompanyName string     `json:"company_name"`
	Position    string     `json:"position"`
	StartDate   time.Time  `json:"start_date" gorm:"type:date"`
	EndDate     *time.Time `json:"end_date" gorm:"type:date"`
	Description string     `json:"description"`
}

type Candidate struct {
	ID        int             `json:"id" gorm:"primaryKey"`
	UserID    int             `json:"user_id" gorm:"index"`
	Status    CandidateStatus `json:"status" gorm:"default:'pending'"`
	Notes     string          `json:"notes" gorm:"type:text"`
	CreatedAt time.Time       `json:"created_at,omitzero" gorm:"autoCreateTime"`
	UpdatedAt time.Time       `json:"updated_at,omitzero" gorm:"autoUpdateTime"`
	Email     string          `json:"email"`
}

type CandidateStatus string

const (
	StatusPending  CandidateStatus = "pending"
	StatusProgress CandidateStatus = "progress"
	StatusApproved CandidateStatus = "approved"
	StatusRejected CandidateStatus = "rejected"
)
