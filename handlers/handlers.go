package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/alainrk/flemq/flep"
	"github.com/alainrk/flemq/store"
	"github.com/alainrk/flemq/topic"
)

type Topics map[string]topic.DefaultTopic

type Handlers struct {
	topics Topics
}

func NewHandlers(queueStore store.QueueStore) Handlers {
	return Handlers{
		topics: make(Topics),
	}
}

func (comm *Handlers) HandlePush(req flep.Request) (uint64, error) {
	tn := string(req.Args[0])

	// XXX: Auto-creates topic if it doesn't exist for now.
	if _, ok := comm.topics[tn]; !ok {
		comm.topics[tn] = topic.New(tn)
	}

	topic := comm.topics[tn]

	offset, err := topic.Write(bytes.NewReader(req.Args[1]))
	if err != nil {
		return 0, err
	}
	return offset, nil
}

func (comm *Handlers) HandlePick(req flep.Request) ([]byte, error) {
	offset, err := strconv.Atoi(string(req.Args[1]))
	if err != nil {
		return nil, err
	}

	tn := string(req.Args[0])

	// Topic must exist.
	topic, ok := comm.topics[tn]
	if !ok {
		return nil, fmt.Errorf("topic %s does not exist", tn)
	}

	var buf bytes.Buffer
	err = topic.Read(uint64(offset), &buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// HandleSubscribe handles the subscribe command.
// This function under normal circumstances should never return.
// It should:
//   - Send all the messages from the offset to the existing offset.
//   - Listen for new messages and send them as they come in (channel/polling on the map/other...)
func (comm *Handlers) HandleSubscribe(conn net.Conn, req flep.Request) error {
	startingOffset, err := strconv.Atoi(string(req.Args[1]))
	if err != nil {
		return err
	}

	tn := string(req.Args[0])

	// Topic must exist.
	topic, ok := comm.topics[tn]
	if !ok {
		return fmt.Errorf("topic %s does not exist", tn)
	}

	var buf bytes.Buffer
	offset := startingOffset
	for {
		buf.Reset()
		err = topic.Read(uint64(offset), &buf)
		if err != nil {
			if errors.Is(err, store.ErrorTopicOffsetNotFound) {
				break
			}
			return err
		}

		// Send previously received message
		fmt.Printf("Sending previous offset %d: %s\n", offset, buf.Bytes())
		// TODO: Decide how to handle the stream, should we just send everything as as it come or prepend the length? Or maybe use a EOF marker?
		fr := flep.SimpleBytesResponse(buf.Bytes())
		conn.Write(fr)
		offset++
	}

	// Send any other incoming message coming from the topic's broker.
	s := topic.Subscribe()
	for msg := range s {
		// TODO: Decide how to handle the stream, should we just send everything as as it come or prepend the length? Or maybe use a EOF marker?
		fr := flep.SimpleBytesResponse(msg)
		conn.Write(fr)
	}

	// Should never get here.
	return nil
}
