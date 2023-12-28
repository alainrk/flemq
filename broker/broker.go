package broker

type Broker[T any] struct {
	stopCh    chan struct{}
	publishCh chan T
	subCh     chan chan T
	unsubCh   chan chan T
}

func NewBroker[T any]() *Broker[T] {
	return &Broker[T]{
		stopCh:    make(chan struct{}),
		publishCh: make(chan T, 1),
		subCh:     make(chan chan T, 1),
		unsubCh:   make(chan chan T, 1),
	}
}

func (b *Broker[T]) Start() {
	subs := map[chan T]struct{}{}
	for {
		select {
		case <-b.stopCh:
			// Close all the channels.
			for msgCh := range subs {
				close(msgCh)
			}
			return
		case msgCh := <-b.subCh:
			subs[msgCh] = struct{}{}
		case msgCh := <-b.unsubCh:
			delete(subs, msgCh)
		case msg := <-b.publishCh:
			for msgCh := range subs {
				// Select + default is a pattern non-blocking sends.
				// Non-blocking send to avoid blocking the broker.
				select {
				case msgCh <- msg:
				default:
				}
			}
		}
	}
}

func (b *Broker[T]) Stop() {
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
