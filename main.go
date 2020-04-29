package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

const (
	MAX_DIST = 10
	MIN_DIST = 1
)

func main() {
	a := loadTestMatrix()
	fmt.Println(a)
	exhaustiveSearch(0, len(a[0]),a)

	return
	generateProblemInstance(5)
}

func generateProblemInstance(n int) {
	d := createDistanceMatrix(n, time.Now().UnixNano())
	// c := createCities(n)

	// dump(d, c)

	exhaustiveSearch(0, n, d)

}

// 	Test
// 	T	N	M
// T	0	7	8
// N	7	0	4
// M	8	4	0
//
//  Start T    7    4    8
//          T -> N -> M -> T
// 	       8    4    7
//          T -> M -> N -> T
//
func createDistanceMatrix(n int, source int64) [][]int {
	r := rand.New(rand.NewSource(source))
	a := make([][]int, n)
	for i := range a {
		a[i] = make([]int, n)
	}

	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[i]); j++ {
			if i == j {
				// Diagonal
				a[i][j] = 0
			} else if j > i {
				a[i][j] = r.Intn(MAX_DIST-MIN_DIST) + MIN_DIST
			} else {
				// Symmetric
				a[i][j] = a[j][i]
			}
		}
	}

	return a
}
func calcPathDist(path []int, dist [][]int) int {

	d := 0
	for i := 0; i < len(path)-1; i++ {
		d += dist[path[i]][path[i+1]]
	}
	fmt.Println(d)
	return d

}
func dump(d [][]int, c []string) {
	fmt.Println("\t \t", c)
	for i := 0; i < len(d); i++ {
		fmt.Println(c[i], "\t \t", d[i])
	}
}
func loadTestMatrix() [][]int{
	clientsFile, err := os.OpenFile("testData2.json", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer clientsFile.Close()

	d, _ := ioutil.ReadAll(clientsFile)
	var a [][]int

	json.Unmarshal(d, &a)

	return a
}
