package routers

import (
	"GIN/models"
	"GIN/settings"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RouterUrl(engine *gin.Engine){

	var store = settings.StoreSession()

	// Главная страница
	engine.GET("/", func(context *gin.Context) {
		context.HTML(200, "index.html", gin.H{
			"title": "Aвторизация",
		})
	})

	// Авторизация аккаунта + профиль
	account := engine.Group("/account")
	{
		// Страница с Авторизацией форма HTML
		account.GET("/", func(context *gin.Context) {
			session, _ := store.Get(context.Request, "secret")
			_, resOk := session.Values["login"]

			context.HTML(200, "accounting.html", gin.H{
				"title" : "Авторизация аккаунта",
				"login" : resOk,
			})
		})

		// Обработчик Авторизации
		account.POST("/Treatment", func(context *gin.Context) {
			login, loginOK := context.GetPostForm("login") // Принимает POST информацию + проверяет на пустые данных
			password, passwordOK := context.GetPostForm("password")

			if loginOK == true && passwordOK == true {
				db := settings.DataBaseOpen() // Подключение к БД
				defer db.Close()

				result := models.DBAccount(db, models.DataAccount{}, login, password)
				if result != nil {
					context.Redirect(302, "/")
				} else {
					secretSHA512 := models.Sha512()

					session, _ := store.Get(context.Request, "secret")
					session.Values["login"] = login
					session.Values["token"] = secretSHA512
					_ = session.Save(context.Request, context.Writer)


					context.Redirect(302, "/account/user")
				}
			} else {
				context.String(200, "Некоторые поля пусты")
			}
		})

		// Профель пользователя
		account.GET("/user", func(context *gin.Context) {
			session, err := store.Get(context.Request, "secret")
			if err != nil {
				context.Redirect(302, "/")
			}

			resLogin, okLogin := session.Values["login"]
			_, okToken := session.Values["token"]

			if okLogin != true && okToken != true {
				context.Redirect(302, "/")
			} else {
				db := settings.DataBaseOpen() // Подключение к БД
				defer db.Close()

				result, err := models.DBExportTable(db, models.DataUser{}, fmt.Sprint(resLogin)) // Выгрузка данных о пользователи... Логин, телефон, год и т.д.
				if err != nil {                                     // Если ответ пришел пустой, то значит это ошибка
					context.String(200, "Нет в наличии")
				} else { // Если нет ошибок, то есть данные
					context.HTML(200, "secret.html", gin.H{
						"id":    result.Id,
						"login": result.Login,
						"phone": result.Phone,
						"year":  result.Year,
					})
				}
			}
		})

		// Удаление аккаунта + удаляет session
		account.GET("/delete/:id", func(context *gin.Context) {
			delete := context.Param("id")

			db := settings.DataBaseOpen()
			defer db.Close()

			_, err := db.Exec("DELETE FROM `people` WHERE `people`.`id`= ?", delete)
			_, err = db.Exec("DELETE FROM `dates` WHERE `dates`.`id` = ?", delete)

			if err != nil {
				context.String(200, "Нет такого аккаунта")
			} else {
				session, err := store.Get(context.Request, "secret")
				if err != nil {
					context.Redirect(302, "/account")
				}

				_, okLogin := session.Values["login"]
				_, okToken:= session.Values["token"]

				if okLogin != true && okToken != true {
					context.Redirect(302, "/account")
				} else {
					session.Options.MaxAge = -1
					_ = session.Save(context.Request, context.Writer)

					context.Redirect(302, "/account")
				}
			}
		})

		// Удаление Session
		account.GET("/session", func(context *gin.Context) {
			session, err := store.Get(context.Request, "secret")
			if err != nil {
				context.Redirect(302, "/account")
			}

			_, okLogin := session.Values["login"]
			_, okToken:= session.Values["token"]

			if okLogin != true && okToken != true {
				context.Redirect(302, "/account")
			} else {
				session.Options.MaxAge = -1
				_ = session.Save(context.Request, context.Writer)

				context.Redirect(302, "/account")
			}
		})
	}

	// Регистрация профиля
	accountRegistr := engine.Group("/accountreg")
	{
		// Регистрация аккаунта Html шаблон
		accountRegistr.GET("/", func(context *gin.Context) {
			context.HTML(200, "accountingReg.html", gin.H{
				"title" : "Регистрация аккаунта",
			})
		})

		// регистрация нового пользователя
		accountRegistr.POST("/user", func(context *gin.Context) {
			login, loginOK := context.GetPostForm("login")
			password, passwordOK := context.GetPostForm("password")
			phone, phoneOK := context.GetPostForm("phone")
			year, yearOK := context.GetPostForm("year")

			if loginOK && passwordOK && phoneOK && yearOK == true {
				db := settings.DataBaseOpen()
				defer db.Close()

				password, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
				if err != nil {
					context.String(200, "Ошибка Хеширование пароля")
				} else {
					_, err = db.Exec("INSERT INTO `people` (`login`, `password`) VALUES (?, ?)", login, password)
					_, err = db.Exec("INSERT INTO `dates` (`login`, `phone`, `year`) VALUES (?, ?, ?)", login, phone, year)
					if err != nil {
						context.String(200,"Ошибка")
					} else {
						context.Redirect(302, "/account")
					}
				}
			} else {
				context.String(200, "Не все заполненые поля")
			}
		})
	}

	// Поиск
	search := engine.Group("/search")
	{
		// Поиск аккаунта HTML шаблона
		search.GET("/", func(context *gin.Context) {
			context.HTML(200, "search.html", gin.H{
				"title" : "Регистрация аккаунта",
			})
		})

		// Поиск аккаунта
		search.GET("/data", func(context *gin.Context) {
			data, dataOK := context.GetQuery("data")	// Получает GET запрос

			if dataOK == true {
				db := settings.DataBaseOpen()
				defer db.Close()

				res, err := models.Search(db, models.User{}, &data)
				if err != nil {
					context.String(200, "Нет данных")
				} else {
					context.String(200, res.Login)
				}
			}
		})
	}

}
