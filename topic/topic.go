package topic

import (
	"io"

	"github.com/alainrk/flemq/store"
)

type Topic struct {
	Name  string
	store store.QueueStore
}

func New(name string) *Topic {
	return &Topic{
		Name:  name,
		store: store.NewMemoryQueueStore(),
	}
}

func (t *Topic) Write(reader io.Reader) (offset uint64, err error) {
	// TODO: Implement any needed topic-specific logic here.
	return t.store.Write(reader)
}

func (t *Topic) Read(offset uint64, writer io.Writer) error {
	// TODO: Implement any needed topic-specific logic here.
	return t.store.Read(offset, writer)
}
