package main

import (
	"fmt"
	"math"

	"github.com/gosuri/uiprogress"
)

func exhaustiveSearch(src, n int, dist [][]int) {
	fmt.Println("exhaustive search, there are ", fac(n-1), " paths to check")

	perms := genPermutations(src, n)

	var path []int
	var shortest int = math.MaxInt32

	for _, perm := range perms {
		// Check if path ist shorter than path
		d := calcPathDist(perm, dist)
		if d < shortest {
			shortest = d
			path = perm
		}
	}
	fmt.Println("shortest path is ", path, "len", shortest)
}

// Do not preallocate all possible paths, be merciful to your memory.
func betteExhaustiveSearch(src, n int, dist [][]int) {
	fmt.Println("better exhaustive search, there are ", fac(n-1), " paths to check")

	var path []int
	var shortest int = math.MaxInt32

	// Fancy
	uiprogress.Start()                   // start rendering
	bar := uiprogress.AddBar(fac(n - 1)) // Add a new bar

	// optionally, append and prepend completion and elapsed time
	bar.AppendCompleted()
	// bar.PrependElapsed()
	bar.AppendElapsed()

	// Indices without/start/end
	left := sliceWithoutSrc(src, n)

	for perm := range permutations(left) {
		bar.Incr()

		// Make a path
		a := make([]int, n+1)
		a[0] = src
		a[n] = src
		k := 1
		for i := 0; i < len(perm); i++ {
			a[k] = perm[i]
			k++
		}

		// Check it
		d := calcPathDist(a, dist)
		if d < shortest {
			shortest = d
			path = a
		}

	}

	fmt.Println("shortest path is ", path, "len", shortest)
}

// This does work, but since all possible paths are generated beforehand
// we may run into a little RAM overflow if the problem is too big...
// In my case (32gb) if n>13
func genPermutations(start int, n int) [][]int {
	a := make([][]int, fac(n-1))

	// Indices without/start/end
	var left []int
	for i := 0; i < n; i++ {
		if i == start {
			continue
		}
		left = append(left, i)
	}

	tmp := 0
	for perm := range permutations(left) {
		a[tmp] = make([]int, n+1)
		a[tmp][0] = start
		a[tmp][n] = start
		k := 1
		for i := 0; i < len(perm); i++ {
			a[tmp][k] = perm[i]
			k++
		}
		tmp++
	}
	return a
}

func fac(n int) (result int) {
	if n > 0 {
		result = n * fac(n-1)
		return result
	}
	return 1
}
func sliceWithoutSrc(v, n int) []int {
	var left []int
	for i := 0; i < n; i++ {
		if i == v {
			continue
		}
		left = append(left, i)
	}
	return left

}
func permutations(data []int) <-chan []int {
	c := make(chan []int)
	go func(c chan []int) {
		defer close(c)
		permutate(c, data)
	}(c)
	return c
}
func permutate(c chan []int, inputs []int) {
	output := make([]int, len(inputs))
	copy(output, inputs)
	c <- output

	size := len(inputs)
	p := make([]int, size+1)
	for i := 0; i < size+1; i++ {
		p[i] = i
	}
	for i := 1; i < size; {
		p[i]--
		j := 0
		if i%2 == 1 {
			j = p[i]
		}
		tmp := inputs[j]
		inputs[j] = inputs[i]
		inputs[i] = tmp
		output := make([]int, len(inputs))
		copy(output, inputs)
		c <- output
		for i = 1; p[i] == 0; i++ {
			p[i] = i
		}
	}
}
