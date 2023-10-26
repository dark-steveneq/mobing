package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/dark-steveneq/mobiapi"
	"github.com/dark-steveneq/mobing/internal"
)

func main() {
	a := app.New()
	api, err := mobiapi.New("", "mobiNG - https://github.com/dark-steveneq")
	if err != nil {
		log.Println("Trouble creating API instance, reason:", err)
		a.SendNotification(fyne.NewNotification("Couldn't start mobiNG!", "Trouble creating API instance, reason: '"+err.Error()+"'."))
		return
	}

	internal.Loginui(a, api)

	a.Run()
}
