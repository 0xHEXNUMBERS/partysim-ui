package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/0xhexnumbers/partysim/mp1"
)

type PlayerConfig struct {
	Name  string
	Color spaceColor
}

func PlayerNameToConfig(name string) PlayerConfig {
	switch name {
	case "Mario":
		return MarioConfig
	case "Luigi":
		return LuigiConfig
	case "Peach":
		return PeachConfig
	case "Yoshi":
		return YoshiConfig
	case "Wario":
		return WarioConfig
	case "Donkey Kong":
		return DKConfig
	}
	return PlayerConfig{}
}

//PlayerConfig section
//Player Color schemes made using https://www.canva.com/colors/color-wheel
//using Monochromatic color combinations.

var MarioConfig = PlayerConfig{
	Name: "Mario",
	Color: spaceColor{
		highlight: color.NRGBA{0xf1, 0x41, 0x41, 0xff},
		dorment:   color.NRGBA{0xee, 0x11, 0x11, 0xff},
	},
}

var LuigiConfig = PlayerConfig{
	Name: "Luigi",
	Color: spaceColor{
		highlight: color.NRGBA{0x47, 0x83, 0xeb, 0xff},
		dorment:   color.NRGBA{0x19, 0x64, 0xe6, 0xff},
	},
}

var PeachConfig = PlayerConfig{
	Name: "Peach",
	Color: spaceColor{
		highlight: color.NRGBA{0xdb, 0x57, 0xc4, 0xff},
		dorment:   color.NRGBA{0xd2, 0x2d, 0xb5, 0xff},
	},
}

var YoshiConfig = PlayerConfig{
	Name: "Yoshi",
	Color: spaceColor{
		highlight: color.NRGBA{0xf3, 0xef, 0x6d, 0xff},
		dorment:   color.NRGBA{0x13, 0xec, 0x48, 0xff},
	},
}

var WarioConfig = PlayerConfig{
	Name: "Wario",
	Color: spaceColor{
		highlight: color.NRGBA{0xca, 0xe5, 0x4d, 0xff},
		dorment:   color.NRGBA{0xbd, 0xdf, 0x20, 0xff},
	},
}

var DKConfig = PlayerConfig{
	Name: "Donkey Kong",
	Color: spaceColor{
		highlight: color.NRGBA{0xbc, 0x73, 0x14, 0xff},
		dorment:   color.NRGBA{0x8e, 0x57, 0x0f, 0xff},
	},
}

type Player struct {
	widget.BaseWidget
	g        *mp1.Game
	pIdx     int
	spaceMap SpaceCirc
	rect     *canvas.Rectangle
}

func NewPlayer(g *mp1.Game, pIdx int, sm SpaceCirc, pc PlayerConfig) *Player {
	p := &Player{
		BaseWidget: widget.BaseWidget{},
		g:          g,
		pIdx:       pIdx,
		spaceMap:   sm,
	}
	p.rect = canvas.NewRectangle(pc.Color.dorment)
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
