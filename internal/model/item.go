package model

type Item struct {
	Id          int64   `json:"id" db:"id" goku:"skipinsert"`
	UserId      int64   `json:"user_id" db:"user_id"`
	Name        string  `json:"name" db:"name"`
	Url         string  `json:"url" db:"url"`
	Category    int     `json:"category" db:"category"`
	SizeNumber  *int    `json:"size_number" db:"size_number"`
	SizeText    *string `json:"size_text" db:"size_text"`
	Description *string `json:"description" db:"description"`
	Color       int     `json:"color" db:"color"`
	FileName    string  `json:"file_name" db:"file_name"`
	File        File    `json:"-" db:"-"`
}
