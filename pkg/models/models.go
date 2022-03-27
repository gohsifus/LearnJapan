package models

import (
	"database/sql"
)

var DB *sql.DB

type JpnCards struct{
	id int
	inJapan string
	inRussian string
	mark int
	dateAdd string
}
