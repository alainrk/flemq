package commands

import (
	"bytes"
	"strconv"

	"github.com/alainrk/flemq/flep"
	"github.com/alainrk/flemq/store"
)

type Commands struct {
	queueStore store.QueueStore
}

func NewCommands(queueStore store.QueueStore) Commands {
	return Commands{
		queueStore: queueStore,
	}
}

func (comm *Commands) HandlePush(req flep.Request) (uint64, error) {
	offset, err := comm.queueStore.Write(bytes.NewReader(req.Args[1]))
	if err != nil {
		return 0, err
	}
	return offset, nil
}

func (comm *Commands) HandlePick(req flep.Request) ([]byte, error) {
	offset, err := strconv.Atoi(string(req.Args[1]))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = comm.queueStore.Read(uint64(offset), &buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
