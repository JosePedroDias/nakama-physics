package main

import (
	cp "github.com/jakecoffman/cp/v2"
)

type PGame struct {
	space *cp.Space
	b1    *cp.Body
}

////

func addHollowRect(space *cp.Space, w, h float64) {
	w2 := w / 2
	h2 := h / 2
	sides := []cp.Vector{
		{X: -w2, Y: -h2}, {X: -w2, Y: h2},
		{X: w2, Y: -h2}, {X: w2, Y: h2},
		{X: -w2, Y: -h2}, {X: w2, Y: -h2},
		{X: -w2, Y: h2}, {X: w2, Y: h2},
	}

	for i := 0; i < len(sides); i += 2 {
		var seg *cp.Shape = space.AddShape(cp.NewSegment(space.StaticBody, sides[i], sides[i+1], 0))
		seg.SetElasticity(1)
		seg.SetFriction(1)
	}
}

func addBox(space *cp.Space, size, mass float64) *cp.Body {
	body := space.AddBody(cp.NewBody(mass, cp.MomentForBox(mass, size, size)))
	shape := space.AddShape(cp.NewBox(body, size, size, 0))
	shape.SetElasticity(0)
	shape.SetFriction(0.7)
	return body
}

////

func newPhysicsGame() *PGame {
	space := cp.NewSpace()
	//space.Iterations = 30 // default to 10
	//space.SetGravity(cp.Vector{X: 0, Y: 9.8})
	space.SleepTimeThreshold = 0.5 // if idle for 0.5 secs, goes to sleep
	space.SetCollisionSlop(0.5)

	addHollowRect(space, 50, 40)
	b1 := addBox(space, 2, 1)

	game := &PGame{
		space: space,
		b1:    b1,
	}

	return game
}

func (g *PGame) update(dt float64) {
	g.space.Step(dt)
	//p := g.b1.Position()
	//fmt.Printf("%#v\n", p)
}

/*func main() {
	g := newPhysicsGame()
	for {
		g.update(1 / 30.0)
		time.Sleep(10 * time.Millisecond)
	}
}*/
