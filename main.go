package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func makeAIUI(w fyne.Window, img *image) fyne.CanvasObject {
	collectablesPanel := container.New(
		layout.NewVBoxLayout(),
		container.New(
			layout.NewHBoxLayout(),
			canvas.NewText("Mario", color.White),
			container.New(
				layout.NewVBoxLayout(),
				canvas.NewText("Stars: 1", color.White),
				canvas.NewText("Coins: 25", color.White),
			),
			canvas.NewText("Luigi", color.White),
			container.New(
				layout.NewVBoxLayout(),
				canvas.NewText("Stars: 0", color.White),
				canvas.NewText("Coins: 42", color.White),
			),
		),
		container.New(
			layout.NewHBoxLayout(),
			canvas.NewText("Peach", color.White),
			container.New(
				layout.NewVBoxLayout(),
				canvas.NewText("Stars: 2", color.White),
				canvas.NewText("Coins: 0", color.White),
			),
			canvas.NewText("Yoshi", color.White),
			container.New(
				layout.NewVBoxLayout(),
				canvas.NewText("Stars: 1", color.White),
				canvas.NewText("Coins: 5", color.White),
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
		widget.NewButton("Run AI", func() { w.SetFullScreen(!w.FullScreen()) }),
	)
	panel := container.New(
		layout.NewHBoxLayout(),
		collectablesPanel,
		layout.NewSpacer(),
		userPanel,
		layout.NewSpacer(),
		aiPanel,
	)
	ui := container.New(
		layout.NewVBoxLayout(),
		img,
		layout.NewSpacer(),
		panel,
	)
	return ui
}

func main() {
	load_images()

	uiApp := app.New()
	window := uiApp.NewWindow("PartySim")

	ui := makeAIUI(window, bmmImage)
	window.SetContent(ui)
	window.ShowAndRun()

	/*uiApp := app.New()
	w := uiApp.NewWindow("Test")
	circ := lerSpaceToPos[makeChainSpace(0, 0)]  //newSpace(color.White, 0, 0, false)
	circ2 := lerSpaceToPos[makeChainSpace(0, 1)] //newSpace(color.White, 100, 0, false)
	circles := container.New(
		&circlesLayout{},
		circ,
		circ2,
	)
	l := container.NewMax(lerImage.Image, circles)
	w.SetContent(l)
	w.ShowAndRun()*/

}
