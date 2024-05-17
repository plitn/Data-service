package model

type Look struct {
	Id        int64  `db:"id" json:"id"`
	UserId    int64  `db:"user_id" json:"user_id"`
	LookName  string `db:"look_name" json:"look_name"`
	StylistId *int64 `db:"stylist_id" json:"stylist_id"`
}

type LooksItems struct {
	LookId int64 `db:"look_id" json:"look_id"`
	ItemId int64 `db:"item_id" json:"item_id"`
}

type LookData struct {
	Id        int64  `db:"id" json:"id"`
	UserId    int64  `db:"user_id" json:"user_id"`
	LookName  string `db:"look_name" json:"look_name"`
	StylistId *int64 `db:"stylist_id" json:"stylist_id"`
	Items     []Item `json:"items"`
}

type Capsule struct {
	Id     int64  `db:"id" json:"id" goku:"skipinsert"`
	Name   string `db:"name" json:"name"`
	UserId int64  `db:"user_id" json:"user_id"`
}

type ItemCapsule struct {
	Id        int64 `db:"id" json:"id" goku:"skipinsert"`
	ItemId    int64 `db:"item_id" json:"item_id"`
	CapsuleId int64 `db:"capsule_id" json:"capsule_id"`
}
