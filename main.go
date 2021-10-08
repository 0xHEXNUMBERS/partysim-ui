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
	"github.com/0xhexnumbers/partysim/mp1/board"
)

func makeAIUI(w fyne.Window, img *image) fyne.CanvasObject {
	g := mp1.InitializeGame(board.LER, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].Char = "Mario"
	g.Players[1].Char = "Luigi"
	g.Players[2].Char = "Peach"
	g.Players[3].Char = "Yoshi"

	p0 := widget.NewLabel("")
	p1 := widget.NewLabel("")
	p2 := widget.NewLabel("")
	p3 := widget.NewLabel("")
	setText := func() {
		p0.SetText(fmt.Sprintf("Stars: %d\nCoins: %d", g.Players[0].Stars, g.Players[0].Coins))
		p1.SetText(fmt.Sprintf("Stars: %d\nCoins: %d", g.Players[1].Stars, g.Players[1].Coins))
		p2.SetText(fmt.Sprintf("Stars: %d\nCoins: %d", g.Players[2].Stars, g.Players[2].Coins))
		p3.SetText(fmt.Sprintf("Stars: %d\nCoins: %d", g.Players[3].Stars, g.Players[3].Coins))
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
		canvas.NewText("Branch: Pick a space to go to?", color.White),
		canvas.NewText("1, 2 | 2, 1", color.White),
	)
	aiPanel := container.New(
		layout.NewVBoxLayout(),
		canvas.NewText("AI says 1, 2", color.White),
		widget.NewButton("Run AI", func() {
			fmt.Println(g.Turn)
			if g.NextEvent == nil {
				return
			}
			g.HandleEvent(g.Responses()[0])
			setText()
		}),
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
