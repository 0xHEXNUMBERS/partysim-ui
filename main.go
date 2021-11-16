package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/0xhexnumbers/partysim/mp1"
)

func makeAIUI(w fyne.Window, boardWdgt *boardWidget) fyne.CanvasObject {
	g := mp1.InitializeGame(boardWdgt.board, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].Char = "Mario"
	g.Players[1].Char = "Luigi"
	g.Players[2].Char = "Peach"
	g.Players[3].Char = "Yoshi"

	gHandler := &GameHandler{Game: g, Controller: &SpaceController{
		Mode: scmNORMAL,
		PlayerPos: [4]mp1.ChainSpace{
			g.Players[0].CurrentSpace,
			g.Players[1].CurrentSpace,
			g.Players[2].CurrentSpace,
			g.Players[3].CurrentSpace,
		},
		Board: boardWdgt,
	}}
	gHandler.Controller.SetPlayerCircPositions()
	boardWdgt.spaceMap.setGameHandler(gHandler)

	p0 := NewPlayer(g, 0, boardWdgt.spaceMap)
	p1 := NewPlayer(g, 1, boardWdgt.spaceMap)
	p2 := NewPlayer(g, 2, boardWdgt.spaceMap)
	p3 := NewPlayer(g, 3, boardWdgt.spaceMap)
	eventText := widget.NewLabel("")
	setText := func() {
		p0.Refresh()
		p1.Refresh()
		p2.Refresh()
		p3.Refresh()
		eventText.SetText(g.NextEvent.Question(g))
	}
	setText()

	collectablesPanel := container.New(
		layout.NewVBoxLayout(),
		container.New(
			layout.NewHBoxLayout(),
			p0,
			p1,
		),
		container.New(
			layout.NewHBoxLayout(),
			p2,
			p3,
		),
	)
	userPanel := container.New(
		layout.NewVBoxLayout(),
		eventText,
	)
	baseResponseContainer := container.New(layout.NewHBoxLayout(), createUserInputUI(gHandler, boardWdgt.spaceMap))
	modeSelector := widget.NewRadioGroup(
		[]string{"Normal", "Show Player Positions", "Show All Spaces"},
		func(s string) {
			switch s {
			case "Normal":
				gHandler.Controller.SetMode(scmNORMAL)
			case "Show Player Positions":
				gHandler.Controller.SetMode(scmSHOW_PLAYERS)
			case "Show All Spaces":
				gHandler.Controller.SetMode(scmSHOW_ALL_SPACES)
			}
		},
	)
	modeSelector.SetSelected("Normal")
	aiSelection := widget.NewLabel("[I will tell you what the AI recommends]")
	aiPanel := container.New(
		layout.NewHBoxLayout(),
		baseResponseContainer,
		modeSelector,
		container.New(
			layout.NewVBoxLayout(),
			aiSelection,
			widget.NewButton("Continue with Selection", func() {
				if g.NextEvent == nil {
					return
				}
				//Execute Event
				err := gHandler.HandleEvent()
				if err != nil { //Tell user no response is selected
					return
				}

				//Update UI with data from new Event
				setText()
				baseResponseContainer.Objects[0] = createUserInputUI(gHandler, boardWdgt.spaceMap)
				baseResponseContainer.Refresh()
			}),
			widget.NewButton("Run AI", func() {
				res := bestMove(*gHandler.Game)
				aiSelection.SetText(fmt.Sprintf("%#v", res))
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
		boardWdgt,
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
