package twitter

import (
	"bytes"
	"html/template"

	"github.com/qwantix/qxeye/camera"
	"github.com/qwantix/qxeye/config"
	"github.com/qwantix/qxeye/service/file"
	"github.com/qwantix/qxeye/util"
)

type TwitterService struct {
	config config.ServiceConfig
	tw     *TwitterClient
	fs     *file.FileService
}

func New(cfg config.ServiceConfig) *TwitterService {
	ts := new(TwitterService)
	var cred TwitterCredentials
	if cfg.Has("file") {
		cred.Load(cfg.String("file"))
	} else {
		cred.ConsumerKey = cfg.String("consumerKey")
		cred.ConsumerSecret = cfg.String("consumerSecret")
		cred.AccessToken = cfg.String("accessToken")
		cred.AccessTokenSecret = cfg.String("accessTokenSecret")
	}
	ts.tw = NewTwitterClient(cred)
	ts.fs = file.New(config.ServiceConfig{"dir": "/tmp"})
	return ts
}

func (ts *TwitterService) Push(ev *camera.CameraEvent, t *config.TriggerConfig) {
	username := t.Params.String("to")
	if username == "" {
		util.Error("Unable to notify, missing user")
		return
	}
	u, err := ts.tw.GetUser(username)
	if util.CheckErr(err, "Unable to fetch this twitter user: ", username) {
		return
	}
	msg := t.Params.String("message")
	tpl := template.New("t1")
	tpl, err = tpl.Parse(msg)
	util.CheckErr(err, "Invalid message")
	if msg == "" || err != nil {
		msg = "What's happen ?!"
		tpl, err = tpl.Parse(msg)
	}

	var data = make(map[string]interface{})
	data["cameraName"] = ev.CameraName
	data["date"] = ev.Date

	var buff bytes.Buffer
	tpl.Execute(&buff, data)
	msg = buff.String()
	util.Log(msg, u)
	filename := ts.fs.CaptureToFile(ev, t)
	util.Log(filename)
	ts.tw.DirectMessageWithMediaFilename(u.IdStr, msg, filename)
	ts.fs.RemoveFile(filename)
}
