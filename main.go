package main

import (
	"runtime"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/0xhexnumbers/partysim/mp1"
)

var CPUS = runtime.NumCPU()

const DEFAULT_AI_TEXT = "[I will tell you what the AI recommends]"

func makeAIUI(w fyne.Window, boardWdgt *boardWidget) fyne.CanvasObject {
	g := mp1.InitializeGame(boardWdgt.board, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].Char = "Mario"
	g.Players[1].Char = "Luigi"
	g.Players[2].Char = "Peach"
	g.Players[3].Char = "Yoshi"

	gHandler := &GameHandler{Game: g, Controller: &SpaceController{
		Mode:  scmNORMAL,
		Board: boardWdgt,
	}}
	gHandler.Controller.SetPlayerCircPositions(
		g.Players[0].CurrentSpace,
		g.Players[1].CurrentSpace,
		g.Players[2].CurrentSpace,
		g.Players[3].CurrentSpace,
	)
	boardWdgt.spaceMap.setGameHandler(gHandler)

	//Event Response Selector
	//The underlying widget may change types, so we use an container
	//that never changes to hold the changing widget.
	baseResponseContainer := container.New(
		layout.NewHBoxLayout(),
		createUserInputUI(gHandler, boardWdgt.spaceMap),
	)

	//Mode Controller
	modeSelector := widget.NewRadioGroup(
		[]string{"Normal", "Show Player Positions", "Show All Spaces", "Show AI Choice"},
		func(s string) {
			switch s {
			case "Normal":
				gHandler.Controller.SetMode(scmNORMAL)
			case "Show Player Positions":
				gHandler.Controller.SetMode(scmSHOW_PLAYERS)
			case "Show All Spaces":
				gHandler.Controller.SetMode(scmSHOW_ALL_SPACES)
			case "Show AI Choice":
				gHandler.Controller.SetMode(scmAI_SPACE_CHOICE)
			}
		},
	)
	modeSelector.SetSelected("Normal")

	aiThreadCountText := widget.NewLabel("Threads: 1/" + strconv.Itoa(CPUS))
	aiThreadCount := widget.NewSlider(1, float64(CPUS))
	aiThreadCount.Step = 1
	aiThreadCount.OnChanged = func(f float64) {
		aiThreadCountText.SetText(
			"Threads: " + strconv.Itoa(int(f)) +
				"/" + strconv.Itoa(CPUS),
		)
	}

	aiMillisecondsText := widget.NewLabel("1 Milliseconds")
	aiMilliseconds := widget.NewSlider(1, 10000)
	aiMilliseconds.Step = 1
	aiMilliseconds.OnChanged = func(f float64) {
		aiMillisecondsText.SetText(
			strconv.Itoa(int(f)) + " Milliseconds",
		)
	}

	selectionButton := widget.NewButton("Continue with Selection", nil)

	//AI Controller
	aiSelection := widget.NewLabel(DEFAULT_AI_TEXT)
	aiButton := widget.NewButton("Run AI", nil)
	aiFunc := func() {
		go func() {
			selectionButton.Disable()
			aiButton.Disable()

			res := bestMove(*gHandler.Game, int(aiThreadCount.Value), int(aiMilliseconds.Value))
			evtType := gHandler.Game.NextEvent.Type()
			if evtType == mp1.CHAINSPACE_EVT_TYPE {
				cs := res.(mp1.ChainSpace)
				gHandler.Controller.SetAISpaceChoice(cs)
				modeSelector.SetSelected("Show AI Choice")
				aiSelection.SetText("AI Response: [Shown on Board]")
			} else {
				var aiText string
				switch evtType {
				case mp1.ENUM_EVT_TYPE:
					aiText = enumToString(res)
				case mp1.RANGE_EVT_TYPE:
					aiText = rangeToString(res)
				case mp1.COIN_EVT_TYPE:
					aiText = coinToString(res)
				case mp1.PLAYER_EVT_TYPE:
					aiText = playerToString(res, gHandler.Game)
				case mp1.MULTIWIN_PLAYER_EVT_TYPE:
					aiText = multiPlayerToString(res, gHandler.Game)
				}
				aiSelection.SetText("AI Response: " + aiText)
			}

			aiButton.Enable()
			selectionButton.Enable()
		}()
	}
	aiButton.OnTapped = aiFunc

	//Player statistics
	p0 := NewPlayer(g, 0, boardWdgt.spaceMap)
	p1 := NewPlayer(g, 1, boardWdgt.spaceMap)
	p2 := NewPlayer(g, 2, boardWdgt.spaceMap)
	p3 := NewPlayer(g, 3, boardWdgt.spaceMap)

	//Updates player stats and event text after each event execution.
	eventText := widget.NewLabel("")
	preEventHandler := func() {
		p0.Refresh()
		p1.Refresh()
		p2.Refresh()
		p3.Refresh()
		eventText.SetText(g.NextEvent.Question(g))
		aiSelection.SetText(DEFAULT_AI_TEXT)

		//Disable AI if CPU Player handles event
		if gHandler.NextEvent.ControllingPlayer() == mp1.CPU_PLAYER {
			aiButton.Disable()
		} else {
			aiButton.Enable()
		}
	}
	preEventHandler()

	//User Simulation Controller
	selectionFunc := func() {
		if g.NextEvent == nil {
			return
		}
		//Execute Event
		err := gHandler.HandleEvent()
		if err != nil { //Tell user no response is selected
			return
		}

		//Update UI with data from new Event
		preEventHandler()
		baseResponseContainer.Objects[0] = createUserInputUI(gHandler, boardWdgt.spaceMap)
		baseResponseContainer.Refresh()

		//Reset UI
		gHandler.Controller.SetMode(gHandler.Controller.Mode)
	}
	selectionButton.OnTapped = selectionFunc

	//Plug user input widgets together
	userInputRegion := container.New(
		layout.NewHBoxLayout(),
		boardWdgt,
		layout.NewSpacer(),
		container.New(
			//Mode Selector
			layout.NewVBoxLayout(),
			widget.NewLabel("Space Options:"),
			modeSelector,

			//AI Controller
			layout.NewSpacer(),
			aiThreadCountText,
			aiThreadCount,
			aiMillisecondsText,
			aiMilliseconds,
			aiSelection,
			aiButton,

			//User Input to Simulation
			layout.NewSpacer(),
			eventText,
			baseResponseContainer,
			selectionButton,
		),
	)

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
