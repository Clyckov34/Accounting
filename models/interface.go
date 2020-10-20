package models

type DataAccount struct { // Для авторизации
	id       int
	login    string
	password string
	secret 	 string
}

type DataUser struct { // Для выгрузки данных
	Id     int
	Login  string
	Phone  string
	Year   string
	Secret string
}

type User struct {	// Для поиска
	Login string
}
