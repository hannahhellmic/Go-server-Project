package main

import (
	"awesomeProject/accounts"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	accountsHandler := accounts.New()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/account", accountsHandler.GetAccount)
	e.POST("/account/create", accountsHandler.CreateAccount)
	e.DELETE("/account", accountsHandler.DeleteAccount)
	e.PATCH("/account/changebalance", accountsHandler.PathAccount)
	e.PATCH("/account/change", accountsHandler.ChangeAccount)
	e.PATCH("/account/transfer", accountsHandler.TransferAccount)
	e.GET("account/all", accountsHandler.ListAccounts)
	e.GET("account/transactions", accountsHandler.TransactionsList)

	e.Logger.Fatal(e.Start(":1323"))
}
