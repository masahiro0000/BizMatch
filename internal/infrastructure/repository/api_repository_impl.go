package repository

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/masahiro0000/BizMatch/internal/domain"
)

type apiRepositoryImpl struct {
	db *sqlx.DB
}

func NewApiRepositoryImpl(db *sqlx.DB) domain.ApiRepository{
	return &apiRepositoryImpl{
		db: db,
	}
}

// GetPrefectures retrieves the data from the prefectures table.
func (r *apiRepositoryImpl) GetPrefectures() ([]domain.Prefectures, error) {
	var prefectures []domain.Prefectures

	query := "SELECT id, prefecture FROM prefectures"
	if err := r.db.Select(&prefectures, query); err != nil {
		log.Printf("Fail to retrieve prefecture data. Query:%q, error:%v", query, err)
		return nil, err
	}
	return prefectures, nil
}

// GetIndustries retrieves the data from the industries table.
func (r *apiRepositoryImpl) GetIndustries() ([]domain.Industries, error) {
	var industries []domain.Industries

	query := "SELECT id, industry FROM industries"
	if err := r.db.Select(&industries, query); err != nil {
		log.Printf("Fail to retrieve industry data. Query:%q, error:%v", query, err)
		return nil, err
	}
	return industries, nil
}

// GetJobs retrieves the data from the jobs table.
func (r *apiRepositoryImpl) GetJobs() ([]domain.Jobs, error) {
	var jobs []domain.Jobs

	query := "SELECT id, job FROM jobs"
	if err := r.db.Select(&jobs, query); err != nil {
		log.Printf("Fail to retrieve job data. Query:%q, error:%v", query, err)
		return nil, err
	}
	return jobs, nil
}

// GetPositions retrieves the data from the positions table.
func (r *apiRepositoryImpl) GetPositions() ([]domain.Positions, error) {
	var positions []domain.Positions

	query := "SELECT id, position FROM positions"
	if err := r.db.Select(&positions, query); err != nil {
		log.Printf("Fail to retrieve position data. Query:%q, error:%v", query, err)
		return nil, err
	}
	return positions, nil
}