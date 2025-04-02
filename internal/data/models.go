package data

import "database/sql"

type Models struct {
	Rooms RoomModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Rooms: RoomModel{DB: db},
	}
}
