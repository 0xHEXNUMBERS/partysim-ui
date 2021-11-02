package main

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/0xhexnumbers/partysim/mp1"
)

func createUserInputUI(g *GameHandler, spaceMap SpaceCirc) fyne.CanvasObject {
	var strs []string

	responses := g.NextEvent.Responses()
	switch g.NextEvent.Type() {
	case mp1.ENUM_EVT_TYPE:
		strs = make([]string, len(responses))
		for i, r := range responses {
			if stringer, ok := r.(fmt.Stringer); ok {
				strs[i] = stringer.String()
			} else {
				strs[i] = fmt.Sprintf("%#v", r)
			}
		}
		selection := widget.NewSelect(strs, func(s string) {
			for i, str := range strs {
				if s == str {
					g.SetResponse(responses[i])
					break
				}
			}
		})
		return selection
	case mp1.RANGE_EVT_TYPE:
		strs = make([]string, len(responses))
		for i, r := range responses {
			strs[i] = fmt.Sprintf("%d", r)
		}
		selection := widget.NewSelect(strs, func(s string) {
			i, _ := strconv.Atoi(s)
			g.SetResponse(i)
		})
		return selection
	case mp1.COIN_EVT_TYPE:
		strs = make([]string, len(responses))
		for i, r := range responses {
			strs[i] = fmt.Sprintf("%d Coins", r)
		}
		selection := widget.NewSelect(strs, func(s string) {
			i, _ := strconv.Atoi(s)
			g.SetResponse(i)
		})
		return selection
	case mp1.PLAYER_EVT_TYPE:
		strs = []string{"", "", "", ""}
		for i, p := range g.Players {
			strs[i] = p.Char
		}
		selection := widget.NewSelect(strs, func(s string) {
			for i, p := range g.Players {
				if s == p.Char {
					g.SetResponse(i)
					return
				}
			}
			panic("Should Be Unreachable")
		})
		return selection
	case mp1.MULTIWIN_PLAYER_EVT_TYPE:
		strs = []string{"None", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""}
		for i := 1; i < len(strs); i++ {
			str := ""
			for p := 0; p < len(g.Players); p++ {
				if i&(1<<p) > 0 {
					str += g.Players[p].Char + " & "
				}
			}
			str = str[:len(str)-3] //Remove last " & "
			strs[i] = str
		}
		selection := widget.NewSelect(strs, func(s string) {
			playerNames := strings.Split(s, " & ")
			res := 0
			for _, pName := range playerNames {
				for p := 0; p < len(g.Players); p++ {
					if pName == g.Players[p].Char {
						res |= (1 << p)
						break
					}
				}
			}
			g.SetResponse(res)
		})
		return selection
	case mp1.CHAINSPACE_EVT_TYPE:
		for _, res := range responses {
			cs := res.(mp1.ChainSpace)
			showSpace(spaceMap, cs)
		}
	}

	return widget.NewLabel("Nothing to put here")
}
