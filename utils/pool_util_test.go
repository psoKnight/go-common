package utils

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	gopool := NewGoPool(WithMaxLimit(10))

	defer gopool.Wait()

	for i := 0; i <= 100; i++ {
		gopool.Submit(func() {

			fmt.Println(fmt.Sprintf("Num: %d.", rand.Int()))

			time.Sleep(time.Duration(10) * time.Second)
		})
	}
}
