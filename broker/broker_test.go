package broker

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func subscriberTest(b *Broker[string], id int, wg *sync.WaitGroup, res chan int) {
	s := b.Subscribe()
	i := 0
	for m := range s {
		fmt.Printf("Client %d got message: %v\n", id, m)
		i++
	}
	fmt.Printf("Client done")
	res <- i
	b.Unsubscribe(s)
	wg.Done()
}

func publisherTest(b *Broker[string], count int, wg *sync.WaitGroup) {
	for i := 0; i < count; i++ {
		b.Publish(fmt.Sprintf("msg#%d", i))
		time.Sleep(100 * time.Millisecond)
	}
	b.Stop()
	wg.Done()
}

func TestOneSubscriber(t *testing.T) {
	var wg sync.WaitGroup
	res := make(chan int)

	b := NewBroker[string]()
	go b.Start()

	wg.Add(1)
	go subscriberTest(b, 1, &wg, res)

	// Start publishing messages:
	wg.Add(1)
	go publisherTest(b, 5, &wg)

	wg.Add(1)
	go func() {
		r := <-res
		if r != 5 {
			t.Errorf("Expected 5 messages, got %d", r)
		}
		wg.Done()
	}()

	wg.Wait()
}
