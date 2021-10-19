package service

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

// DirWorkFunc ...
type DirWorkFunc func(path os.FileInfo) error

// FSExists ...
func FSExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// FSRemove ...
func FSRemove(path string) error {
	if _, err := os.Stat(path); err == nil {
		return err
	}
	if err := os.Remove(path); err != nil {
		return err
	}
	return nil
}

// FSWorkOnDirectoryInPath ...
func FSWorkOnDirectoryInPath(p string, workFunc DirWorkFunc) error {
	fileInfo, err := ioutil.ReadDir(p)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"operation": "internal.GetDirectoryInPath",
			"path":      p,
			"error":     err}).Error("Scan directory")
		return err
	}
	for _, file := range fileInfo {
		if file.IsDir() {
			if err := workFunc(file); err != nil {
				return err
			}
		}
	}
	return nil
}

// FSWorkOnItemInPath ...
func FSWorkOnItemInPath(p string, workFunc DirWorkFunc) error {
	fileInfo, err := ioutil.ReadDir(p)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"operation": "internal.GetDirectoryInPath",
			"path":      p,
			"error":     err}).Error("Scan directory")
		return err
	}
	for _, file := range fileInfo {
		if err := workFunc(file); err != nil {
			return err
		}
	}
	return nil
}

// FSGetDirectories ...
func FSGetDirectories(p string) (*[]string, error) {
	var dir []string
	fileInfo, err := ioutil.ReadDir(p)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"operation": "internal.FSGetDirectories",
			"path":      p,
			"error":     err}).Error("Scan directory")
		return nil, err
	}
	for _, file := range fileInfo {
		if file.IsDir() {
			dir = append(dir, file.Name())
		}
	}
	return &dir, nil
}

// FSCreateDirectory ...
func FSCreateDirectory(p string) error {
	err := os.MkdirAll(p, 0755)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"path":  p,
			"error": err,
		}).Errorln("Error creating path")
	}
	return err
}

func CreateFile(f string) (*os.File, error) {
	file, err := os.Create(f)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"file":  f,
			"error": err,
		}).Errorln("Error creating file")
	}
	return file, err
}

// LoadFromJSONFile ...
func LoadFromJSONFile(f string, i interface{}) error {
	jsonFile, err := os.Open(f)
	// if we os.Open returns an error then handle it
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"file":  f,
			"error": err,
		}).Errorln("Error opening json file")
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"file":  f,
			"error": err,
		}).Errorln("Error reading data from json file")
	}
	if err = json.Unmarshal([]byte(byteValue), i); err != nil {
		logrus.WithFields(logrus.Fields{
			"file":  f,
			"error": err,
		}).Errorln("Error unmarshalling json file")
	}
	return nil
}

const (
	osRead       = 04
	osWrite      = 02
	osEx         = 01
	osUserShift  = 6
	osGroupShift = 3
	osOthShift   = 0

	osUserR = osRead << osUserShift
	osUserW = osWrite << osUserShift
	osUserX = osEx << osUserShift

	osGroupR = osRead << osGroupShift
	osGroupW = osWrite << osGroupShift
	osGroupX = osEx << osGroupShift

	osOthR = osRead << osOthShift
	osOthW = osWrite << osOthShift
	osOthX = osEx << osOthShift
)

func setFileAsExecutable(f string) error {
	err := os.Chmod(f, os.FileMode(os.ModeDir|osUserR|osUserW|osUserX|osGroupR|osGroupX|osOthR|osOthX))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"file":  f,
			"error": err,
		}).Errorln("Error opening file")
	}
	return err
}

// ReadCSV read csv per line and send current one into channel
func ReadCSV(rc io.Reader) (ch chan []string) {
	ch = make(chan []string, 10)
	go func() {
		r := csv.NewReader(rc)
		if _, err := r.Read(); err != nil { //read header
			log.Fatal(err)
		}
		defer close(ch)
		for {
			rec, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatal(err)

			}
			ch <- rec
		}
	}()
	return
}
