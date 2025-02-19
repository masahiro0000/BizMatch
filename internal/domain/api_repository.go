package domain

type ApiRepository interface {
	GetPrefectures()([]Prefectures, error)
	GetIndustries()([]Industries, error)
	GetJobs()([]Jobs, error)
	GetPositions()([]Positions, error)
}