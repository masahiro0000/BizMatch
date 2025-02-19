package usecase

import (
	"log"

	"github.com/masahiro0000/BizMatch/internal/domain"
)

type ApiUsecase struct {
	apiRepo domain.ApiRepository
}

func NewApiUsecase(repo domain.ApiRepository) *ApiUsecase {
	return &ApiUsecase{
		apiRepo: repo,
	}
}

// GetPrefectures retrieves prefecture data by calling the repository method.
func (u *ApiUsecase) GetPrefectures() ([]domain.Prefectures, error){
	prefectures, err := u.apiRepo.GetPrefectures()
	if err != nil {
		log.Printf("Fail to fetch prefectures from repository: %v", err)
		return nil, err
	}
	return prefectures, nil
}

// GetIndustries retrieves prefecture data by calling the repository method.
func (u *ApiUsecase) GetIndustries() ([]domain.Industries, error) {
	industries, err := u.apiRepo.GetIndustries()
	if err != nil {
		log.Printf("Fail to fetch industries from repository: %v", err)
		return nil, err
	}
	return industries, nil
}

// GetJobs retrieves prefecture data by calling the repository method.
func (u *ApiUsecase) GetJobs() ([]domain.Jobs, error) {
	jobs, err := u.apiRepo.GetJobs()
	if err != nil {
		log.Printf("Fail to fetch jobs from repository: %v", err)
		return nil, err
	}
	return jobs, nil
}

// GetPositions retrieves prefecture data by calling the repository method.
func (u *ApiUsecase) GetPositions() ([]domain.Positions, error) {
	positions, err := u.apiRepo.GetPositions()
	if err != nil {
		log.Printf("Fail to fetch positions from repository: %v", err)
		return nil, err
	}
	return positions, nil
}