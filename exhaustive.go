package main

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/gosuri/uiprogress"
)

func search(dist [][]int, fast bool) {
	if fast {
		fastExhaustive(len(dist[0]), dist)
	} else {
		betterExhaustive(len(dist[0]), dist)
	}

}

// First try, overflows memory on large problems
// since all permutations are created and stored in memory.
func exhaustiveSearch(src, n int, dist [][]int) {
	fmt.Println("exhaustive search, there are ", fac(n), " paths to check")

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
func betterExhaustive(n int, dist [][]int) {
	fmt.Println("better exhaustive search, there are ", fac(n), " paths to check \nsearching ...\n")

	var (
		perm         = make([]int, 0, n)
		path         = make([]int, n, n)
		shortest     = math.MaxInt32
		shortestTemp int
		rate         float64
		tLeft        float64
		i, j         int
	)

	for i := 0; i < n; i++ {
		perm = append(perm, i)
	}
	shortest = calcPathDist(perm, dist) // dont forget the first one...

	// Uhhh yes a fancy progressbar

	uiprogress.Start()
	bar := uiprogress.AddBar(fac(n))
	bar.Width = 42
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		tLeft = float64(bar.Total-bar.Current()) / rate
		return "\tpaths/s\t\t" + strconv.FormatFloat(rate, 'f', 0, 64) + "\n" +
		    "\ttime left\t" + time.Duration(tLeft*1000000000).String() + "\n" +
		    "\tshortest\t" + strconv.Itoa(shortest) + " : " + fmt.Sprint(path) + "\n" +
		    " " + strconv.FormatFloat(time.Since(b.TimeStarted).Seconds(), 'f', 0, 64) + "s"
	})
	bar.AppendCompleted()

	for c := 1; c < fac(n); c++ {
		bar.Incr()

		// For the sake of performance permutations are generated inline
		i = (n - 1) - 1
		j = (n - 1)
		for perm[i] > perm[i+1] {
			i--
		}
		for perm[j] < perm[i] {
			j--
		}
		perm[i], perm[j] = perm[j], perm[i]
		j = (n - 1)
		i += 1
		for i < j {
			perm[i], perm[j] = perm[j], perm[i]
			i++
			j--
		}

		// Check the path

		shortestTemp = calcPathDist(perm, dist)
		if shortestTemp < shortest {
			shortest = shortestTemp
			copy(path, perm)
		}

		// Compute the current computation rate every now and then,
		// this and the progress bar does slow down things
		// by about 100% ...
		if bar.Current()%100000 == 0 {
			rate = float64(bar.Current()) / time.Since(bar.TimeStarted).Seconds()
		}

	}

	fmt.Println("shortest path is ", path, "len", shortest, " took ", time.Since(bar.TimeStarted))
}

// Fast version of the above, almost no prints or other fancy crap
func fastExhaustive(n int, dist [][]int) {
	fmt.Println("fast exhaustive search, there are ", fac(n), " paths to check \nsearching ...\n")
	var (
		perm         = make([]int, 0, n)
		path         = make([]int, n, n)
		shortest     = math.MaxInt32
		shortestTemp int
		rate         float64
		tLeft        float64
		i, j, k      int
	)

	for i := 0; i < n; i++ {
		perm = append(perm, i)
	}
	shortest = calcPathDist(perm, dist) // dont forget the first one...

	start := time.Now()
	for c := 1; c < fac(n); c++ {
		// For the sake of performance permutations are generated inline
		i = (n - 1) - 1
		j = (n - 1)
		for perm[i] > perm[i+1] {
			i--
		}
		for perm[j] < perm[i] {
			j--
		}
		perm[i], perm[j] = perm[j], perm[i]
		j = (n - 1)
		i += 1
		for i < j {
			perm[i], perm[j] = perm[j], perm[i]
			i++
			j--
		}

		// Check the path

		shortestTemp = 0
		for k = 0; k < len(path)-1; k++ {
			shortestTemp += dist[perm[k]][perm[k+1]]
		}
		shortestTemp += dist[perm[len(perm)-1]][0]

		if shortestTemp < shortest {
			shortest = shortestTemp
			copy(path, perm)
		}

		// A little clue about expected time might be helpful...
		if c == 10_000_000 {
			rate = float64(c) / time.Since(start).Seconds()
			tLeft = float64(fac(n)-c) / rate
			fmt.Println("estimated time left\t" + time.Duration(tLeft*1000000000).String())
		}

	}

	fmt.Println("shortest path is ", path, "len", shortest, " took ", time.Since(start))
}

// This does work, but since all possible paths are generated beforehand
// we may run into a little RAM overflow if the problem is too big...
func genPermutations(start int, n int) [][]int {
	a := make([][]int, fac(n))

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

// Based on on the QuickPerm algorithm,
// unfortunately the overhead of the go runtime mitigate
// any performance improvements compared to a simple inline heaps algorithm.
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
