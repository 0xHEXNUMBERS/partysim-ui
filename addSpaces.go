package main

//This source file is not used anywhere in the repository. It's only here
//in case I need to go back and replace the space coordinates in MP maps.

/*import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)*/

/*func addCirclesUI(w fyne.Window, img *image) fyne.CanvasObject {
	curSpaceLabel := widget.NewLabel("mp1.ChainSpace{0, 0}")
	setCircleColor := func(s spaceColor) func() {
		return func() {
			curSpaceLabel.SetText(fmt.Sprintf("%#v", img.curSpace))
			img.mouseCircle = canvas.NewCircle(s.highlight)
			img.circles.Add(img.mouseCircle)
		}
	}
	red := widget.NewButton("Create Red Circle", setCircleColor(redSpace))
	green := widget.NewButton("Create Green Circle", setCircleColor(greenSpace))
	blue := widget.NewButton("Create Blue Circle", setCircleColor(blueSpace))
	gray := widget.NewButton("Create Gray Circle", setCircleColor(graySpace))
	printPos := widget.NewButton("Print Space Info", func() {
		for cs, space := range img.spaceMap {
			fmt.Printf("makeChainSpace(%d, %d): newSpace(%s, %d, %d),\n",
				cs.Chain, cs.Space, colorToSpaceColor(space.col.highlight), int(space.pos.X), int(space.pos.Y))
		}
	})
	buttonPanel := container.New(
		layout.NewHBoxLayout(),
		layout.NewSpacer(),
		red,
		green,
		blue,
		gray,
		printPos,
		curSpaceLabel,
		layout.NewSpacer(),
	)
	image := container.New(
		layout.NewHBoxLayout(),
		layout.NewSpacer(),
		img,
		layout.NewSpacer(),
	)
	ui := container.New(
		layout.NewVBoxLayout(),
		image,
		layout.NewSpacer(),
		buttonPanel,
	)
	return ui
}

func (i *image) MouseMoved(m *desktop.MouseEvent) {
	if i.mouseCircle != nil {
		i.mouseCircle.Move(m.Position)
		i.mouseCircle.Refresh()
	}
}

func (i *image) Tapped(p *fyne.PointEvent) {
	if i.mouseCircle != nil {
		i.mouseCircle.Move(p.Position)
		i.spaceMap[i.curSpace] = &space{
			widget.BaseWidget{},
			nil,
			spaceColor{
				highlight: i.mouseCircle.FillColor.(color.NRGBA),
			},
			i.mouseCircle.Position(),
		}

		i.mouseCircle = nil
		i.curSpace.Space++
		if i.curSpace.Space >= len((*i.board.Chains)[i.curSpace.Chain]) {
			i.curSpace.Chain++
			i.curSpace.Space = 0
		}
		i.mouseCircle.Refresh()
	}
}*/
