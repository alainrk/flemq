package broker

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func subscriberTest(b *Broker[string], id int, res chan int) {
	s := b.Subscribe()
	i := 0
	for m := range s {
		fmt.Printf("Client %d got message: %v\n", id, m)
		i++
	}

	fmt.Printf("Client %d done with %d messages\n", id, i)

	res <- i
	b.Unsubscribe(s)
}

func publisherTest(b *Broker[string], count int) {
	for i := 0; i < count; i++ {
		b.Publish(fmt.Sprintf("msg#%d", i))
		time.Sleep(100 * time.Millisecond)
	}
}

func TestOneSubscriber(t *testing.T) {
	var (
		wg  sync.WaitGroup
		res = make(chan int)
		b   = NewBroker[string]()
	)

	go b.Start()

	wg.Add(1)
	go func() {
		defer wg.Done()
		subscriberTest(b, 1, res)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		publisherTest(b, 5)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(1 * time.Second)
		b.Stop()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		r := <-res
		if r != 5 {
			t.Errorf("Expected 5 messages, got %d", r)
		}
	}()

	wg.Wait()
}
