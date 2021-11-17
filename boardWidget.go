package main

import (
	"fmt"

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
	imageSize              = fyne.NewSize(height*aspectRatio, height)
	ytiImage  *boardWidget = &boardWidget{board: board.YTI}
	dkjaImage *boardWidget = &boardWidget{board: board.DKJA}
	pbcImage  *boardWidget = &boardWidget{board: board.PBC}
	wbcImage  *boardWidget = &boardWidget{board: board.WBC}
	lerImage  *boardWidget = &boardWidget{board: board.LER}
	mrcImage  *boardWidget = &boardWidget{board: board.MRC}
	bmmImage  *boardWidget = &boardWidget{board: board.BMM}
	esImage   *boardWidget = &boardWidget{board: board.ES}
)

type boardWidget struct {
	widget.BaseWidget
	image         *canvas.Image
	spaceMap      SpaceCirc
	circles       *fyne.Container
	playerCircles [4]*space
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

type imageLoader struct {
	err error
}

func (i *imageLoader) loadImage(filePath string) (img *canvas.Image) {
	if i.err == nil {
		var r fyne.Resource
		r, i.err = fyne.LoadResourceFromPath(filePath)
		if i.err != nil {
			return
		}
		img = canvas.NewImageFromResource(r)
		img.FillMode = canvas.ImageFillContain
	}
	return
}

func (i imageLoader) initImage(img *boardWidget, filePath string, spaceMap SpaceCirc) {
	if i.err != nil {
		return
	}
	img.image = i.loadImage(filePath)
	img.spaceMap = spaceMap
	img.ExtendBaseWidget(img)
	circs := container.New(&circlesLayout{})
	if spaceMap != nil {
		for _, circle := range img.spaceMap {
			circs.Add(circle)
		}
	}
	img.playerCircles[0] = newSpace(scPlayer1Colors, 0, 0, 0, 0)
	img.playerCircles[1] = newSpace(scPlayer2Colors, 0, 0, 0, 0)
	img.playerCircles[2] = newSpace(scPlayer3Colors, 0, 0, 0, 0)
	img.playerCircles[3] = newSpace(scPlayer4Colors, 0, 0, 0, 0)
	circs.Add(img.playerCircles[0])
	circs.Add(img.playerCircles[1])
	circs.Add(img.playerCircles[2])
	circs.Add(img.playerCircles[3])
	img.circles = circs
}

func load_images() {
	i := imageLoader{}
	i.initImage(ytiImage, "./img/YoshisTropicalIsland.png", ytiSpaceToPos)
	i.initImage(dkjaImage, "./img/DKsJungleAdventure.png", dkjaSpaceToPos)
	i.initImage(pbcImage, "./img/PeachsBirthdayCake.png", pbcSpaceToPos)
	i.initImage(wbcImage, "./img/WariosBattleCanyon.png", wbcSpaceToPos)
	i.initImage(lerImage, "./img/LuigisEngineRoom.png", lerSpaceToPos)
	i.initImage(mrcImage, "./img/MariosRainbowCastle.png", mrcSpaceToPos)
	i.initImage(bmmImage, "./img/BowsersMagmaMountain.png", bmmSpaceToPos)
	i.initImage(esImage, "./img/EternalStar.png", esSpaceToPos)
	if i.err != nil {
		panic(fmt.Errorf("cannot load image: %w", i.err))
	}
}
