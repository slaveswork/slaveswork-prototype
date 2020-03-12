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
				{Label: astikit.StrPtr("Normal 1")},
				{Label: astikit.StrPtr("Normal 2"),
					OnClick: func(e astilectron.Event) (deleteListener bool) {
						w.SendMessage("hello", func(m *astilectron.EventMessage) {
							// Unmarshal
							var s string
							m.Unmarshal(&s)

							// Process message
							log.Printf("received %s\n", s)
						})
						return
					},
				},
				{Type: astilectron.MenuItemTypeSeparator},
				{Label: astikit.StrPtr("dev tool"), OnClick: func(e astilectron.Event) (deleteListener bool) {
					w.OpenDevTools()
					return
				}},
			},
		},
		{
			Label: astikit.StrPtr("Checkbox"),
			SubMenu: []*astilectron.MenuItemOptions{
				{Checked: astikit.BoolPtr(true), Label: astikit.StrPtr("Checkbox 1"), Type: astilectron.MenuItemTypeCheckbox},
				{Label: astikit.StrPtr("Checkbox 2"), Type: astilectron.MenuItemTypeCheckbox},
				{Label: astikit.StrPtr("Checkbox 3"), Type: astilectron.MenuItemTypeCheckbox},
			},
		},
		{
			Label: astikit.StrPtr("Radio"),
			SubMenu: []*astilectron.MenuItemOptions{
				{Checked: astikit.BoolPtr(true), Label: astikit.StrPtr("Radio 1"), Type: astilectron.MenuItemTypeRadio},
				{Label: astikit.StrPtr("Radio 2"), Type: astilectron.MenuItemTypeRadio},
				{Label: astikit.StrPtr("Radio 3"), Type: astilectron.MenuItemTypeRadio},
			},
		},
		{
			Label: astikit.StrPtr("Roles"),
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
