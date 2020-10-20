// Package settings настройка сервера
package settings

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/securecookie"
	"log"

	"github.com/gorilla/sessions"

	"database/sql"
	"net/http"
	"time"
)

// Server настройка сервера
func Server() (*gin.Engine, *http.Server) {
	router := gin.New() // New - выключает промежуточное ПО, Default - Включает промежуточное ПО

	settings := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return router, settings
}

// StoreSession создание и настройка сессии
func StoreSession() *sessions.CookieStore {
	sessionKey := sessions.NewCookieStore(securecookie.GenerateRandomKey(32))
	//sessionKey := sessions.NewCookieStore([]byte("SECRET_KEY"))
	sessionKey.Options = &sessions.Options{
		Domain:   "localhost",
		Path:     "/",
		Secure:   false,
		HttpOnly: true, // Защита от XSS атак, но это слабая защита , перекрывая доступ к кукам из JavaScript
		//SameSite: http.SameSiteStrictMode,  // Подделка межсайтовых запросов (CSRF) запрещает использовать эти куки на других сайтах
	}

	return sessionKey
}

// DataBaseOpen подключение к БД
func DataBaseOpen() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/accounting")
	ErrorPanic(err, "Ошибка в подключении к БД")
	return db
}

// ErrorPanic обрабатывает вызов функции panic
func ErrorPanic(err error, description string) {
	if err != nil {
		log.Println(description)
		panic(err)
	}
}

// ErrorFatal обрабатывает ошибку
func ErrorFatal(err error, description string) {
	if err != nil {
		log.Fatal(err, ": ", description)
	}
}
