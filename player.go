package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/0xhexnumbers/partysim/mp1"
)

type Player struct {
	widget.BaseWidget
	g        *mp1.Game
	pIdx     int
	spaceMap SpaceCirc
	rect     *canvas.Rectangle
}

func NewPlayer(g *mp1.Game, pIdx int, sm SpaceCirc) *Player {
	p := &Player{
		BaseWidget: widget.BaseWidget{},
		g:          g,
		pIdx:       pIdx,
		spaceMap:   sm,
	}
	switch pIdx {
	case 0:
		p.rect = canvas.NewRectangle(scPlayer1Colors.dorment)
	case 1:
		p.rect = canvas.NewRectangle(scPlayer2Colors.dorment)
	case 2:
		p.rect = canvas.NewRectangle(scPlayer3Colors.dorment)
	case 3:
		p.rect = canvas.NewRectangle(scPlayer4Colors.dorment)
	}
	p.ExtendBaseWidget(p)
	return p
}

func (p *Player) CreateRenderer() fyne.WidgetRenderer {
	return &PlayerRenderer{Player: p}
}

type PlayerRenderer struct {
	Player   *Player
	objects  []fyne.CanvasObject
	charName *widget.Label
	count    *widget.Label
}

// Destroy is for internal use.
func (p *PlayerRenderer) Destroy() {
}

// Layout is a hook that is called if the widget needs to be laid out.
// This should never call Refresh.
func (p *PlayerRenderer) Layout(_ fyne.Size) {
}

// MinSize returns the minimum size of the widget that is rendered by this renderer.
func (p *PlayerRenderer) MinSize() (size fyne.Size) {
	size.Width += p.Player.rect.Size().Width + p.count.Size().Width + p.charName.Size().Width
	size.Height += p.count.Size().Height
	return size
}

// Objects returns all objects that should be drawn.
func (p *PlayerRenderer) Objects() []fyne.CanvasObject {
	return p.objects
}

// Refresh is a hook that is called if the widget has updated and needs to be redrawn.
// This might trigger a Layout.
func (p *PlayerRenderer) Refresh() {
	p.charName = widget.NewLabel(p.Player.g.Players[p.Player.pIdx].Char)
	p.count = widget.NewLabel(
		fmt.Sprintf("Stars: %d\nCoins: %d",
			p.Player.g.Players[p.Player.pIdx].Stars,
			p.Player.g.Players[p.Player.pIdx].Coins,
		),
	)
	p.Player.rect.SetMinSize(fyne.NewSize(10, 10))
	p.objects = []fyne.CanvasObject{
		container.New(
			layout.NewHBoxLayout(),
			p.Player.rect,
			p.charName,
			container.New(
				layout.NewVBoxLayout(),
				p.count,
			),
		),
	}
}
