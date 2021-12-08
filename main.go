package main

import (
	"errors"
	"go_photos/internal/service"
	"go_photos/internal/transport"
	"io"

	"github.com/sirupsen/logrus"
)

type mock struct{}

func (m mock) Detect(reader io.Reader) (string, error) {
	return "", errors.New("test")
}

func (m mock) Save(img []byte, cType string) (id int, err error) {
	return 0, errors.New("test")
}

func (m mock) Get(id int) (img []byte, cType string, err error) {
	return nil, "", errors.New("test")
}

func (m mock) Check(reader io.Reader) error {
	return errors.New("test")
}

func main() {
	log := logrus.New()
	log.ReportCaller = true
	log.SetLevel(logrus.DebugLevel)
	m := mock{}
	imageService := service.NewImageService(m, m, m, log)

	t := transport.NewTransport(":9090",
		imageService,
		imageService,
		log,
	)
	log.Info("transport initializing complete")
	if err := t.Run(); err != nil {
		log.WithError(err).Fatal("Could not start transport")
	}
}
