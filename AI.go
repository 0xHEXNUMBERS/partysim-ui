package main

import (
	"sync"
	"time"

	"github.com/0xhexnumbers/gmcts/v2"
	"github.com/0xhexnumbers/partysim/mp1"
)

type GameState struct {
	g   mp1.Game
	res []mp1.Response
}

func (g GameState) Len() int {
	return len(g.res)
}

func (g GameState) ApplyAction(i int) (gmcts.Game, error) {
	g1 := g.g
	g1.HandleEvent(g.res[i])
	return GameState{g1, g1.Responses()}, nil
}

func (g GameState) Hash() interface{} {
	return g.g
}

func (g GameState) Player() gmcts.Player {
	if g.g.NextEvent != nil {
		return gmcts.Player(g.g.NextEvent.ControllingPlayer())
	}
	return -1
}

func (g GameState) IsTerminal() bool {
	return g.g.NextEvent == nil
}

func (g GameState) Winners() []gmcts.Player {
	if g.g.NextEvent != nil {
		return nil
	}
	mpWinners := g.g.Winners()
	gmctsWinners := make([]gmcts.Player, len(mpWinners))
	for i, w := range mpWinners {
		gmctsWinners[i] = gmcts.Player(w)
	}
	return gmctsWinners
}

func bestMove(g mp1.Game, threads, milliseconds int) mp1.Response {
	mcts := gmcts.NewMCTS(GameState{g, g.Responses()})
	mcts.SetSeed(time.Now().Unix())
	var wait sync.WaitGroup
	wait.Add(threads)
	for i := 0; i < threads; i++ {
		go func() {
			tree := mcts.SpawnTree()
			tree.Search(time.Millisecond * time.Duration(milliseconds))
			mcts.AddTree(tree)
			wait.Done()
		}()
	}
	wait.Wait()
	i, _ := mcts.BestAction()

	return g.Responses()[i]
}
