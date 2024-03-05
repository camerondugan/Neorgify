package main

import (
	"image/color"
	"log"
	"os"
	"time"

	"gioui.org/app"
	_ "gioui.org/app/permission/storage"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
	"gioui.org/x/notify"
)

func main() {
	go func() {
		w := app.NewWindow()
		err := run(w)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window) error {
	th := material.NewTheme()
	notifier, err := notify.NewNotifier()
	if err != nil {
		panic(err)
	}

	go func() {}
	var notification notify.Notification
	notificationSent := false
	if !notificationSent {
		notification, err = notifier.CreateNotification("App Started", "something")
		if err != nil {
			panic(err)
		}
		notificationSent = true
	} else {
		time.Sleep(time.Second)
		notification.Cancel()
		notificationSent = false
	}

	var ops op.Ops
	for {
		switch e := w.NextEvent().(type) {
		case app.DestroyEvent:
			// Doesn't work
			_, err := notifier.CreateNotification("App Closed", "")
			if err != nil {
				panic(err)
			}
			return e.Err
		case app.FrameEvent:
			// This graphics context is used for managing the rendering state.
			gtx := app.NewContext(&ops, e)

			// Define an large label with an appropriate text:
			title := material.H3(th, "Cameron Is Amazing")

			// Change the color of the label.
			title.Color = color.NRGBA{R: 255, G: 0, B: 255, A: 255}

			// Change the position of the label.
			title.Alignment = text.Start

			// Draw the label to the graphics context.
			title.Layout(gtx)

			// Pass the drawing operations to the GPU.
			e.Frame(gtx.Ops)
		}
	}
}
