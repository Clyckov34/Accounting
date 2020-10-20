package models

type DataAccount struct { // Для авторизации
	id       int
	login    string
	password string
}

type DataUser struct { // Для выгрузки данных
	Id     int
	Login  string
	Phone  string
	Year   string
}

type User struct {	// Для поиска
	Login string
}
