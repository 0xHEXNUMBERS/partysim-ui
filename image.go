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

var (
	imageSize        = fyne.NewSize(1200, 912)
	ytiImage  *image = &image{board: board.YTI}
	dkjaImage *image = &image{board: board.DKJA}
	pbcImage  *image = &image{board: board.PBC}
	wbcImage  *image = &image{board: board.WBC}
	lerImage  *image = &image{board: board.LER}
	mrcImage  *image = &image{board: board.MRC}
	bmmImage  *image = &image{board: board.BMM}
	esImage   *image = &image{board: board.ES}
)

type image struct {
	widget.BaseWidget
	image    *canvas.Image
	spaceMap SpaceCirc
	circles  *fyne.Container
	curSpace mp1.ChainSpace
	board    mp1.Board
}

// MouseIn is a hook that is called if the mouse pointer enters the element.
func (i *image) MouseIn(_ *desktop.MouseEvent) {
}

// MouseMoved is a hook that is called if the mouse pointer moved over the element.
func (i *image) MouseMoved(m *desktop.MouseEvent) {
}

// MouseOut is a hook that is called if the mouse pointer leaves the element.
func (i *image) MouseOut() {
}

func (i *image) Tapped(p *fyne.PointEvent) {
}

func (i *image) CreateRenderer() fyne.WidgetRenderer {
	return imageRenderer{i}
}

type imageRenderer struct {
	img *image
}

// Destroy is for internal use.
func (i imageRenderer) Destroy() {
}

// Layout is a hook that is called if the widget needs to be laid out.
// This should never call Refresh.
func (i imageRenderer) Layout(s fyne.Size) {
	i.img.image.Resize(imageSize)
}

// MinSize returns the minimum size of the widget that is rendered by this renderer.
func (i imageRenderer) MinSize() fyne.Size {
	return imageSize
}

// Objects returns all objects that should be drawn.
func (i imageRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{i.img.image, i.img.circles}
}

// Refresh is a hook that is called if the widget has updated and needs to be redrawn.
// This might trigger a Layout.
func (i imageRenderer) Refresh() {
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

func (i imageLoader) initImage(img *image, filePath string, spaceMap SpaceCirc) {
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
