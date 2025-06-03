package admin

import "hrms/models"

type CandidateResponse struct {
	models.Candidate
	Profile *models.Profile `json:"profile"`
	Jobs    []*models.Job   `json:"jobs"`
}
