package transport

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type ImageStorer interface {
	Store(io.Reader) (id int, err error)
}

type ImageGetter interface {
	Get(id int) (img io.Reader, contentType string, err error)
}

type Transport struct {
	storer ImageStorer
	getter ImageGetter
	router *mux.Router
	serv   *http.Server
	log    *logrus.Logger
}

func NewTransport(addr string, storer ImageStorer, getter ImageGetter, log *logrus.Logger) *Transport {
	t := &Transport{
		storer: NewImageStorerWithLogrus(storer, log.WithField("trace", true)),
		getter: NewImageGetterWithLogrus(getter, log.WithField("trace", true)),
	}
	r := mux.NewRouter()
	r.HandleFunc("/upload", t.Upload).Methods(http.MethodPost)
	r.HandleFunc("/download/{id}", t.Download).Methods(http.MethodGet)
	t.router = r
	srv := &http.Server{Handler: r, Addr: addr}
	t.serv = srv
	t.log = log
	return t
}

func (t *Transport) Run() error {
	//return errors.New("test")
	return t.serv.ListenAndServe()
}

type ErrorResponse struct {
	Error string
	Code  int
}

const (
	ErrCodeWrongInput   = 100
	ErrCodeFailGetImage = 101
)

const (
	maxMemory = 2 * 1024 * 1024 * 1024
)

type UploadResponse struct {
	Id int
}

func (t *Transport) Upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(maxMemory)
	if err != nil {
		t.log.WithError(err).Warn("could not pars multipart")
		t.errorResponse(http.StatusBadRequest, err, w)
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		t.log.WithError(err).Warn("could not get file from request")
		t.errorResponse(http.StatusBadRequest, err, w)
		return
	}
	id, err := t.storer.Store(file)
	if err != nil {
		t.log.WithError(err).Error("could not save file")
		t.errorResponse(http.StatusInternalServerError, err, w)
		return
	}
	jsResp, err := json.Marshal(UploadResponse{Id: id})
	if err != nil {
		t.log.WithError(err).Error("could not marhsall UploadResponse")
		t.errorResponse(http.StatusInternalServerError, err, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsResp)
	if err != nil {
		t.log.WithError(err).Error("could not write response")
	}
}

func (t *Transport) errorResponse(httpCode int, err error, w http.ResponseWriter) {
	erResp := ErrorResponse{
		Code:  httpCode,
		Error: err.Error(),
	}
	jsErr, mErr := json.Marshal(erResp)
	if mErr != nil {
		t.log.WithError(mErr).WithField("target_error", err).Error("fail to marshal error")
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	_, wErr := w.Write(jsErr)
	if wErr != nil {
		t.log.WithError(wErr).WithField("target_error", err).Error("fail send error response")
	}

}

func (t *Transport) Download(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		t.log.WithError(err).Warn("incorrect user input id isn't integer")
		t.errorResponse(ErrCodeWrongInput, err, w)
		return
	}
	body, cType, err := t.getter.Get(id)
	if err != nil {
		t.log.WithError(err).Warn("ImageGetter return error")
		t.errorResponse(ErrCodeFailGetImage, err, w)
		return
	}
	w.Header().Set("Content-Type", cType)
	_, err = io.Copy(w, body)
	if err != nil {
		t.log.WithError(err).Error("could not write response")
	}
}
