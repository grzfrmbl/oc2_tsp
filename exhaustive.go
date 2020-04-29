package main

import (
	"fmt"
	"math"
)

// Perm calls f with each permutation of a.
func Perm(a []int, f func([]int)) {
	perm(a, f, 0)
}

// Permute the values at index i to len(a)-1.
func perm(a []int, f func([]int), i int) {
	if i > len(a) {
		f(a)
		return
	}
	perm(a, f, i+1)
	for j := i + 1; j < len(a); j++ {
		a[i], a[j] = a[j], a[i]
		perm(a, f, i+1)
		a[i], a[j] = a[j], a[i]
	}
}

func exhaustiveSearch(src, n int, dist [][]int) {
	// fmt.Println("exhaustive search from ", src, cities[src])
	perms := genPermutations(src, n)
	fmt.Println("there are ", len(perms), " different paths to check")

	var path []int
	var shortest int = math.MaxInt32

	for _, perm := range perms {
		// Check if path ist shorter than path
		d := calcPathDist(perm, dist)
		if d < shortest {
			fmt.Println(path,"new shortest found ", d, "is shorter than",shortest)

			shortest = d
			path = perm
		}
	}
	fmt.Println("shortest path is ", path)
}
func genPermutations(start int, n int) [][]int {
	// Example Three cities:
	// 	Start = 0
	//	0 1 2 0
	// 	0 2 1 0
	// 	Start = 1
	// 	1 0 2 1
	// 	1 2 0 1
	fmt.Println("there are ",Factorial(n-1)," paths to check")
	a := make([][]int, Factorial(n-1))

	// Indices without/start/end
	var left []int
	for i := 0; i < n; i++ {
		if i == start {
			continue
		}
		left = append(left, i)
	}

	tmp := 0
	for perm := range GeneratePermutations(left) {
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

func Factorial(n int) (result int) {
	if n > 0 {
		result = n * Factorial(n-1)
		return result
	}
	return 1
}

func GeneratePermutations(data []int) <-chan []int {
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
