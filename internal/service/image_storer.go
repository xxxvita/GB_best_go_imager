package service

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

type ImageChecker interface {
	Check(io.Reader) error
}

type ImageStorage interface {
	Save(img []byte, cType string) (id int, err error)
	Get(id int) (img []byte, cType string, err error)
}

type ContentTypeDetector interface {
	Detect(io.Reader) (string, error)
}

type ImageService struct {
	checker  ImageChecker
	storage  ImageStorage
	detector ContentTypeDetector
	log      *logrus.Logger
}

func NewImageService(checker ImageChecker, storage ImageStorage, detector ContentTypeDetector, log *logrus.Logger) *ImageService {
	return &ImageService{checker: checker, storage: storage, detector: detector, log: log}
}

func (i *ImageService) Get(id int) (img io.Reader, contentType string, err error) {
	imgBytes, contentType, err := i.storage.Get(id)
	if err != nil {
		i.log.WithError(err).Error("could not get image")
		return nil, "", err
	}
	img = bytes.NewBuffer(imgBytes)
	return img, contentType, nil
}

func (i *ImageService) Store(r io.Reader) (id int, err error) {
	err = i.checker.Check(r)
	if err != nil {
		i.log.WithError(err).Warn("image invalid")
		return 0, err
	}
	cType, err := i.detector.Detect(r)
	if err != nil {
		i.log.WithError(err).Warn("detector return error")
		return 0, err
	}
	img, err := ioutil.ReadAll(r)
	if err != nil {
		i.log.WithError(err).Error("could not read image")
		return 0, err
	}
	id, err = i.storage.Save(img, cType)
	if err != nil {
		i.log.WithError(err).Error("could not save image")
		return 0, err
	}

	return id, nil
}
