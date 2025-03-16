package domain

import "errors"

type Prefectures struct {
	ID		int64	`json:"id" db:"id"`
	Name	string	`json:"name" db:"prefecture"`
}

type Industries struct {
	ID		int64	`json:"id" db:"id"`
	Name 	string  `json:"name" db:"industry"`
}

type Jobs struct {
	ID  	int64	`json:"id" db:"id"`
	Name	string	`json:"name" db:"job"`
}

type Positions struct {
	ID 		int64	`json:"id" db:"id"`
	Name 	string `json:"name" db:"position"`
}

var (
	ErrGetPrefecturesFailed = errors.New("都道府県情報の取得に失敗しました")
	ErrGetIndustriesFailed  = errors.New("業種情報の取得に失敗しました")
	ErrGetJobsFailed        = errors.New("職種情報の取得に失敗しました")
	ErrGetPositionsFailed   = errors.New("ポジション情報の取得に失敗しました")
)