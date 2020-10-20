package models

import (
	"database/sql"
)


// Search поиск
func Search(db *sql.DB, user User, name *string) (*User, error) {
	p := user

	inquiry := "%" + *name + "%"
	err := db.QueryRow("SELECT `people`.`login` FROM `people` WHERE `login` LIKE ?", inquiry).Scan(&p.Login)
	return &p, err
}
