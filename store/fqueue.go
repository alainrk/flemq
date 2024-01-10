package store

import (
	"io"
	"os"
	"sync"
)

/*
	Index file (16 bytes each entry (offset)):
	Offset i -> i * 8+8 (pos, size in data file)

	Data file:
	Offset p -> data [p, p+size-1]
*/

// Needed to store the (data_offset, data_size) for each entry in the index file
// 8 bytes for data_file_offset + 8 bytes for the size of the entry
const offsetMapEntrySize uint64 = 16

type FileQueue struct {
	mu        sync.RWMutex
	dataFile  *os.File
	offset    uint64
	offsetMap map[uint64][2]uint64
	indexFile *os.File
}

// NewFileQueue creates a new file queue.
// If the given folder does not exist, it will be created.
func NewFileQueue(folderPath string) *FileQueue {
	err := createFolderIfNotExists(folderPath)
	if err != nil {
		panic(err)
	}

	var (
		dataFilePath  = folderPath + "/data"
		indexFilePath = folderPath + "/index"
	)

	dataFile, err := os.OpenFile(dataFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	indexFile, err := os.OpenFile(indexFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	// TODO: Initialize the offsetMap
	s := &FileQueue{
		dataFile:  dataFile,
		indexFile: indexFile,
		mu:        sync.RWMutex{},
	}

	s.offset, err = s.getOffsetAtStartup()
	if err != nil {
		panic(err)
	}

	return s
}

// TODO
func (s *FileQueue) Write(reader io.Reader) (offset uint64, err error) {
	return 0, nil
}

// TODO
func (s *FileQueue) Read(offset uint64, writer io.Writer) error {
	return nil
}

// getOffsetAtStartup returns the offset at startup atomically
func (s *FileQueue) getOffsetAtStartup() (uint64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Get indexFile size
	stat, err := s.indexFile.Stat()
	if err != nil {
		return 0, err
	}

	offset := uint64(stat.Size()) / offsetMapEntrySize
	return offset, nil
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
