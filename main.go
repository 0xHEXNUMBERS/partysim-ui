package main

import (
	"runtime"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/0xhexnumbers/partysim/mp1"
)

var CPUS = runtime.NumCPU()

const DEFAULT_AI_TEXT = "[I will tell you what the AI recommends]"

func makeMainUI(boardWdgt *boardWidget, gc mp1.GameConfig,
	p1Conf, p2Conf, p3Conf, p4Conf PlayerConfig) fyne.CanvasObject {
	g := mp1.InitializeGame(boardWdgt.board, gc)
	g.Players[0].Char = p1Conf.Name
	g.Players[1].Char = p2Conf.Name
	g.Players[2].Char = p3Conf.Name
	g.Players[3].Char = p4Conf.Name

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
	p1 := NewPlayer(g, 0, boardWdgt.spaceMap, p1Conf)
	p2 := NewPlayer(g, 1, boardWdgt.spaceMap, p2Conf)
	p3 := NewPlayer(g, 2, boardWdgt.spaceMap, p3Conf)
	p4 := NewPlayer(g, 3, boardWdgt.spaceMap, p4Conf)

	//Updates player stats and event text after each event execution.
	eventText := widget.NewLabel("")
	preEventHandler := func() {
		p1.Refresh()
		p2.Refresh()
		p3.Refresh()
		p4.Refresh()
		eventText.SetText(g.NextEvent.Question(g))

		aiSelection.SetText(DEFAULT_AI_TEXT)
		//Disable AI if CPU Player handles event
		if gHandler.NextEvent.ControllingPlayer() == mp1.CPU_PLAYER {
			aiButton.Disable()
		} else {
			aiButton.Enable()
		}

		baseResponseContainer.Objects[0] = createUserInputUI(gHandler, boardWdgt.spaceMap)
		baseResponseContainer.Refresh()

		//Reset UI
		gHandler.Controller.SetMode(gHandler.Controller.Mode)
	}
	preEventHandler()

	postGameHandler := func() {
		p1.Refresh()
		p2.Refresh()
		p3.Refresh()
		p4.Refresh()

		//Reuse Event Text to show winners
		winnerIDs := g.Winners()
		winnersText := ""
		for _, pID := range winnerIDs {
			winnersText += g.Players[pID].Char + " & "
		}
		winnersText = winnersText[:len(winnersText)-2]
		winnersText += "have won the game!"
		eventText.SetText(winnersText)

		//Update user input UI to say the game has finished.
		baseResponseContainer.Objects[0] = widget.NewLabel("The game has finished!")
		baseResponseContainer.Refresh()

		//Disable all buttons
		aiButton.Disable()
		selectionButton.Disable()
	}

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

		if g.NextEvent != nil {
			//Update UI with data from new Event
			preEventHandler()
		} else {
			postGameHandler()
		}
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
			p1,
			p2,
		),
		container.New(
			layout.NewHBoxLayout(),
			p3,
			p4,
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

func makePlayerSelectEntry() *widget.Select {
	return widget.NewSelect([]string{
		"Mario", "Luigi", "Peach", "Yoshi", "Wario", "Donkey Kong",
	}, nil)
}

func makeConfigScreen(canvas fyne.Canvas) fyne.CanvasObject {
	p1 := makePlayerSelectEntry()
	p2 := makePlayerSelectEntry()
	p3 := makePlayerSelectEntry()
	p4 := makePlayerSelectEntry()

	p1.SetSelected("Mario")
	p2.SetSelected("Luigi")
	p3.SetSelected("Peach")
	p4.SetSelected("Yoshi")

	playerSelectors := container.New(
		layout.NewVBoxLayout(),
		container.New(
			layout.NewHBoxLayout(),
			widget.NewLabel("Player 1:"),
			p1,
		),
		container.New(
			layout.NewHBoxLayout(),
			widget.NewLabel("Player 2:"),
			p2,
		),
		container.New(
			layout.NewHBoxLayout(),
			widget.NewLabel("Player 3:"),
			p3,
		),
		container.New(
			layout.NewHBoxLayout(),
			widget.NewLabel("Player 4:"),
			p4,
		),
	)
	boardSelector := widget.NewRadioGroup(
		[]string{
			"Mario's Rainbow Castle",
			"DK's Jungle Adventure",
			"Peach's Birthday Cake",
			"Yoshi's Tropical Island",
			"Wario's Battle Canyon",
			"Luigi's Engine Room",
			"Bowser's Magma Mountain",
			"Eternal Star",
		},
		nil,
	)
	boardSelector.SetSelected("Mario's Rainbow Castle")

	maxTurnsInput := widget.NewSelect([]string{"20", "35", "50"}, nil)
	maxTurnsInput.SetSelected("20")

	bonusStarsInput := widget.NewCheck("Bonus Stars", nil)
	bonusStarsInput.SetChecked(true)

	koopaInput := widget.NewCheck("Koopa on Board", nil)
	koopaInput.SetChecked(true)

	booInput := widget.NewCheck("Boo on Board", nil)
	booInput.SetChecked(true)

	redDiceInput := widget.NewCheck("Red Dice", nil)
	redDiceInput.SetChecked(false)

	blueDiceInput := widget.NewCheck("Blue Dice", nil)
	blueDiceInput.SetChecked(false)

	warpDiceInput := widget.NewCheck("Warp Dice", nil)
	warpDiceInput.SetChecked(false)

	eventsDiceInput := widget.NewCheck("Events Dice", nil)
	eventsDiceInput.SetChecked(false)

	gameConfigContainer := container.New(
		layout.NewVBoxLayout(),
		container.New(
			layout.NewHBoxLayout(),
			widget.NewLabel("Max Turns:"),
			maxTurnsInput,
		),
		bonusStarsInput,
		koopaInput,
		booInput,
		redDiceInput,
		blueDiceInput,
		warpDiceInput,
		eventsDiceInput,
	)

	startGameButton := widget.NewButton("Start Game", func() {
		var mp1BoardConfig boardConfig
		switch boardSelector.Selected {
		case "Mario's Rainbow Castle":
			mp1BoardConfig = MRC
		case "DK's Jungle Adventure":
			mp1BoardConfig = DKJA
		case "Peach's Birthday Cake":
			mp1BoardConfig = PBC
		case "Yoshi's Tropical Island":
			mp1BoardConfig = YTI
		case "Wario's Battle Canyon":
			mp1BoardConfig = WBC
		case "Luigi's Engine Room":
			mp1BoardConfig = LER
		case "Bowser's Magma Mountain":
			mp1BoardConfig = BMM
		case "Eternal Star":
			mp1BoardConfig = ES
		}

		p1Conf := PlayerNameToConfig(p1.Selected)
		p2Conf := PlayerNameToConfig(p2.Selected)
		p3Conf := PlayerNameToConfig(p3.Selected)
		p4Conf := PlayerNameToConfig(p4.Selected)

		mp1Board := initImage(
			mp1BoardConfig, p1Conf, p2Conf, p3Conf, p4Conf,
		)

		maxTurns, _ := strconv.Atoi(maxTurnsInput.Selected)

		gc := mp1.GameConfig{
			MaxTurns:     uint8(maxTurns),
			NoBonusStars: !bonusStarsInput.Checked,
			NoKoopa:      !koopaInput.Checked,
			NoBoo:        !booInput.Checked,
			RedDice:      redDiceInput.Checked,
			BlueDice:     blueDiceInput.Checked,
			WarpDice:     warpDiceInput.Checked,
			EventsDice:   eventsDiceInput.Checked,
		}

		canvas.SetContent(
			makeMainUI(
				mp1Board, gc, p1Conf, p2Conf, p3Conf, p4Conf,
			),
		)
	})

	ui := container.New(
		layout.NewVBoxLayout(),
		container.New(
			layout.NewHBoxLayout(),
			playerSelectors,
			gameConfigContainer,
		),
		boardSelector,
		startGameButton,
	)

	return ui
}

func main() {
	uiApp := app.NewWithID("me.hexnumbers.partysim")
	uiApp.Settings().SetTheme(theme.DarkTheme())

	window := uiApp.NewWindow("PartySim")
	mainMenu := fyne.NewMainMenu(
		fyne.NewMenu("File"),
		fyne.NewMenu(
			"Themes",
			fyne.NewMenuItem("Dark Theme", func() {
				uiApp.Settings().SetTheme(theme.DarkTheme())
			}),
			fyne.NewMenuItem("Light Theme", func() {
				uiApp.Settings().SetTheme(theme.LightTheme())
			}),
		),
	)
	window.SetMainMenu(mainMenu)

	ui := makeConfigScreen(window.Canvas())
	window.SetContent(ui)
	window.ShowAndRun()
}
