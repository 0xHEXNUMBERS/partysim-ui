package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/0xhexnumbers/partysim/mp1"
	"github.com/0xhexnumbers/partysim/mp1/board"
)

const (
	aspectRatio = 1.315789474
	height      = 600
)

var (
	imageSize = fyne.NewSize(height*aspectRatio, height)
)

type boardConfig struct {
	board    mp1.Board
	filePath string
	spaceMap SpaceCirc
}

var YTI = boardConfig{
	board:    board.YTI,
	filePath: "./img/YoshisTropicalIsland.png",
	spaceMap: ytiSpaceToPos,
}

var DKJA = boardConfig{
	board:    board.DKJA,
	filePath: "./img/DKsJungleAdventure.png",
	spaceMap: dkjaSpaceToPos,
}

var PBC = boardConfig{
	board:    board.PBC,
	filePath: "./img/PeachsBirthdayCake.png",
	spaceMap: pbcSpaceToPos,
}

var WBC = boardConfig{
	board:    board.WBC,
	filePath: "./img/WariosBattleCanyon.png",
	spaceMap: wbcSpaceToPos,
}

var LER = boardConfig{
	board:    board.LER,
	filePath: "./img/LuigisEngineRoom.png",
	spaceMap: lerSpaceToPos,
}

var MRC = boardConfig{
	board:    board.MRC,
	filePath: "./img/MariosRainbowCastle.png",
	spaceMap: mrcSpaceToPos,
}

var BMM = boardConfig{
	board:    board.BMM,
	filePath: "./img/BowsersMagmaMountain.png",
	spaceMap: bmmSpaceToPos,
}

var ES = boardConfig{
	board:    board.ES,
	filePath: "./img/EternalStar.png",
	spaceMap: esSpaceToPos,
}

type boardWidget struct {
	widget.BaseWidget
	image         *canvas.Image
	spaceMap      SpaceCirc
	circles       *fyne.Container
	playerCircles [4]*space
	aiCircle      *space
	board         mp1.Board
}

// MouseIn is a hook that is called if the mouse pointer enters the element.
func (i *boardWidget) MouseIn(_ *desktop.MouseEvent) {
}

// MouseMoved is a hook that is called if the mouse pointer moved over the element.
func (i *boardWidget) MouseMoved(m *desktop.MouseEvent) {
}

// MouseOut is a hook that is called if the mouse pointer leaves the element.
func (i *boardWidget) MouseOut() {
}

func (i *boardWidget) Tapped(p *fyne.PointEvent) {
}

func (i *boardWidget) CreateRenderer() fyne.WidgetRenderer {
	return boardRenderer{i}
}

type boardRenderer struct {
	img *boardWidget
}

// Destroy is for internal use.
func (i boardRenderer) Destroy() {
}

// Layout is a hook that is called if the widget needs to be laid out.
// This should never call Refresh.
func (i boardRenderer) Layout(s fyne.Size) {
	i.img.image.Resize(imageSize)
}

// MinSize returns the minimum size of the widget that is rendered by this renderer.
func (i boardRenderer) MinSize() fyne.Size {
	return imageSize
}

// Objects returns all objects that should be drawn.
func (i boardRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		i.img.image,
		i.img.circles,
	}
}

// Refresh is a hook that is called if the widget has updated and needs to be redrawn.
// This might trigger a Layout.
func (i boardRenderer) Refresh() {
	i.img.circles.Refresh()
}

func loadImage(filePath string) (img *canvas.Image, err error) {
	var r fyne.Resource
	r, err = fyne.LoadResourceFromPath(filePath)
	if err != nil {
		return
	}
	img = canvas.NewImageFromResource(r)
	img.FillMode = canvas.ImageFillContain
	return
}

func initImage(conf boardConfig,
	p1Conf, p2Conf, p3Conf, p4Conf PlayerConfig) *boardWidget {
	var err error
	boardWdgt := &boardWidget{}
	boardWdgt.board = conf.board
	boardWdgt.image, err = loadImage(conf.filePath)
	if err != nil {
		panic(fmt.Errorf("Cannot load image: %w", err.Error()))
	}
	boardWdgt.spaceMap = conf.spaceMap
	boardWdgt.ExtendBaseWidget(boardWdgt)

	circs := container.New(&circlesLayout{})
	if conf.spaceMap != nil {
		for _, circle := range boardWdgt.spaceMap {
			circs.Add(circle)
		}
	}

	boardWdgt.playerCircles[0] = newSpace(p1Conf.Color, 0, 0, 0, 0)
	boardWdgt.playerCircles[1] = newSpace(p2Conf.Color, 0, 0, 0, 0)
	boardWdgt.playerCircles[2] = newSpace(p3Conf.Color, 0, 0, 0, 0)
	boardWdgt.playerCircles[3] = newSpace(p4Conf.Color, 0, 0, 0, 0)
	circs.Add(boardWdgt.playerCircles[0])
	circs.Add(boardWdgt.playerCircles[1])
	circs.Add(boardWdgt.playerCircles[2])
	circs.Add(boardWdgt.playerCircles[3])

	aiColors := spaceColor{
		dorment:   color.NRGBA{0xc5, 0x88, 0x3a, 0xff},
		highlight: color.NRGBA{0xe3, 0x83, 0x1c, 0xff},
	}

	boardWdgt.aiCircle = newSpace(aiColors, 0, 0, 0, 0)
	circs.Add(boardWdgt.aiCircle)
	boardWdgt.circles = circs

	return boardWdgt
}
