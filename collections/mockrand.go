package collections

import "fmt"

type mockrnd struct {
	sequence []int
	index    int
}

func (rnd *mockrnd) Intn(v int) int {
	if rnd.index >= len(rnd.sequence) {
		panic("mocrnd sequence exhausted")
	}
	ans := rnd.sequence[rnd.index] % v
	fmt.Println("RAND from ", rnd.sequence, ": ", ans, ",  idx: ", rnd.index)
	rnd.index++
	return ans
}
