package topic

import (
	"bytes"
	"io"

	"github.com/alainrk/flemq/broker"
	"github.com/alainrk/flemq/store"
)

type Topic struct {
	Name   string
	Broker *broker.Broker[[]byte]
	store  store.QueueStore
}

// New creates a new topic with the given name.
// It also start a broker for the topic.
func New(name string) *Topic {
	broker := broker.NewBroker[[]byte](name, false)
	go broker.Start()
	return &Topic{
		Name:   name,
		Broker: broker,
		store:  store.NewMemoryQueueStore(),
	}
}

func (t *Topic) Write(reader io.Reader) (offset uint64, err error) {
	// TODO: Maybe I can do this stuff in parallel?
	// Also I'm converting the reader into bytes twice,
	// on the other hand if the I get huge messages is
	// a problem for memory until I don't go with io.Copy anyway

	var buf bytes.Buffer
	tee := io.TeeReader(reader, &buf)

	// NOTE: First read from the tee, otherwise the buffer will be empty.
	b, err := io.ReadAll(tee)
	if err != nil {
		return 0, err
	}
	t.Broker.Publish(b)

	return t.store.Write(&buf)
}

func (t *Topic) Read(offset uint64, writer io.Writer) error {
	return t.store.Read(offset, writer)
}
