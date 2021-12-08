package transport

// DO NOT EDIT!
// This code is generated with http://github.com/hexdigest/gowrap tool
// using https://raw.githubusercontent.com/hexdigest/gowrap/master/templates/logrus template

//go:generate gowrap gen -p go_photos/internal/transport -i ImageGetter -t https://raw.githubusercontent.com/hexdigest/gowrap/master/templates/logrus -o tracegetter.go

import (
	"io"

	"github.com/sirupsen/logrus"
)

// ImageGetterWithLogrus implements ImageGetter that is instrumented with logrus logger
type ImageGetterWithLogrus struct {
	_log  *logrus.Entry
	_base ImageGetter
}

// NewImageGetterWithLogrus instruments an implementation of the ImageGetter with simple logging
func NewImageGetterWithLogrus(base ImageGetter, log *logrus.Entry) ImageGetterWithLogrus {
	return ImageGetterWithLogrus{
		_base: base,
		_log:  log,
	}
}

// Get implements ImageGetter
func (_d ImageGetterWithLogrus) Get(id int) (img io.Reader, contentType string, err error) {
	_d._log.WithFields(logrus.Fields(map[string]interface{}{
		"id": id})).Debug("ImageGetterWithLogrus: calling Get")
	defer func() {
		if err != nil {
			_d._log.WithFields(logrus.Fields(map[string]interface{}{
				"img":         img,
				"contentType": contentType,
				"err":         err})).Error("ImageGetterWithLogrus: method Get returned an error")
		} else {
			_d._log.WithFields(logrus.Fields(map[string]interface{}{
				"img":         img,
				"contentType": contentType,
				"err":         err})).Debug("ImageGetterWithLogrus: method Get finished")
		}
	}()
	return _d._base.Get(id)
}
