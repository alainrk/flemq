package store

import (
	"os"
	"sync"
)

// Needed to store the (data_offset, data_size) for each entry in the index file
// 8 bytes for data_file_offset + 8 bytes for the size of the entry
const offsetMapEntrySize = 16

type FileQueue struct {
	mu       sync.RWMutex
	dataFile *os.File

	offsetMu  sync.RWMutex
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

	// TODO: Initialize the current offset depending on the index file size
	offset := 0

	// TODO: We could just load the needed ones (last X, or until a certain offset)

	return &FileQueue{
		dataFile:  dataFile,
		indexFile: indexFile,
		offset:    uint64(offset),
		// offsetMap: offsetMap,
		offsetMu: sync.RWMutex{},
		mu:       sync.RWMutex{},
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

// Read reads the data from the given offset and writes it to the given writer.
// XXX: For now it's based on the assumption that the requested data are already in memory.
// func (s *FileQueue) Read(offset uint64, writer io.Writer) error {
// 	s.mu.RLock()
// 	defer s.mu.RUnlock()

// 	// Retrieve data offset from the offset map
// 	d, ok := s.offsetMap[offset]
// 	if !ok {
// 		return errors.New("offset not found")
// 	}

// 	offset, size := d[0], d[1]

// 	// Read data from the data file
// 	_, err := s.dataFile.Seek(int64(offset), io.SeekStart)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = io.Copy(writer, io.LimitReader(s.dataFile, 1))
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
