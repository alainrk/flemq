package broker

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func subscriberTest(t *testing.T, s chan string, id int, res *atomic.Int32, subReady *sync.WaitGroup) {
	i := 0
	subReady.Done()
	for range s {
		// for m := range s {
		i++
		res.Add(1)
		// t.Logf("Client %d got message: %v\n", id, m)
	}
	// t.Logf("Client %d done with %d messages\n", id, i)
}

func publisherTest(t *testing.T, b Broker[string], count int) {
	for i := 0; i < count; i++ {
		m := fmt.Sprintf("msg#%d", i)
		t.Logf("Publishing message: %s\n", m)
		b.Publish(m)
	}
}

func TestPublishSubscribeBlocking(t *testing.T) {
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

			subTestPublishSubscribeBlocking(t, tc.name, tc.nSub, tc.nMsg)
		})
	}
}

func subTestPublishSubscribeBlocking(t *testing.T, name string, nSub, nMsg int) {
	var (
		wg       sync.WaitGroup
		subReady sync.WaitGroup
		res      atomic.Int32
		b        = New[string]("test", true)
		totMsg   = nSub * nMsg
	)

	t.Logf("Running test: %s\n", name)

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
			subscriberTest(t, s, i, &res, &subReady)
		}(i, s)
	}

	// Maybe not strictly necessary, but wait for all subscribers to be ready.
	subReady.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()
		publisherTest(t, b, nMsg)
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
			t.Logf("Got %d messages so far...\n", v)
		}
	}()

	wg.Wait()

	v := res.Load()
	if v != int32(totMsg) {
		t.Errorf("Expected %d messages, got %d", totMsg, v)
	}
}
