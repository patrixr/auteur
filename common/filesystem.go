package common

import (
	"fmt"
	"os"

	"github.com/golang-cz/textcase"
)

func Mkdirp(path string) error {
	info, err := os.Stat(path)

	if err == nil {
		if info.IsDir() {
			return nil
		} else {
			return fmt.Errorf("%s already exists and is not a directory", path)
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func Rmdir(path string) error {
	info, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}

	err = os.RemoveAll(path)
	if err != nil {
		return err
	}
	return nil
}

func ToSlug(text string) string {
	return textcase.KebabCase(text)
}
