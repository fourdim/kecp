package kecpchannel

// https://github.com/leolara/conveyor

import (
	"sync"
)

type Channel[T any] struct {
	ch chan T

	closingCh      chan interface{}
	writersWG      sync.WaitGroup
	writersWGMutex sync.Mutex
}

// New creates a Channel
func New[T any]() *Channel[T] {
	return &Channel[T]{
		ch:        make(chan T),
		closingCh: make(chan interface{}),
	}
}

// Read returns the channel to write
func (p *Channel[T]) Read() <-chan T {
	return p.ch
}

// Write into the channel in a different goroutine
func (p *Channel[T]) Write(data T) {
	go func(data T) {
		p.writersWGMutex.Lock()
		p.writersWG.Add(1)
		p.writersWGMutex.Unlock()
		defer p.writersWG.Done()

		select {
		case <-p.closingCh:
			return
		default:
		}

		select {
		case <-p.closingCh:
		case p.ch <- data:
		}
	}(data)
}

// Closes channel, draining any blocked writes
func (p *Channel[T]) Close() {
	close(p.closingCh)

	go func() {
		for range p.ch {
		}
	}()

	p.writersWGMutex.Lock()
	p.writersWG.Wait()
	p.writersWGMutex.Unlock()

	close(p.ch)
}

// CloseWithoutDraining closes channel, without draining any pending writes, this method
// will block until all writes have been unblocked by reads
func (p *Channel[T]) CloseWithoutDraining() {
	close(p.closingCh)

	p.writersWGMutex.Lock()
	p.writersWG.Wait()
	p.writersWGMutex.Unlock()

	close(p.ch)
}
