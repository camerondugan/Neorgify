package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

var pickedFolder = false

func main() {
	a := app.NewWithID("com.camerondugan.neorgify")
	w := a.NewWindow("Neorgify")
	w.Resize(fyne.NewSize(600, 500))

	theme := a.Settings().Theme()

	homeLabel := "Home"
	settingsLabel := "Settings"
	if fyne.CurrentDevice().IsMobile() {
		homeLabel = ""
		settingsLabel = ""
	}

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon(homeLabel, theme.Icon("home"), homeContent(w)),
		container.NewTabItemWithIcon(settingsLabel, theme.Icon("settings"), settingsContent(w)),
	)

	tabs.SetTabLocation(container.TabLocationBottom)

	// go notifyTest(a)
	w.SetContent(tabs)

	w.ShowAndRun()
}
