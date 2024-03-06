package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hello World")
	w.Resize(fyne.NewSize(600, 500))

	// go notifyTest(a)
	w.SetContent(homeContent(w))

	w.ShowAndRun()
}
