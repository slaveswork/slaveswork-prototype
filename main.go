package main

import (
	"fmt"
	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	"log"
)

const StaticPath string = "static"

func main() {
	logger := setLogger()
	app := createAstilectronAPP(logger)

	// Handle signals
	//a.HandleSignals()

	err := app.Start()
	errorChecker(err, logger, "main: starting astilectron failed: %w")
	defer app.Close()

	window, err := app.NewWindow(StaticPath+"/index.html", &astilectron.WindowOptions{
		Center: astikit.BoolPtr(true),
		Height: astikit.IntPtr(700),
		Width:  astikit.IntPtr(700),
	})
	errorChecker(err, logger, "main: new window failed: %w")

	err = window.Create()
	errorChecker(err, logger, "main: creating window failed: %w")

	addMenu(app, window)
	OnMessage(window)

	// Blocking pattern
	app.Wait()
}

// Set logger
func setLogger() *log.Logger {
	return log.New(log.Writer(), log.Prefix(), log.Flags())
}

// Create astilectron
func createAstilectronAPP(logger *log.Logger) *astilectron.Astilectron {
	a, err := astilectron.New(logger, astilectron.Options{
		AppName:           "Test",
		BaseDirectoryPath: "",
	})
	errorChecker(err, logger, "main: creating astilectron failed: %w")
	return a
}

// Error Checker
func errorChecker(err error, logger *log.Logger, message string) {
	if err != nil {
		logger.Fatal(fmt.Errorf(message, err))
	}
}

// NewMenu
func addMenu(a *astilectron.Astilectron, w *astilectron.Window) {
	menu := a.NewMenu([]*astilectron.MenuItemOptions{
		{
			Label: astikit.StrPtr("Separator"),
			SubMenu: []*astilectron.MenuItemOptions{
				{Label: astikit.StrPtr("dev tool"), OnClick: func(e astilectron.Event) (deleteListener bool) {
					w.OpenDevTools()
					return
				}},
			},
		},
		{
			Label: astikit.StrPtr("Window"),
			SubMenu: []*astilectron.MenuItemOptions{
				{Label: astikit.StrPtr("Minimize"), Role: astilectron.MenuItemRoleMinimize},
				{Label: astikit.StrPtr("Close"), Role: astilectron.MenuItemRoleClose},
			},
		},
		{
			Label: astikit.StrPtr("Help"),
			SubMenu: []*astilectron.MenuItemOptions{
				{Label: astikit.StrPtr("Minimize"), Role: astilectron.MenuItemRoleMinimize},
				{Label: astikit.StrPtr("Close"), Role: astilectron.MenuItemRoleClose},
			},
		},
	})

	err := menu.Create()
	if err != nil {
		fmt.Println("Menu Error!")
	}
}

// OnMessage
func OnMessage(w *astilectron.Window) {
	w.OnMessage(func(m *astilectron.EventMessage) interface{} {
		// Unmarshal
		var s string
		m.Unmarshal(&s)
		log.Println(s)
		// Process message
		if s == "start" {
			log.Println("server start")
			return "server start!!"
		}
		return nil
	})
}
