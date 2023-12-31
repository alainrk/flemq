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
	for range s {
		i++
		res.Add(1)
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
	tests := []struct {
		name string
		nSub int
		nMsg int
	}{
		{"1 message 1 subscriber", 1, 1},
		{"1 message 10 subscribers", 10, 1},
		{"10 messages 1 subscriber", 1, 10},
		{"10 messages 10 subscribers", 10, 10},
		{"100 messages 1 subscriber", 1, 100},
		{"100 messages 10 subscribers", 10, 100},
		{"100 messages 100 subscribers", 100, 100},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			subTestPublishSubscribe(t, tc.name, tc.nSub, tc.nMsg)
		})
	}
}

func subTestPublishSubscribe(t *testing.T, name string, nSub, nMsg int) {
	var (
		wg       sync.WaitGroup
		subReady sync.WaitGroup
		res      atomic.Int32
		b        = NewBroker[string]("test")
		totMsg   = nSub * nMsg
	)

	fmt.Printf("Running test: %s\n", name)

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

	// Maybe not strictly necessary, but wait for all subscribers to be ready.
	subReady.Wait()

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
		}
	}()

	wg.Wait()

	v := res.Load()
	if v != int32(totMsg) {
		t.Errorf("Expected %d messages, got %d", totMsg, v)
	}
}
