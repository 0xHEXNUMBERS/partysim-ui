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

	//Player statistics
	p0 := NewPlayer(g, 0, boardWdgt.spaceMap)
	p1 := NewPlayer(g, 1, boardWdgt.spaceMap)
	p2 := NewPlayer(g, 2, boardWdgt.spaceMap)
	p3 := NewPlayer(g, 3, boardWdgt.spaceMap)

	//Updates player stats and event text after each event execution.
	eventText := widget.NewLabel("")
	setText := func() {
		p0.Refresh()
		p1.Refresh()
		p2.Refresh()
		p3.Refresh()
		eventText.SetText(g.NextEvent.Question(g))
	}
	setText()

	//Collectables View
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

	//Event Response Selector
	//The underlying widget may change types, so we use an container
	//that never changes to hold the changing widget.
	baseResponseContainer := container.New(
		layout.NewHBoxLayout(),
		createUserInputUI(gHandler, boardWdgt.spaceMap),
	)

	//Mode Controller
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

	//User Simulation Controller
	selectionButton := widget.NewButton("Continue with Selection", func() {
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
	})

	//AI Controller
	aiSelection := widget.NewLabel("[I will tell you what the AI recommends]")
	aiButton := widget.NewButton("Run AI", nil)
	aiFunc := func() {
		go func() {
			selectionButton.Disable()
			aiButton.Disable()

			res := bestMove(*gHandler.Game)
			aiSelection.SetText(fmt.Sprintf("%#v", res))

			aiButton.Enable()
			selectionButton.Enable()
		}()
	}
	aiButton.OnTapped = aiFunc

	//Plug user input widgets together
	userInputRegion := container.New(
		layout.NewHBoxLayout(),
		boardWdgt,
		layout.NewSpacer(),
		container.New(
			//Mode Selector
			layout.NewVBoxLayout(),
			modeSelector,

			//AI Controller
			layout.NewSpacer(),
			aiSelection,
			aiButton,

			//User Input to Simulation
			layout.NewSpacer(),
			eventText,
			baseResponseContainer,
			selectionButton,
		),
	)

	//Plug everything together
	ui := container.New(
		layout.NewVBoxLayout(),
		userInputRegion,
		layout.NewSpacer(),
		collectablesPanel,
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
