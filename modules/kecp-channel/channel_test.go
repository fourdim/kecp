package kecpchannel_test

// https://gist.github.com/leolara/f6fb5dfc04d64947487f16764d6b37b6
// In order to detect race conditions run the test with:
// go test -cpu=1,9,55,99 -race -count=100 -failfast

import (
	"math/big"
	"sync"
	"testing"

	. "github.com/fourdim/kecp/modules/kecp-channel"
)

func run(p *Channel[big.Int], n int) {
	for i := 0; i < n; i++ {
		p.Write(*big.NewInt(int64(i)))
	}
}

func TestSimple(t *testing.T) {
	consumer := func(pub *Channel[big.Int], n int, wg *sync.WaitGroup, result chan *big.Int) {
		ch := pub.Read()

		acc := big.NewInt(0)

		for i := 0; i < n; i++ {
			val := <-ch
			t.Log(&val)
			acc.Add(acc, &val)
		}

		wg.Done()
		result <- acc
	}

	producer := func(pub *Channel[big.Int], n int, wg *sync.WaitGroup) {
		run(pub, n)
		wg.Done()
	}

	precalc := func(n int) *big.Int {
		acc := big.NewInt(0)
		for i := 0; i < n; i++ {
			acc.Add(acc, big.NewInt(int64(i)))
		}

		return acc
	}

	p := New[big.Int]()
	var wg sync.WaitGroup
	resultCh := make(chan *big.Int)

	wg.Add(2)
	go consumer(p, 100, &wg, resultCh)
	go producer(p, 100, &wg)
	wg.Wait()
	p.CloseWithoutDraining()

	result := <-resultCh
	t.Log(result)
	t.Log(precalc(100))

	if result.Cmp(precalc(100)) != 0 {
		t.Error("wrong result")
	}
}

func TestIntermediate(t *testing.T) {
	consumer := func(pub *Channel[big.Int], n int, wg *sync.WaitGroup, result chan *big.Int) {
		ch := pub.Read()

		acc := big.NewInt(0)

		for i := 0; i < n; i++ {
			val := <-ch
			t.Log(&val)
			acc.Add(acc, &val)
		}

		wg.Done()
		result <- acc
	}

	producer := func(pub *Channel[big.Int], n int, wg *sync.WaitGroup) {
		run(pub, n)
		wg.Done()
	}

	p := New[big.Int]()
	var wg sync.WaitGroup
	resultCh := make(chan *big.Int)

	wg.Add(3)
	go consumer(p, 100, &wg, resultCh)
	go producer(p, 100, &wg)
	go producer(p, 100, &wg)

	<-resultCh
	p.Close()

	wg.Wait()
}
