package http

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/qwantix/qxeye/camera"
	"github.com/qwantix/qxeye/config"
	"github.com/qwantix/qxeye/service/file"
	"github.com/qwantix/qxeye/util"
)

type HttpService struct {
	config config.ServiceConfig
	url    string
	fs     *file.FileService
}

func New(cfg config.ServiceConfig) *HttpService {
	hs := new(HttpService)
	hs.url = cfg.Params.String("url")
	if hs.url == "" {
		util.Error("http: Missing url in http service")
		return nil
	}
	hs.fs = file.New(config.ServiceConfig{Service: "file", Params: config.Hmap{"dir": "/tmp", "traces": cfg.Params.String("traces")}})
	return hs
}

func (hs *HttpService) Push(ev *camera.CameraEvent, t *config.TriggerConfig) {
	filename := hs.fs.CaptureToFile(ev, t)
	util.Log("http: Push to ", hs.url)
	defer hs.fs.RemoveFile(filename)

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Add image
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()
	fw, err := w.CreateFormFile("capture", filename)
	if util.CheckErr(err, "http: CreateFormField capture") {
		return
	}
	if _, err = io.Copy(fw, f); util.CheckErr(err, "http: io.Copy") {
		return
	}

	writeField := func(name, value string) {
		if fw, err = w.CreateFormField(name); util.CheckErr(err, "http: CreateFormField ", name) {
			return
		}
		if _, err = fw.Write([]byte(value)); util.CheckErr(err, "http: write ", name) {
			return
		}
	}

	writeField("camera", ev.CameraName)
	// TODO write matches

	w.Close()

	// Send
	req, err := http.NewRequest("POST", hs.url, &b)
	if util.CheckErr(err, "http: NewRequest") {
		return
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	client := &http.Client{}
	res, err := client.Do(req)
	if util.CheckErr(err, "http: doing request") {
		return
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
		util.Error("bad status: ", res.Status)
	}
	return
}
