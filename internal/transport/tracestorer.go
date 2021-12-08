package transport

// DO NOT EDIT!
// This code is generated with http://github.com/hexdigest/gowrap tool
// using https://raw.githubusercontent.com/hexdigest/gowrap/master/templates/logrus template

//go:generate gowrap gen -p go_photos/internal/transport -i ImageStorer -t https://raw.githubusercontent.com/hexdigest/gowrap/master/templates/logrus -o tracestorer.go

import (
	"io"

	"github.com/sirupsen/logrus"
)

// ImageStorerWithLogrus implements ImageStorer that is instrumented with logrus logger
type ImageStorerWithLogrus struct {
	_log  *logrus.Entry
	_base ImageStorer
}

// NewImageStorerWithLogrus instruments an implementation of the ImageStorer with simple logging
func NewImageStorerWithLogrus(base ImageStorer, log *logrus.Entry) ImageStorerWithLogrus {
	return ImageStorerWithLogrus{
		_base: base,
		_log:  log,
	}
}

// Store implements ImageStorer
func (_d ImageStorerWithLogrus) Store(r1 io.Reader) (id int, err error) {
	_d._log.WithFields(logrus.Fields(map[string]interface{}{
		"r1": r1})).Debug("ImageStorerWithLogrus: calling Store")
	defer func() {
		if err != nil {
			_d._log.WithFields(logrus.Fields(map[string]interface{}{
				"id":  id,
				"err": err})).Error("ImageStorerWithLogrus: method Store returned an error")
		} else {
			_d._log.WithFields(logrus.Fields(map[string]interface{}{
				"id":  id,
				"err": err})).Debug("ImageStorerWithLogrus: method Store finished")
		}
	}()
	return _d._base.Store(r1)
}
