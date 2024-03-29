package store

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/alainrk/flemq/common"
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
	mu sync.RWMutex
	// offset is the next offset to be written
	offset uint64
	// folderPath is the root folder where the data and index files are stored
	folderPath string
	// dataFile is the file where the data is stored
	dataFile *os.File
	// indexFile is the file where the index is stored
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

	dataFile, err := os.OpenFile(dataFilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	indexFile, err := os.OpenFile(indexFilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	s := &FileQueue{
		mu:         sync.RWMutex{},
		folderPath: folderPath,
		dataFile:   dataFile,
		indexFile:  indexFile,
	}

	s.offset, err = s.getOffsetAtStartup()
	if err != nil {
		panic(err)
	}

	return s
}

// Write writes the data from the reader into the queue.
func (s *FileQueue) Write(reader io.Reader) (offset uint64, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	o, err := s.writeItem(reader)
	if err != nil {
		return 0, err
	}

	s.offset++
	return o, nil
}

// Read reads the data at the given offset into the writer.
func (s *FileQueue) Read(offset uint64, writer io.Writer) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// If the offset is equal or greater than the current offset (the next to be writter)
	if offset >= s.offset {
		return common.OffsetNotFoundError{Err: fmt.Errorf("offset %d not found", offset)}
	}

	return s.getItem(offset, writer)
}

// Close closes the queue files.
func (s *FileQueue) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.dataFile.Close()
	if err != nil {
		return err
	}

	err = s.indexFile.Close()
	if err != nil {
		return err
	}

	return nil
}

// getItem reads from the index file the offset and size of the data
// and then reads the data from the data file.
// It writes the data to the given writer.
// Note: It does not lock the mutex, it is the caller responsibility.
func (s *FileQueue) getItem(offset uint64, w io.Writer) error {
	// Read the offset and size
	var (
		dataOffset uint64
		dataSize   uint64
		err        error
	)

	// Get data position and size from the index file
	dataPos := offset * offsetMapEntrySize
	_, err = s.indexFile.Seek(int64(dataPos), 0)
	if err != nil {
		return err
	}

	err = binary.Read(s.indexFile, binary.BigEndian, &dataOffset)
	if err != nil {
		return err
	}

	err = binary.Read(s.indexFile, binary.BigEndian, &dataSize)
	if err != nil {
		return err
	}

	_, err = s.dataFile.Seek(int64(dataOffset), 0)
	if err != nil {
		return err
	}

	n, err := io.CopyN(w, s.dataFile, int64(dataSize))
	if err != nil {
		return err
	}
	if n != int64(dataSize) {
		return io.ErrShortWrite
	}

	return nil
}

// writeItem writes the data to the data file and the offset and size to the index file.
// It returns the offset of the data in the index file.
// Note: It does not lock the mutex, it is the caller responsibility.
func (s *FileQueue) writeItem(r io.Reader) (offset uint64, err error) {
	// TODO: This can be done way better, like storing it in memory and update it accordingly at each write
	// Get current data file offset through size of the data file
	stat, err := s.dataFile.Stat()
	if err != nil {
		return 0, err
	}
	dataOffset := uint64(stat.Size())

	// Write data to the data file
	dataSize, err := io.Copy(s.dataFile, r)
	if err != nil {
		return 0, err
	}

	// Write data offset and size to the index file
	err = binary.Write(s.indexFile, binary.BigEndian, dataOffset)
	if err != nil {
		return 0, err
	}
	err = binary.Write(s.indexFile, binary.BigEndian, uint64(dataSize))
	if err != nil {
		return 0, err
	}

	return s.offset, nil
}

// getOffsetAtStartup returns the offset at startup
func (s *FileQueue) getOffsetAtStartup() (uint64, error) {
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
