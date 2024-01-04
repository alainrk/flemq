package topic

import (
	"bytes"
	"io"

	"github.com/alainrk/flemq/broker"
	"github.com/alainrk/flemq/store"
)

type Topic interface {
	Write(reader io.Reader) (offset uint64, err error)
	Read(offset uint64, writer io.Writer) error
	Subscribe() chan []byte
	Unsubscribe() chan []byte
}

// DefaultTopic implements the Topic interface.
// It both stores and retrieve messages from the provided store and
// allows to subscribe to the topic.
type DefaultTopic struct {
	Name   string
	broker broker.Broker[[]byte]
	store  store.QueueStore
}

// New creates a new topic with the given name.
// It also start a broker for the topic.
func New(name string) DefaultTopic {
	broker := broker.New[[]byte](name, false)
	go broker.Start()
	return DefaultTopic{
		Name:   name,
		broker: broker,
		store:  store.NewMemoryQueueStore(),
	}
}

func (t DefaultTopic) Write(reader io.Reader) (offset uint64, err error) {
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
	t.broker.Publish(b)

	return t.store.Write(&buf)
}

func (t DefaultTopic) Read(offset uint64, writer io.Writer) error {
	return t.store.Read(offset, writer)
}

func (t DefaultTopic) Subscribe() chan []byte {
	return t.broker.Subscribe()
}

func (t DefaultTopic) Unsubscribe(c chan []byte) {
	t.broker.Unsubscribe(c)
}
