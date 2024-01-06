package store

import (
	"encoding/binary"
	"io"
	"os"
	"sync"
)

// Needed to store the offset map
const offsetMapEntrySize = 16 // 8 bytes for key + 8 bytes for value

type FileQueue struct {
	dataFile     *os.File
	indexFile    *os.File
	offset       uint64
	offsetMap    map[uint64]uint64
	offsetFile   *os.File
	offsetFileMu sync.Mutex
	mu           sync.Mutex
}

// NewFileQueue creates a new file queue.
// If the given folder does not exist, it will be created.
func NewFileQueue(folderPath string) *FileQueue {
	err := createFolderIfNotExists(folderPath)
	if err != nil {
		panic(err)
	}

	var (
		dataFilePath   = folderPath + "/data"
		indexFilePath  = folderPath + "/index"
		offsetFilePath = folderPath + "/offset"
	)

	dataFile, err := os.OpenFile(dataFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	indexFile, err := os.OpenFile(indexFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	offsetFile, err := os.OpenFile(offsetFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	// Initialize offset from the index file size
	offset, err := offsetFile.Seek(0, io.SeekEnd)
	if err != nil {
		panic(err)
	}

	// Load the offset map from the offset file
	// TODO: We could just load the needed ones (last X, or until a certain offset)
	offsetMap := make(map[uint64]uint64)
	for {
		var key, value uint64
		err := binary.Read(offsetFile, binary.LittleEndian, &key)
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		err = binary.Read(offsetFile, binary.LittleEndian, &value)
		if err != nil {
			panic(err)
		}
		offsetMap[key] = value
	}

	return &FileQueue{
		dataFile:     dataFile,
		indexFile:    indexFile,
		offset:       uint64(offset),
		offsetMap:    offsetMap,
		offsetFile:   offsetFile,
		offsetFileMu: sync.Mutex{},
		mu:           sync.Mutex{},
	}
}

func createFolderIfNotExists(folderPath string) error {
	// Create folder if it does not exist
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err = os.MkdirAll(folderPath, 0700)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *FileQueue) Write(reader io.Reader) (offset uint64, err error) {
	return 0, nil
}
