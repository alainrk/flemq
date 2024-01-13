package topic

import (
	"bytes"
	"fmt"
	"io"

	"github.com/alainrk/flemq/broker"
	"github.com/alainrk/flemq/store"
)

const (
	usePersistence     = true
	defaultStoreFolder = "/tmp/flemq"
)

type Topic interface {
	Write(reader io.Reader) (offset uint64, err error)
	Read(offset uint64, writer io.Writer) error
	Subscribe() chan []byte
	Unsubscribe() chan []byte
	Close() error
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
	var s store.QueueStore

	broker := broker.New[[]byte](name, false)
	go broker.Start()

	if usePersistence {
		s = store.NewFileQueue(fmt.Sprintf("%s/%s", defaultStoreFolder, name))
	} else {
		s = store.NewMemoryQueue()
	}

	return DefaultTopic{
		Name:   name,
		broker: broker,
		store:  s,
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

func (t DefaultTopic) Close() error {
	return t.store.Close()
}
