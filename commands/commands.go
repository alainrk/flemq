package commands

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/alainrk/flemq/flep"
	"github.com/alainrk/flemq/store"
)

type Commands struct {
	topics map[string]*Topic
}

type Topic struct {
	Name  string
	store store.QueueStore
}

func NewCommands(queueStore store.QueueStore) Commands {
	return Commands{
		topics: make(map[string]*Topic),
	}
}

func (comm *Commands) HandlePush(req flep.Request) (uint64, error) {
	tn := string(req.Args[0])

	// XXX: Auto-create topic if it doesn't exist for now.
	if _, ok := comm.topics[tn]; !ok {
		comm.topics[tn] = &Topic{
			Name:  tn,
			store: store.NewMemoryQueueStore(),
		}
	}

	topic := comm.topics[tn]

	offset, err := topic.store.Write(bytes.NewReader(req.Args[1]))
	if err != nil {
		return 0, err
	}
	return offset, nil
}

func (comm *Commands) HandlePick(req flep.Request) ([]byte, error) {
	var topic *Topic
	var ok bool
	var buf bytes.Buffer

	offset, err := strconv.Atoi(string(req.Args[1]))
	if err != nil {
		return nil, err
	}

	tn := string(req.Args[0])

	// Topic must exist.
	if topic, ok = comm.topics[tn]; !ok {
		return nil, fmt.Errorf("topic %s does not exist", tn)
	}

	err = topic.store.Read(uint64(offset), &buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// TODO: Implement this.
// - Send all the messages from the offset to the existing offset.
// - Listen for new messages and send them as they come in (channel/polling on the map/other...)
func (comm *Commands) HandleSubscribe(conn net.Conn, req flep.Request) error {
	var topic *Topic
	var ok bool
	var buf bytes.Buffer

	offset, err := strconv.Atoi(string(req.Args[1]))
	if err != nil {
		return err
	}

	tn := string(req.Args[0])

	// Topic must exist.
	if topic, ok = comm.topics[tn]; !ok {
		return fmt.Errorf("topic %s does not exist", tn)
	}

	for i := offset; ; i++ {
		err = topic.store.Read(uint64(offset), &buf)
		if err != nil {
			// if err == store.ErrorTopicOffsetNotFound {
			if errors.Is(err, store.ErrorTopicOffsetNotFound) {
				// TODO: Remove this print
				fmt.Printf("offset %d not found\n", i)
				break
			}
			return err
		}

		conn.Write(buf.Bytes())
	}
	return nil
}
