package broker

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func subscriberTest(s chan string, id int, res *atomic.Int32) {
	i := 0
	for m := range s {
		i++
		res.Add(1)
		fmt.Printf("Client %d got message: %v\n", id, m)
	}

	fmt.Printf("Client %d done with %d messages\n", id, i)
}

func publisherTest(b *Broker[string], count int) {
	for i := 0; i < count; i++ {
		b.Publish(fmt.Sprintf("msg#%d", i))
	}
}

func TestPublishSubscribe(t *testing.T) {
	var (
		wg     sync.WaitGroup
		res    atomic.Int32
		b      = NewBroker[string]()
		nSub   = 10
		nMsg   = 5
		totMsg = nSub * nMsg
	)

	res.Store(0)

	go b.Start()
	time.Sleep(1 * time.Second)

	for i := 0; i < nSub; i++ {
		s := b.Subscribe()
		wg.Add(1)
		go func(i int, s chan string) {
			defer wg.Done()
			subscriberTest(s, i, &res)
		}(i, s)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		publisherTest(b, nMsg)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			time.Sleep(100 * time.Millisecond)
			v := res.Load()
			if v > int32(totMsg) {
				t.Errorf("Expected %d messages, got %d", totMsg, v)
			}
			if v == int32(totMsg) {
				b.Stop()
				break
			}
		}
	}()

	wg.Wait()

	v := res.Load()
	if v != int32(totMsg) {
		t.Errorf("Expected %d messages, got %d", totMsg, v)
	}
}
