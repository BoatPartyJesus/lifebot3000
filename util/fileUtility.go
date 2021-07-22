package util

import (
	"io/ioutil"
	"os"
)

type FileUtility interface {
	ReadFile(filepath string) ([]byte, error)
	WriteFile(filepath string, fileContents []byte) error
	CreateFile(filepath string) error
	DoesFileExist(filepath string) bool
}

type OsFileUtility struct{}

func (o OsFileUtility) ReadFile(filepath string) ([]byte, error) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (o OsFileUtility) WriteFile(filepath string, fileContents []byte) error {
	err := ioutil.WriteFile(filepath, fileContents, 0644)

	if err != nil {
		return err
	}

	return nil
}

func (o OsFileUtility) CreateFile(filepath string) error {
	_, err := os.Create(filepath)

	if err != nil {
		return err
	}

	return nil
}

func (o OsFileUtility) DoesFileExist(filepath string) bool {
	_, err := os.Stat(filepath)
	return err == nil
}
