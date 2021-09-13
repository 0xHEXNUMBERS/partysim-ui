package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

const (
	highlightAlpha = 200
	dormentAlpha   = 100
	spaceLen       = 25
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
	circ *canvas.Circle
	col  spaceColor
	pos  fyne.Position
}

func newSpace(col spaceColor, x, y float32) *space {
	space := &space{
		circ: canvas.NewCircle(col.dorment),
		col:  col,
		pos:  fyne.NewPos(x, y),
	}
	space.circ.Resize(spaceSize)
	space.circ.StrokeColor = color.Black
	space.circ.StrokeWidth = 3
	space.ExtendBaseWidget(space)
	return space
}

// MouseIn is a hook that is called if the mouse pointer enters the element.
func (s *space) MouseIn(_ *desktop.MouseEvent) {
	s.circ.FillColor = s.col.highlight
	s.circ.Refresh()
}

// MouseMoved is a hook that is called if the mouse pointer moved over the element.
func (s *space) MouseMoved(_ *desktop.MouseEvent) {
}

// MouseOut is a hook that is called if the mouse pointer leaves the element.
func (s *space) MouseOut() {
	s.circ.FillColor = s.col.dorment
	s.circ.Refresh()
}

func (s *space) Tapped(_ *fyne.PointEvent) {
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
}

type circlesLayout struct{}

// Layout will manipulate the listed CanvasObjects Size and Position
// to fit within the specified size.
func (c *circlesLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, o := range objects {
		o.Resize(spaceSize)
		circ, ok := o.(*space)
		if ok {
			o.Move(circ.pos)
		}
	}
}

// MinSize calculates the smallest size that will fit the listed
// CanvasObjects using this Layout algorithm.
func (c *circlesLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return imageSize
}
