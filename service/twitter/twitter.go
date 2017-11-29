package twitter

// https://github.com/kreshikhin/twitter-media-uploader/blob/master/main.go
//
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/mrjones/oauth"
	"github.com/qwantix/qxeye/util"
)

const StatusUpdateEndpoint string = "https://api.twitter.com/1.1/statuses/update.json"
const MediaUploadEndpoint string = "https://upload.twitter.com/1.1/media/upload.json"
const DirectMessageEndpoint string = "https://api.twitter.com/1.1/direct_messages/events/new.json"
const GetUserEndpoint string = "https://api.twitter.com/1.1/users/show.json"

type TwitterClient struct {
	client *http.Client
}

type TwitterCredentials struct {
	ConsumerKey       string `json:"consumerKey"`
	ConsumerSecret    string `json:"consumerSecret"`
	AccessToken       string `json:"accessToken"`
	AccessTokenSecret string `json:"accessTokenSecret"`
}

func (tc *TwitterCredentials) Load(filename string) error {
	util.Log("twitter: Load credential file ", filename)
	file, err := os.Open(filename)
	if util.CheckErr(err, "twitter: Unable to load credential") {
		return err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(tc)
	util.Log(tc)
	if util.CheckErr(err, "twitter: Unable to parse credential file") {
		return err
	}
	return nil
}

type GetUserResponse struct {
	Id         uint64 `json:"id"`
	IdStr      string `json:"id_str"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}

type MediaInitResponse struct {
	MediaId          uint64 `json:"media_id"`
	MediaIdString    string `json:"media_id_string"`
	ExpiresAfterSecs uint64 `json:"expires_after_secs"`
}

type DMEvent struct {
	Event DMEventRequest `json:"event"`
}
type DMEventRequest struct {
	Type          string               `json:"type"`
	MessageCreate DMEventMessageCreate `json:"message_create,omitempty"`
}
type DMEventMessageCreate struct {
	Target  DMEventMessageCreateTarget `json:"target"`
	Message DMEventMessageCreateData   `json:"message_data"`
}
type DMEventMessageCreateTarget struct {
	RecipientId string `json:"recipient_id"`
}
type DMEventMessageCreateData struct {
	Text       string        `json:"text"`
	Attachment DMAttachement `json:"attachment,omitempty"`
}
type DMAttachement struct {
	Type  string  `json:"type"`
	Media DMMedia `json:"media"`
}
type DMMedia struct {
	Id uint64 `json:"id"`
}

var getUserCache map[string]*GetUserResponse = make(map[string]*GetUserResponse)

func NewTwitterClient(creds TwitterCredentials) *TwitterClient {
	t := new(TwitterClient)

	c := oauth.NewConsumer(
		creds.ConsumerKey,
		creds.ConsumerSecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
			AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
			AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
		})

	var access oauth.AccessToken
	access.Token = creds.AccessToken
	access.Secret = creds.AccessTokenSecret
	client, err := c.MakeHttpClient(&access)
	if err != nil {
		return nil
	}
	t.client = client
	return t
}

func (t *TwitterClient) GetUser(screenName string) (*GetUserResponse, error) {
	if getUserCache[screenName] != nil {
		return getUserCache[screenName], nil
	}
	req, err := http.NewRequest("GET", GetUserEndpoint, nil)

	q := req.URL.Query()
	q.Add("screen_name", screenName)
	req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL)
	res, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	var response GetUserResponse
	err = json.Unmarshal(body, &response)
	getUserCache[screenName] = &response
	return getUserCache[screenName], nil
}

func (t *TwitterClient) DirectMessageWithMedia(userId, message string, media []byte) error {
	fmt.Println("bytes", len(media))

	mediaInitResponse, err := t.MediaInit(media, "image/jpeg", "dm_image")
	if err != nil {
		fmt.Println("Can't init media", err)
	}

	fmt.Println(mediaInitResponse)

	mediaId := mediaInitResponse.MediaId

	if t.MediaAppend(mediaId, media) != nil {
		fmt.Println("Cant't append media")
	}

	if t.MediaFinilize(mediaId) != nil {
		fmt.Println("Cant't fin media")
	}
	if err != nil {
		return err
	}
	dm := DMEvent{}
	dm.Event.Type = "message_create"
	dm.Event.MessageCreate.Target.RecipientId = userId
	dm.Event.MessageCreate.Message.Text = message
	dm.Event.MessageCreate.Message.Attachment.Type = "media"
	dm.Event.MessageCreate.Message.Attachment.Media.Id = mediaId
	b, err := json.Marshal(&dm)
	fmt.Println("Send: ", string(b))
	req, err := http.NewRequest("POST", DirectMessageEndpoint, bytes.NewBuffer(b))
	req.Header.Add("Content-Type", "application/json")
	res, err := t.client.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(res.Body)
	//
	fmt.Println(string(body))
	fmt.Println(err)
	return nil
}

func (t *TwitterClient) DirectMessageWithMediaFilename(userId, message string, filename string) error {
	media, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return t.DirectMessageWithMedia(userId, message, media)
}

func (t *TwitterClient) MediaUpload(media []byte) (*MediaInitResponse, error) {
	var bodyInput bytes.Buffer
	w := multipart.NewWriter(&bodyInput)
	fw, err := w.CreateFormFile("media", "capture.jpg")
	n, err := fw.Write(media)
	fmt.Println("Write: ", n)
	w.Close()
	req, err := http.NewRequest("POST", MediaUploadEndpoint, &bodyInput)
	req.Header.Add("Content-Type", w.FormDataContentType())

	res, err := t.client.Do(req)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	fmt.Println("response", string(body))

	var mediaInitResponse MediaInitResponse
	err = json.Unmarshal(body, &mediaInitResponse)

	if err != nil {
		return nil, err
	}

	fmt.Println("Initialized media: ", mediaInitResponse)

	return &mediaInitResponse, nil
}

func (t *TwitterClient) MediaInit(media []byte, mimeType, category string) (*MediaInitResponse, error) {
	form := url.Values{}
	form.Add("command", "INIT")
	form.Add("media_type", mimeType)
	form.Add("media_category", category)
	form.Add("total_bytes", fmt.Sprint(len(media)))

	fmt.Println(form.Encode())

	req, err := http.NewRequest("POST", MediaUploadEndpoint, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := t.client.Do(req)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	fmt.Println("response", string(body))

	var mediaInitResponse MediaInitResponse
	err = json.Unmarshal(body, &mediaInitResponse)

	if err != nil {
		return nil, err
	}

	fmt.Println("Initialized media: ", mediaInitResponse)

	return &mediaInitResponse, nil
}

func (t *TwitterClient) MediaAppend(mediaId uint64, media []byte) error {
	step := 500 * 1024
	for s := 0; s*step < len(media); s++ {
		var body bytes.Buffer
		rangeBegining := s * step
		rangeEnd := (s + 1) * step
		if rangeEnd > len(media) {
			rangeEnd = len(media)
		}

		fmt.Println("try to append ", rangeBegining, "-", rangeEnd)

		w := multipart.NewWriter(&body)

		w.WriteField("command", "APPEND")
		w.WriteField("media_id", fmt.Sprint(mediaId))
		w.WriteField("segment_index", fmt.Sprint(s))

		fw, err := w.CreateFormFile("media", "capture.jpg")

		fmt.Println(body.String())

		n, err := fw.Write(media[rangeBegining:rangeEnd])

		fmt.Println("len ", n)

		w.Close()

		req, err := http.NewRequest("POST", MediaUploadEndpoint, &body)

		req.Header.Add("Content-Type", w.FormDataContentType())

		res, err := t.client.Do(req)
		if err != nil {
			return err
		}

		resBody, err := ioutil.ReadAll(res.Body)
		fmt.Println("append response ", string(resBody))
	}

	return nil
}

func (t *TwitterClient) MediaFinilize(mediaId uint64) error {
	form := url.Values{}
	form.Add("command", "FINALIZE")
	form.Add("media_id", fmt.Sprint(mediaId))

	req, err := http.NewRequest("POST", MediaUploadEndpoint, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := t.client.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(res.Body)
	fmt.Println("final response ", string(body))

	return nil
}
