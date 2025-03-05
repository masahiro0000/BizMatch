package usecase

import (
	"github.com/masahiro0000/BizMatch/internal/domain"
)

type MatchUsecase struct {
	matchRepo domain.MatchRepository
}

func NewMatchUsecase(mr domain.MatchRepository) *MatchUsecase{
	return &MatchUsecase{
		matchRepo: mr,
	}
}

func (u *MatchUsecase) ListMatch(userID int64) ([]*domain.Match, error) {
	return u.matchRepo.GetMatchByUserID(userID)
}