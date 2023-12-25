package commands

import (
	"bytes"
	"log"
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

func (comm *Commands) HandlePush(req flep.Request) {
	offset, err := comm.queueStore.Write(bytes.NewReader(req.Args[1]))
	if err != nil {
		log.Println("Error:", err)
		return
	}
	log.Println("Offset:", offset)
}

func (comm *Commands) HandlePick(req flep.Request) {
	offset, err := strconv.Atoi(string(req.Args[1]))
	if err != nil {
		log.Println("Error:", err)
		return
	}

	var buf bytes.Buffer
	err = comm.queueStore.Read(uint64(offset), &buf)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	log.Println("Message:", buf.String())
}