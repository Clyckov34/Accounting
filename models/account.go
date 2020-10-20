package models

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)


// DBAccount активация аккаунта... Логин: admin Пароль: admin24
func DBAccount(db *sql.DB, account DataAccount, login, password string) error {
	query := db.QueryRow("SELECT `people`.`login`, `people`.`password`, `people`.`secret` FROM `people` WHERE `login` = ?", login).Scan(&account.login, &account.password, &account.secret)
	if query != nil {
		return errors.New("не верный логин")
	} else {
		err := bcrypt.CompareHashAndPassword([]byte(account.password), []byte(password)) // Хеширование - расшифровывает
		if err != nil {
			return errors.New("не верный пароль")
		} else {
			return nil
		}
	}
}

// DBExportAccountTable выгрузка из двух таблиц
func DBExportTable(db *sql.DB, user DataUser, login string) (*DataUser, error) {
	err := db.QueryRow("SELECT `people`.`id`, `people`.`login`, `people`.`secret`, `dates`.`phone`, `dates`.`year` FROM `people` INNER JOIN `dates` ON `people`.`login` = `dates`.`login` WHERE `people`.`login` = ?", login).Scan(&user.Id, &user.Login, &user.Secret ,&user.Phone, &user.Year)
	return &user, err
}
