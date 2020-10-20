package models

import (
	"crypto/sha512"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"time"
)


// DBAccount активация аккаунта... Логин: admin Пароль: admin24
func DBAccount(db *sql.DB, account DataAccount, login, password string) error {
	query := db.QueryRow("SELECT `people`.`login`, `people`.`password` FROM `people` WHERE `login` = ?", login).Scan(&account.login, &account.password)
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
	err := db.QueryRow("SELECT `people`.`id`, `people`.`login`, `dates`.`phone`, `dates`.`year` FROM `people` INNER JOIN `dates` ON `people`.`login` = `dates`.`login` WHERE `people`.`login` = ?", login).Scan(&user.Id, &user.Login, &user.Phone, &user.Year)
	return &user, err
}


// Sha256 генерирует аунтификационный код
func Sha512() string {
	h := time.Now().String()
	sha := sha512.New()
	sha.Write([]byte(h))
	return fmt.Sprintf("%x", sha.Sum(nil))
}