package domain

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