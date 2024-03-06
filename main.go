package main

import (
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hello World")

	// go notifyTest(a)
	w.SetContent(homeContent(w))

	w.ShowAndRun()
}
