package storage

import (
	"os"
)

type FileRepository struct {
}

func New() *FileRepository {
	return &FileRepository{}
}

func (r *FileRepository) Save(bytes []byte, key string) error {

	file, err := os.Create(key)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(bytes)
	return err
}
