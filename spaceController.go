package main

import (
	"image/color"
	"log"

	"github.com/0xhexnumbers/partysim/mp1"
)

type SpaceControllerMode int

const (
	scmNORMAL SpaceControllerMode = iota
	scmSHOW_PLAYERS
	scmSHOW_ALL_SPACES
)

func (s SpaceControllerMode) String() string {
	switch s {
	case scmNORMAL:
		return "scmNORMAL"
	case scmSHOW_PLAYERS:
		return "scmSHOW_PLAYERS"
	case scmSHOW_ALL_SPACES:
		return "scmSHOW_ALL_SPACES"
	}
	return ""
}

var (
	scPlayer1Colors = spaceColor{
		dorment:   color.NRGBA{0xc5, 0x88, 0x3a, 0xff},
		highlight: color.NRGBA{0xe3, 0x83, 0x1c, 0xff},
	}
	scPlayer2Colors = spaceColor{
		dorment:   color.NRGBA{0x3a, 0xc5, 0x43, 0xff},
		highlight: color.NRGBA{0x1c, 0xe3, 0x20, 0xff},
	}
	scPlayer3Colors = spaceColor{
		dorment:   color.NRGBA{0x3a, 0x77, 0xc5, 0xff},
		highlight: color.NRGBA{0x1c, 0x7c, 0xe3, 0xff},
	}
	scPlayer4Colors = spaceColor{
		dorment:   color.NRGBA{0xc5, 0x3a, 0xbc, 0xff},
		highlight: color.NRGBA{0xe3, 0x1c, 0xdf, 0xff},
	}
)

type SpaceController struct {
	Mode        SpaceControllerMode
	PlayerPos   [4]mp1.ChainSpace
	Board       *boardWidget
	NormalCircs []*space
}

func (s *SpaceController) SetPlayerCircPositions() {
	log.Println("SpaceController.SetPlayerCircPositions: Setting space pos")
	s.Board.playerCircles[0].pos = s.Board.spaceMap[s.PlayerPos[0]].pos
	s.Board.playerCircles[1].pos = s.Board.spaceMap[s.PlayerPos[1]].pos
	s.Board.playerCircles[2].pos = s.Board.spaceMap[s.PlayerPos[2]].pos
	s.Board.playerCircles[3].pos = s.Board.spaceMap[s.PlayerPos[3]].pos
}

func (s *SpaceController) SetNormalCircs(chainSpaces []mp1.ChainSpace) {
	log.Println("SpaceController.SetNormalCircs: Setting normal circs")
	if chainSpaces == nil || len(chainSpaces) == 0 {
		log.Println("SpaceController.SetNormalCircs: No circles to show")
		s.NormalCircs = nil
		return
	}

	s.NormalCircs = make([]*space, len(chainSpaces))
	for i, cs := range chainSpaces {
		log.Printf(
			"SpaceController.SetNormalCircs: Adding chainspace %#v",
			cs,
		)
		s.NormalCircs[i] = s.Board.spaceMap[cs]
	}
}

func (s *SpaceController) HideAllBoardSpaces() {
	log.Printf("SpaceController.HideAllBoardSpaces: Hiding all circles")
	for _, circ := range s.Board.spaceMap {
		circ.Hide()
	}
}

func (s *SpaceController) SetMode(sm SpaceControllerMode) {
	log.Printf("spaceController.SetMode: Setting mode from %s to %s",
		s.Mode, sm)

	s.HideAllBoardSpaces()
	switch sm {
	case scmNORMAL:
		s.Board.playerCircles[0].Hide()
		s.Board.playerCircles[1].Hide()
		s.Board.playerCircles[2].Hide()
		s.Board.playerCircles[3].Hide()
		for _, space := range s.NormalCircs {
			space.Show()
		}
	case scmSHOW_PLAYERS:
		s.Board.playerCircles[0].Show()
		s.Board.playerCircles[1].Show()
		s.Board.playerCircles[2].Show()
		s.Board.playerCircles[3].Show()
	case scmSHOW_ALL_SPACES:
		s.Board.playerCircles[0].Hide()
		s.Board.playerCircles[1].Hide()
		s.Board.playerCircles[2].Hide()
		s.Board.playerCircles[3].Hide()
		for _, space := range s.Board.spaceMap {
			space.Show()
		}
	}
	s.Mode = sm
}
