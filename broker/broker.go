package broker

import (
	"log"
)

type Broker[T any] struct {
	name      string
	stopCh    chan struct{}
	publishCh chan T
	subCh     chan chan T
	unsubCh   chan chan T
}

func NewBroker[T any](name string) *Broker[T] {
	return &Broker[T]{
		name:      name,
		stopCh:    make(chan struct{}),
		publishCh: make(chan T, 1),
		subCh:     make(chan chan T, 1),
		unsubCh:   make(chan chan T, 1),
	}
}

func (b *Broker[T]) Start() {
	log.Printf("[Broker %s] Starting...\n", b.name)
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
			// TODO: Experiment with non-blocking sends.
			//
			// i := 0
			// log.Printf("[Broker %s] Count Start - Sending \"%v\" to %d channels\n", b.name, msg, len(subs))
			// for msgCh := range subs {
			// 	// Select + default is a pattern non-blocking sends.
			// 	// Non-blocking send to avoid blocking the broker.
			// 	select {
			// 	case msgCh <- msg:
			// 		i++
			// 	default:
			// 		// Client not listening yet, drop the message.
			// 		log.Printf("[Broker %s] Count - Skipping send to channel\n", b.name)
			// 	}
			// }
			// log.Printf("[Broker %s] Count End \"%v\" to channel %d times\n", b.name, msg, i)

			for msgCh := range subs {
				msgCh <- msg
			}
		}
	}
}

func (b *Broker[T]) Stop() {
	log.Printf("[Broker %s] Stopping...\n", b.name)
	close(b.stopCh)
}

func (b *Broker[T]) Subscribe() chan T {
	msgCh := make(chan T, 5)
	b.subCh <- msgCh
	return msgCh
}

func (b *Broker[T]) Unsubscribe(msgCh chan T) {
	b.unsubCh <- msgCh
}

func (b *Broker[T]) Publish(msg T) {
	b.publishCh <- msg
}
