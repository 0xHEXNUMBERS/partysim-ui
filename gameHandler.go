package main

import (
	"errors"
	"log"

	"github.com/0xhexnumbers/partysim/mp1"
)

type GameHandler struct {
	*mp1.Game
	Response     mp1.Response
	CurrentSpace *space
}

func (g *GameHandler) SetResponse(r mp1.Response) {
	g.Response = r
	log.Printf("Set Response: %#v", r)
}

func (g *GameHandler) HandleEvent() error {
	if g.Response == nil {
		return errors.New("Please choose a response")
	}

	log.Printf("Executing event with response: %#v", g.Response)
	g.Game.HandleEvent(g.Response)
	g.Response = nil

	if g.CurrentSpace != nil {
		g.CurrentSpace.isSelected = false
		g.CurrentSpace.darken()
		g.CurrentSpace = nil
	}
	return nil
}

func (g *GameHandler) SetSpace(s *space) {
	if g.CurrentSpace != nil {
		g.CurrentSpace.isSelected = false
		g.CurrentSpace.darken()
	}
	g.CurrentSpace = s
	g.CurrentSpace.isSelected = true
	g.CurrentSpace.highlight()

	g.SetResponse(s.chainSpace)
}
