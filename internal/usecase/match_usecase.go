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
	var MatchList []*domain.Match
	// Get all matches for the user.
	MatchList, err := u.matchRepo.GetMatchByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Filter out blocked matches.
	var filteredMatchList []*domain.Match
	for _, match := range MatchList {
		if match.Status != "BLOCKED" {
			filteredMatchList = append(filteredMatchList, match)
		}
	}
	return filteredMatchList, nil
}

func (u *MatchUsecase) GetMatchByID(matchID int64) (*domain.Match, error) {
	return u.matchRepo.GetMatchByID(matchID)
}

func (u *MatchUsecase) UpdateMatchStatus(matchID int64, status string) error {
	return u.matchRepo.UpdateMatchStatus(matchID, "BLOCKED")
}