package main

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/gosuri/uiprogress"
)

// First try, overflows memory on large problems
// since all permutations are created and stored in memory.
func exhaustiveSearch(src, n int, dist [][]int) {
	fmt.Println("exhaustive search, there are ", fac(n-1), " paths to check")

	perms := genPermutations(src, n)

	var (
		path     []int
		shortest int = math.MaxInt32
	)

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
	fmt.Println("better exhaustive search, there are ", fac(n-1), " paths to check \nsearching ...\n")

	// These variables are allocated here to increase performance.
	// If we do this every iteration the GC strain is quite noticeable.
	var (
		path         []int
		shortest     int = math.MaxInt32
		shortestTemp int
		a            = make([]int, n+1)
		rate         float64
		tLeft        float64
	)

	// Uhhh yes a fancy progressbar
	uiprogress.Start()
	bar := uiprogress.AddBar(fac(n - 1))
	bar.Width = 42
	// bar.PrependElapsed()
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		tLeft = float64(bar.Total-bar.Current()) / rate
		return "\tpaths/s\t\t" + strconv.FormatFloat(rate, 'f', 0, 64) + "\n" +
		    "\ttime left\t" + time.Duration(tLeft*1000000000).String() + "\n" +
		    "\tshortest\t" + strconv.Itoa(shortest) + "\n" +
		    " " + strconv.FormatFloat(time.Since(b.TimeStarted).Seconds(), 'f', 0, 64) + "s"
	})
	bar.AppendCompleted()

	// Indices without/start/end
	left := sliceWithoutSrc(src, n)

	for perm := range permutations(left) {
		bar.Incr()

		// Make a path

		a[0] = src
		a[n] = src

		for i := 0; i < len(perm); i++ {
			a[i+1] = perm[i]
		}

		// Check it

		shortestTemp = calcPathDist(a, dist)
		if shortestTemp < shortest {
			shortest = shortestTemp
			path = a
		}

		// Compute the current computation rate every now and then

		if bar.Current()%100000 == 0 {
			rate = float64(bar.Current()) / time.Since(bar.TimeStarted).Seconds()
		}
	}

	fmt.Println("shortest path is ", path, "len", shortest, " took ", time.Since(bar.TimeStarted))
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
