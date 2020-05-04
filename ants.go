package main

import "math/rand"

type Colony struct {
	α    float64 // pheromonce importance
	β    float64 // heuristic information importance
	τ    float64
	N    int
	dM   [][]int
	visibilities [][]int
	pheromones [][]int
	ants []ant
}
type ant struct {
	loc location
	currentTour []int
}
type location struct {
	x int
	y int
}

// Randomly pace ants on nodes
func (c *Colony) Init(nAnts, n int ) {
	c.dM = createDistanceMatrix(n, 0)

	c.α  = 0.1
	c.β = 0.1

	c.ants = make([]ant, nAnts)
	for _, a := range c.ants {
		a.loc = location{
			x: rand.Intn(len(c.dM)),
			y: rand.Intn(len(c.dM)),
		}
	}

}

func (c *Colony) Train()  {

}
