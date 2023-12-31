package broker

import (
	"log"
)

type Broker[T any] interface {
	Start()
	Stop()
	Subscribe() chan T
	Unsubscribe(chan T)
	Publish(T)
}

type DefaultBroker[T any] struct {
	name      string
	blocking  bool
	stopCh    chan struct{}
	publishCh chan T
	subCh     chan chan T
	unsubCh   chan chan T
}

// New creates a new DefaultBroker.
// If set to blocking, the send will block until the subscriber read the message.
//
//	Be careful with this setting as can slow down the broker, the publisher or even block the broker.
//
// If set to non-blocking, the send will drop the message if the subscriber is not ready.
func New[T any](name string, blocking bool) DefaultBroker[T] {
	return DefaultBroker[T]{
		name:      name,
		blocking:  blocking,
		stopCh:    make(chan struct{}),
		publishCh: make(chan T, 1),
		subCh:     make(chan chan T, 1),
		unsubCh:   make(chan chan T, 1),
	}
}

func (b DefaultBroker[T]) Start() {
	log.Printf("[Broker %s] Starting...\n", b.name)

	if b.blocking {
		b.startBlockingLoop()
	} else {
		b.startNonBlockingLoop()
	}
}

func (b DefaultBroker[T]) startNonBlockingLoop() {
	subs := map[chan T]struct{}{}
	for {
		select {
		case <-b.stopCh:
			// Close all the channels.
			for msgCh := range subs {
				log.Printf("[Broker %s] Closing all the channels\n", b.name)
				close(msgCh)
			}
			return
		case msgCh := <-b.subCh:
			subs[msgCh] = struct{}{}
		case msgCh := <-b.unsubCh:
			delete(subs, msgCh)
		case msg := <-b.publishCh:
			i := 0
			for msgCh := range subs {
				// Select + default is a pattern non-blocking sends.
				// Non-blocking send to avoid blocking the broker.
				select {
				case msgCh <- msg:
					i++
				default:
				}
			}
		}
	}
}

func (b DefaultBroker[T]) startBlockingLoop() {
	subs := map[chan T]struct{}{}
	for {
		select {
		case <-b.stopCh:
			// Close all the channels.
			for msgCh := range subs {
				log.Printf("[Broker %s] Closing all the channels\n", b.name)
				close(msgCh)
			}
			return
		case msgCh := <-b.subCh:
			subs[msgCh] = struct{}{}
		case msgCh := <-b.unsubCh:
			delete(subs, msgCh)
		case msg := <-b.publishCh:
			for msgCh := range subs {
				msgCh <- msg
			}
		}
	}
}

func (b DefaultBroker[T]) Stop() {
	log.Printf("[Broker %s] Stopping...\n", b.name)
	close(b.stopCh)
}

func (b DefaultBroker[T]) Subscribe() chan T {
	msgCh := make(chan T, 5)
	b.subCh <- msgCh
	return msgCh
}

func (b DefaultBroker[T]) Unsubscribe(msgCh chan T) {
	b.unsubCh <- msgCh
}

func (b DefaultBroker[T]) Publish(msg T) {
	b.publishCh <- msg
}
