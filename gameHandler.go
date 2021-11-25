package main

import (
	"errors"
	"log"

	"github.com/0xhexnumbers/partysim/mp1"
)

type GameHandler struct {
	*mp1.Game
	Response     mp1.Response
	Controller   *SpaceController
	CurrentSpace *space
}

func (g *GameHandler) SetResponse(r mp1.Response) {
	g.Response = r
	log.Printf("GameHandler.SetResponse: Settting Response to %#v", r)
}

func (g *GameHandler) HandleEvent() error {
	if g.Response == nil {
		return errors.New("Please choose a response")
	}

	log.Printf("GameHandler.HandleEvent: Executing event with response %#v", g.Response)
	g.Game.HandleEvent(g.Response)
	g.Response = nil

	if g.CurrentSpace != nil {
		g.CurrentSpace.isSelected = false
		g.CurrentSpace.darken()
		g.CurrentSpace = nil
	}

	g.Controller.Reset(
		g.Players[0].CurrentSpace,
		g.Players[1].CurrentSpace,
		g.Players[2].CurrentSpace,
		g.Players[3].CurrentSpace,
	)

	return nil
}

func (g *GameHandler) SetSpace(s *space) {
	if g.Controller.Mode != scmNORMAL {
		log.Println("GameHandler.SetSpace: Mode is not scmNORMAL")
		return
	}

	log.Printf("GameHandler.SetSpace: Highlighting space %#v", s.chainSpace)

	if g.CurrentSpace != nil {
		g.CurrentSpace.isSelected = false
		g.CurrentSpace.darken()
	}
	g.CurrentSpace = s
	g.CurrentSpace.isSelected = true
	g.CurrentSpace.highlight()

	g.SetResponse(s.chainSpace)
}
