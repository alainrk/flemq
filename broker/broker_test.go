package broker

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func subscriberTest(s chan string, id int, res *atomic.Int32, subReady *sync.WaitGroup) {
	i := 0
	subReady.Done()
	for m := range s {
		i++
		res.Add(1)
		fmt.Printf("Client %d got message: %v\n", id, m)
	}

	fmt.Printf("Client %d done with %d messages\n", id, i)
}

func publisherTest(b *Broker[string], count int) {
	for i := 0; i < count; i++ {
		m := fmt.Sprintf("msg#%d", i)
		fmt.Printf("Publishing message: %s\n", m)
		b.Publish(m)
	}
}

func TestPublishSubscribe(t *testing.T) {
	var (
		wg       sync.WaitGroup
		subReady sync.WaitGroup
		res      atomic.Int32
		b        = NewBroker[string]("test")
		nSub     = 20
		nMsg     = 90
		totMsg   = nSub * nMsg
	)

	res.Store(0)

	go b.Start()
	// Run go scheduler to allow broker to start.
	runtime.Gosched()

	wg.Add(nSub)
	subReady.Add(nSub)

	for i := 0; i < nSub; i++ {
		s := b.Subscribe()
		go func(i int, s chan string) {
			defer wg.Done()
			subscriberTest(s, i, &res, &subReady)
		}(i, s)
	}

	fmt.Println("Waiting for subscribers to be ready...")
	subReady.Wait()
	fmt.Println("All subscribers ready!")

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
				b.Stop()
				break
			}
			if v == int32(totMsg) {
				b.Stop()
				break
			}
			fmt.Printf("Got %d messages so far...\n", v)
		}
	}()

	wg.Wait()

	v := res.Load()
	if v != int32(totMsg) {
		t.Errorf("Expected %d messages, got %d", totMsg, v)
	}
}
