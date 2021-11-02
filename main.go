package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/0xhexnumbers/partysim/mp1"
)

func makeAIUI(w fyne.Window, img *image) fyne.CanvasObject {
	hideAllSpaces(img.spaceMap)

	g := mp1.InitializeGame(img.board, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].Char = "Mario"
	g.Players[1].Char = "Luigi"
	g.Players[2].Char = "Peach"
	g.Players[3].Char = "Yoshi"

	gHandler := &GameHandler{Game: g}
	img.spaceMap.setGameHandler(gHandler)

	p0 := widget.NewLabel("")
	p1 := widget.NewLabel("")
	p2 := widget.NewLabel("")
	p3 := widget.NewLabel("")
	eventText := widget.NewLabel("")
	setText := func() {
		p0.SetText(fmt.Sprintf("Stars: %d\nCoins: %d", g.Players[0].Stars, g.Players[0].Coins))
		p1.SetText(fmt.Sprintf("Stars: %d\nCoins: %d", g.Players[1].Stars, g.Players[1].Coins))
		p2.SetText(fmt.Sprintf("Stars: %d\nCoins: %d", g.Players[2].Stars, g.Players[2].Coins))
		p3.SetText(fmt.Sprintf("Stars: %d\nCoins: %d", g.Players[3].Stars, g.Players[3].Coins))
		eventText.SetText(g.NextEvent.Question(g))
	}
	setText()

	collectablesPanel := container.New(
		layout.NewVBoxLayout(),
		container.New(
			layout.NewHBoxLayout(),
			canvas.NewText(g.Players[0].Char, color.White),
			container.New(
				layout.NewVBoxLayout(),
				p0,
			),
			canvas.NewText(g.Players[1].Char, color.White),
			container.New(
				layout.NewVBoxLayout(),
				p1,
			),
		),
		container.New(
			layout.NewHBoxLayout(),
			canvas.NewText(g.Players[2].Char, color.White),
			container.New(
				layout.NewVBoxLayout(),
				p2,
			),
			canvas.NewText(g.Players[3].Char, color.White),
			container.New(
				layout.NewVBoxLayout(),
				p3,
			),
		),
	)
	userPanel := container.New(
		layout.NewVBoxLayout(),
		eventText,
	)
	baseResponseContainer := container.New(layout.NewHBoxLayout(), createUserInputUI(gHandler, img.spaceMap))
	aiPanel := container.New(
		layout.NewHBoxLayout(),
		baseResponseContainer,
		container.New(
			layout.NewVBoxLayout(),
			canvas.NewText("AI says 1, 2", color.White),
			widget.NewButton("Continue with Selection", func() {
				if g.NextEvent == nil {
					return
				}
				//Execute Event
				err := gHandler.HandleEvent()
				if err != nil { //Tell user no response is selected
					return
				}

				//Reset Space
				hideAllSpaces(img.spaceMap)

				//Update UI with data from new Event
				setText()
				baseResponseContainer.Objects[0] = createUserInputUI(gHandler, img.spaceMap)
				baseResponseContainer.Refresh()
			}),
		),
	)
	panel := container.New(
		layout.NewHBoxLayout(),
		collectablesPanel,
		layout.NewSpacer(),
		userPanel,
		layout.NewSpacer(),
		aiPanel,
	)
	imgLayout := container.New(
		layout.NewHBoxLayout(),
		layout.NewSpacer(),
		img,
		layout.NewSpacer(),
	)
	ui := container.New(
		layout.NewVBoxLayout(),
		imgLayout,
		layout.NewSpacer(),
		panel,
	)
	return ui
}

func main() {
	load_images()

	uiApp := app.New()
	window := uiApp.NewWindow("PartySim")

	ui := makeAIUI(window, lerImage)
	window.SetContent(ui)
	window.ShowAndRun()
}
