package main

import "fmt"

func main() {
	segmentCount := 50
	segments := make([][2]int, segmentCount)
	bytes := 2340958301
	eachFragment := bytes / segmentCount
	for i, _ := range segments {
		if i == 0 {
			segments[i][0] = 0
		} else {
			segments[i][0] = segments[i-1][1] + 1
		}
		if i < segmentCount-1 {
			segments[i][1] = segments[i][0] + eachFragment
		} else {
			segments[i][1] = bytes - 1
		}
	}
	fmt.Printf("%v", segments)
}
