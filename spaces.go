package main

import (
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/0xhexnumbers/partysim/mp1"
)

const (
	highlightAlpha = 200
	dormentAlpha   = 100
	spaceLen       = 15
)

type spaceColor struct {
	dorment   color.NRGBA
	highlight color.NRGBA
}

func colorToSpaceColor(c color.NRGBA) string {
	switch c {
	case redSpace.highlight:
		return "redSpace"
	case blueSpace.highlight:
		return "blueSpace"
	case greenSpace.highlight:
		return "greenSpace"
	case graySpace.highlight:
		return "graySpace"
	}
	return ""
}

var (
	spaceSize = fyne.NewSize(spaceLen, spaceLen)
	redSpace  = spaceColor{
		dorment:   color.NRGBA{R: 255, G: 0, B: 0, A: dormentAlpha},
		highlight: color.NRGBA{R: 255, G: 0, B: 0, A: highlightAlpha},
	}
	greenSpace = spaceColor{
		dorment:   color.NRGBA{R: 0, G: 255, B: 0, A: dormentAlpha},
		highlight: color.NRGBA{R: 0, G: 255, B: 0, A: highlightAlpha},
	}
	blueSpace = spaceColor{
		dorment:   color.NRGBA{R: 0, G: 0, B: 255, A: dormentAlpha},
		highlight: color.NRGBA{R: 0, G: 0, B: 255, A: highlightAlpha},
	}
	graySpace = spaceColor{
		dorment:   color.NRGBA{R: 128, G: 128, B: 128, A: dormentAlpha},
		highlight: color.NRGBA{R: 128, G: 128, B: 128, A: highlightAlpha},
	}
)

type space struct {
	widget.BaseWidget
	circ       *canvas.Circle
	col        spaceColor
	pos        fyne.Position
	chainSpace mp1.ChainSpace
	gHandler   *GameHandler
	isSelected bool
}

func newSpace(col spaceColor, x, y float32, c, s int) *space {
	space := &space{
		circ:       canvas.NewCircle(col.dorment),
		col:        col,
		pos:        fyne.NewPos(x, y),
		chainSpace: mp1.NewChainSpace(c, s),
	}
	space.circ.Resize(spaceSize)
	space.circ.StrokeColor = color.Black
	space.circ.StrokeWidth = 3
	space.ExtendBaseWidget(space)
	return space
}

func (s *space) highlight() {
	s.circ.FillColor = s.col.highlight
	s.Refresh()
}

func (s *space) darken() {
	if !s.isSelected { //Only darken if the space is not selected
		s.circ.FillColor = s.col.dorment
		s.Refresh()
	}
}

// MouseIn is a hook that is called if the mouse pointer enters the element.
func (s *space) MouseIn(_ *desktop.MouseEvent) {
	s.highlight()
}

// MouseMoved is a hook that is called if the mouse pointer moved over the element.
func (s *space) MouseMoved(_ *desktop.MouseEvent) {
}

// MouseOut is a hook that is called if the mouse pointer leaves the element.
func (s *space) MouseOut() {
	s.darken()
}

//TODO: create selectSpace(*space/chainspace) func
//will require either data stored in *space or package
//scope.
func (s *space) Tapped(p *fyne.PointEvent) {
	s.gHandler.SetSpace(s)
}

func (s *space) CreateRenderer() fyne.WidgetRenderer {
	return s
}

// Destroy is for internal use.
func (s *space) Destroy() {
}

// Layout is a hook that is called if the widget needs to be laid out.
// This should never call Refresh.
func (s *space) Layout(_ fyne.Size) {
}

// MinSize returns the minimum size of the widget that is rendered by this renderer.
func (s *space) MinSize() fyne.Size {
	return spaceSize
}

// Objects returns all objects that should be drawn.
func (s *space) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{s.circ}
}

// Refresh is a hook that is called if the widget has updated and needs to be redrawn.
// This might trigger a Layout.
func (s *space) Refresh() {
	s.circ.Refresh()
}

type circlesLayout struct{}

// Layout will manipulate the listed CanvasObjects Size and Position
// to fit within the specified size.
func (c *circlesLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, o := range objects {
		o.Resize(spaceSize)
		circ, ok := o.(*space)
		if ok {
			o.Move(fyne.NewPos(circ.pos.X*size.Width, circ.pos.Y*size.Height))
		}
	}
}

// MinSize calculates the smallest size that will fit the listed
// CanvasObjects using this Layout algorithm.
func (c *circlesLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return imageSize
}

func showSpace(spaceMap SpaceCirc, c mp1.ChainSpace) {
	if circ, ok := spaceMap[c]; ok {
		log.Printf("Showing circle at %#v", c)
		circ.Show()
	}
}

func hideSpace(spaceMap SpaceCirc, c mp1.ChainSpace) {
	if circ, ok := spaceMap[c]; ok {
		circ.Hide()
	}
}

func hideAllSpaces(spaceMap SpaceCirc) {
	for _, circ := range spaceMap {
		circ.Hide()
	}
}
