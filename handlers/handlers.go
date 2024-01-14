package handlers

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/alainrk/flemq/common"
	"github.com/alainrk/flemq/config"
	"github.com/alainrk/flemq/flep"
	"github.com/alainrk/flemq/topic"
)

type Handlers struct {
	config config.Config
	topics map[string]topic.DefaultTopic
}

func NewHandlers(c config.Config) Handlers {
	t := topic.RestoreDefaultTopics(c.Store.Folder)

	return Handlers{
		config: c,
		topics: t,
	}
}

func (h *Handlers) Close() error {
	log.Println("Closing handlers...")
	for _, t := range h.topics {
		err := t.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *Handlers) HandlePush(req flep.Request) (uint64, error) {
	tn := string(req.Args[0])

	// XXX: Auto-creates topic if it doesn't exist for now.
	if _, ok := h.topics[tn]; !ok {
		h.topics[tn] = topic.New(tn, h.config.Store.Folder)
	}

	topic := h.topics[tn]

	offset, err := topic.Write(bytes.NewReader(req.Args[1]))
	if err != nil {
		return 0, err
	}
	return offset, nil
}

func (h *Handlers) HandlePick(req flep.Request) ([]byte, error) {
	offset, err := strconv.Atoi(string(req.Args[1]))
	if err != nil {
		return nil, err
	}

	tn := string(req.Args[0])

	// Topic must exist.
	topic, ok := h.topics[tn]
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
func (h *Handlers) HandleSubscribe(conn net.Conn, req flep.Request) error {
	startingOffset, err := strconv.Atoi(string(req.Args[1]))
	if err != nil {
		return err
	}

	tn := string(req.Args[0])

	// Topic must exist.
	topic, ok := h.topics[tn]
	if !ok {
		return fmt.Errorf("topic %s does not exist", tn)
	}

	var buf bytes.Buffer
	offset := startingOffset
	for {
		buf.Reset()
		err = topic.Read(uint64(offset), &buf)
		if err != nil {
			if _, ok := err.(common.OffsetNotFoundError); ok {
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
