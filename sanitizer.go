package main

import (
	"errors"
	"html/template"
	"log"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/abadojack/whatlanggo"
	guuid "github.com/google/uuid"
)

var (
	hasCriticalError = false
	message          SanitizedMessage
)

// SanitizedMessage - contains the message after all sanitization
type SanitizedMessage struct {
	Headers HeadersMessage `json:"headers"`
	Content ContentMessage `json:"content"`
}

// HeadersMessage - contains the headers from rabbit
type HeadersMessage struct {
	CorrelationID string `json:"correlationId"`
	AppID         string `json:"appId"`
}

// ContentMessage - contains the content
type ContentMessage struct {
	Origem         string `json:"origem"`
	IDLive         string `json:"idLive"`
	DataLive       string `json:"dataLive"`
	DataPublicacao string `json:"dataPublicacao"`
	EmbedHTML      string `json:"embedHTML"`
	IDCategoria    string `json:"idCategoria"`
	IDCanal        string `json:"idCanal"`
	TituloCanal    string `json:"tituloCanal"`
	DescricaoLive  string `json:"descricaoLive"`
	ThumbDefault   string `json:"thumbDefault"`
	ThumbHigh      string `json:"thumbHigh"`
	ThumbMaxRes    string `json:"thumbMaxRes"`
	ThumbMedium    string `json:"thumbMedium"`
	ThumbStandard  string `json:"thumbStandard"`
	TituloLive     string `json:"tituloLive"`
	Likes          string `json:"likes"`
}

// Sanitizer - start the sanitizer proccess
func (m MessageResponse) Sanitizer() (SanitizedMessage, error) {

	hasCriticalError = false

	// define variable
	// var message SanitizedMessage

	// sanitize headers
	message.Headers.CorrelationID = m.sanitizeUUID()
	message.Headers.AppID = m.sanitizeAppID()

	// sanitize content
	message.Content.Origem = m.sanitizeSource()
	message.Content.IDLive = m.sanitizeIDLive()
	message.Content.DataLive = m.sanitizeDataLive()
	message.Content.DataPublicacao = m.sanitizeDataPublicacao()
	message.Content.EmbedHTML = m.sanitizeEmbedHTML()
	message.Content.IDCategoria = m.sanitizeIDCategoria()
	message.Content.IDCanal = m.sanitizeIDCanal()
	message.Content.TituloCanal = m.sanitizeTituloCanal()
	message.Content.DescricaoLive = m.sanitizeDescricaoLive()
	message.Content.ThumbDefault = m.sanitizeThumbDefault()
	message.Content.ThumbHigh = m.sanitizeThumbHigh()
	message.Content.ThumbMaxRes = m.sanitizeThumbMaxRes()
	message.Content.ThumbMedium = m.sanitizeThumbMedium()
	message.Content.ThumbStandard = m.sanitizeThumbStandard()
	message.Content.TituloLive = m.sanitizeTituloLive()
	message.Content.Likes = m.sanitizeLikes()
	// message = m.sanitizeLang()

	// check if there are any critical error during the process
	if hasCriticalError {
		return message, errors.New("An undefined error has happened during the sanitization, check log files")
	}

	return message, nil

}

func (m MessageResponse) sanitizeSource() string {
	return m.Source
}

func (m MessageResponse) sanitizeLikes() string {

	err := ReflectStructField(m.Videos.Items[0].Statistics, "LikeCount")

	// if does not existe
	if err != nil {
		return ""
	}

	if string(m.Videos.Items[0].Statistics.LikeCount) == "" {
		return ""
	}

	l := strconv.FormatUint(m.Videos.Items[0].Statistics.LikeCount, 10)

	return l
}

// Contains critical information
func (m MessageResponse) sanitizeTituloLive() string {

	err := ReflectStructField(m.Videos.Items[0].Snippet, "Title")

	// if does not existe
	if err != nil {
		return ""
	}

	if m.Videos.Items[0].Snippet.Title == "" {
		// create function to save log
		hasCriticalError = true
		log.Println(" [-] Error to sanitize Titulo live ")
		return ""
	}

	return m.Videos.Items[0].Snippet.Title
}

func (m MessageResponse) sanitizeThumbStandard() string {

	err := ReflectStructField(m.Videos.Items[0].Snippet.Thumbnails, "Standard")

	// if does not existe
	if err != nil {
		return ""
	}

	if m.Videos.Items[0].Snippet.Thumbnails.Standard.Url == "" {
		// create function to save log
		return ""
	}

	u := url.PathEscape(m.Videos.Items[0].Snippet.Thumbnails.Standard.Url)
	return u
}

func (m MessageResponse) sanitizeThumbMedium() string {

	err := ReflectStructField(m.Videos.Items[0].Snippet.Thumbnails, "Medium")

	// if does not existe
	if err != nil {
		return ""
	}

	if m.Videos.Items[0].Snippet.Thumbnails.Medium.Url == "" {
		// create function to save log
		return ""
	}
	u := url.PathEscape(m.Videos.Items[0].Snippet.Thumbnails.Medium.Url)
	return u
}

func (m MessageResponse) sanitizeThumbMaxRes() string {

	err := ReflectStructField(m.Videos.Items[0].Snippet.Thumbnails, "Maxres")

	// if does not existe
	if err != nil {

		return ""
	}

	if m.Videos.Items[0].Snippet.Thumbnails.Maxres.Url == "" {
		// create function to save log
		return ""
	}

	u := url.PathEscape(m.Videos.Items[0].Snippet.Thumbnails.Maxres.Url)
	return u
}

func (m MessageResponse) sanitizeThumbHigh() string {

	err := ReflectStructField(m.Videos.Items[0].Snippet.Thumbnails, "High")

	// if does not existe
	if err != nil {
		return ""
	}

	if m.Videos.Items[0].Snippet.Thumbnails.High.Url == "" {
		// create function to save log
		return ""
	}

	u := url.PathEscape(m.Videos.Items[0].Snippet.Thumbnails.High.Url)
	return u
}

func (m MessageResponse) sanitizeThumbDefault() string {

	err := ReflectStructField(m.Videos.Items[0].Snippet.Thumbnails, "Default")

	// if does not existe
	if err != nil {
		return ""
	}

	if m.Videos.Items[0].Snippet.Thumbnails.Default.Url == "" {
		// create function to save log
		return ""
	}

	u := url.PathEscape(m.Videos.Items[0].Snippet.Thumbnails.Default.Url)
	return u
}

func (m MessageResponse) sanitizeDescricaoLive() string {

	err := ReflectStructField(m.Videos.Items[0].Snippet, "Description")

	// if does not existe
	if err != nil {
		return ""
	}

	if m.Videos.Items[0].Snippet.Description == "" {
		// create function to save log
		return ""
	}

	return template.HTMLEscapeString(m.Videos.Items[0].Snippet.Description)
}

func (m MessageResponse) sanitizeTituloCanal() string {

	err := ReflectStructField(m.Videos.Items[0].Snippet, "ChannelTitle")

	// if does not existe
	if err != nil {
		return ""
	}

	if m.Videos.Items[0].Snippet.ChannelTitle == "" {
		// create function to save log
		return ""
	}

	return m.Videos.Items[0].Snippet.ChannelTitle
}

func (m MessageResponse) sanitizeIDCanal() string {

	err := ReflectStructField(m.Videos.Items[0].Snippet, "ChannelId")

	// if does not existe
	if err != nil {
		return ""
	}

	if m.Videos.Items[0].Snippet.ChannelId == "" {
		// create function to save log
		return ""
	}

	return m.Videos.Items[0].Snippet.ChannelId
}

func (m MessageResponse) sanitizeIDCategoria() string {

	err := ReflectStructField(m.Videos.Items[0].Snippet, "CategoryId")

	// if does not existe
	if err != nil {
		return ""
	}

	if m.Videos.Items[0].Snippet.CategoryId == "" {
		// create function to save log
		return ""
	}

	// check category if it's religion
	isReligion := isRelegionCategory(m.Videos.Items[0].Snippet.Title)

	// return ID 1000. it comes from the web database
	if isReligion {
		return "1000"
	}

	return m.Videos.Items[0].Snippet.CategoryId
}

func (m MessageResponse) sanitizeEmbedHTML() string {

	err := ReflectStructField(m.Videos.Items[0].Player, "EmbedHtml")

	// if does not existe
	if err != nil {
		return ""
	}

	if m.Videos.Items[0].Player.EmbedHtml == "" {
		// create function to save log
		return ""
	}

	return m.Videos.Items[0].Player.EmbedHtml
}

// Contains critical information
func (m MessageResponse) sanitizeDataPublicacao() string {

	err := ReflectStructField(m.Videos.Items[0].Snippet, "PublishedAt")

	// if does not existe
	if err != nil {
		hasCriticalError = true
		return ""
	}

	if m.Videos.Items[0].Snippet.PublishedAt == "" {
		// create function to save log
		hasCriticalError = true
		log.Println(" [-] Error to sanitize Data Publicação")
		return ""
	}

	return m.Videos.Items[0].Snippet.PublishedAt
}

// Contains critical information
func (m MessageResponse) sanitizeDataLive() string {

	err := ReflectStructField(m.Videos.Items[0].LiveStreamingDetails, "ScheduledStartTime")

	// if does not existe
	if err != nil {
		hasCriticalError = true
		return ""
	}

	if m.Videos.Items[0].LiveStreamingDetails.ScheduledStartTime == "" {
		// create function to save log
		hasCriticalError = true
		log.Println(" [-] Error to sanitize Live Data ")
		return ""
	}

	return m.Videos.Items[0].LiveStreamingDetails.ScheduledStartTime
}

// Contains critical information
func (m MessageResponse) sanitizeIDLive() string {

	err := ReflectStructField(m.Videos.Items[0], "Id")

	// if does not existe
	if err != nil {
		hasCriticalError = true
		return ""
	}

	if m.Videos.Items[0].Id == "" {
		// create function to save log
		hasCriticalError = true
		log.Println(" [-] Error to sanitize Live ID ")
		return ""
	}

	return m.Videos.Items[0].Id
}

func (m MessageResponse) sanitizeAppID() string {

	if m.Interal.AppID == "" {
		return "processor"
	}

	return m.Interal.AppID
}

func (m MessageResponse) sanitizeUUID() string {

	if m.Interal.CorrelationID == "" {
		return guuid.New().String()
	}

	return m.Interal.CorrelationID
}

func (m MessageResponse) sanitizeLang() string {

	err := ReflectStructField(m.Videos.Items[0].Snippet, "Title")

	// if does not existe
	if err != nil {
		return ""
	}

	if m.Videos.Items[0].Snippet.Title == "" {
		// create function to save log
		hasCriticalError = true
		log.Println(" [-] Error to sanitize Titulo live ")
		return ""
	}

	l, r := identifyLanguage(m.Videos.Items[0].Snippet.Title)

	// check if is realiable
	if r {
		return l
	}

	return ""
}

// ReflectStructField if an interface is either a struct or a pointer to a struct
// and has the defined member field, if error is nil, the given
// FieldName exists and is accessible with reflect.
func ReflectStructField(Iface interface{}, FieldName string) error {
	ValueIface := reflect.ValueOf(Iface)

	// Check if the passed interface is a pointer
	if ValueIface.Type().Kind() != reflect.Ptr {
		// Create a new type of Iface's Type, so we have a pointer to work with
		ValueIface = reflect.New(reflect.TypeOf(Iface))
	}

	// 'dereference' with Elem() and get the field by name
	Field := ValueIface.Elem().FieldByName(FieldName)

	if !Field.IsValid() {
		return errors.New("Element does not exist in the struct")
		// return fmt.Errorf("Interface `%s` does not have the field `%s`", ValueIface.Type(), FieldName)
	}

	// check if is empty
	if Field.Type().Kind() == reflect.String {
		if Field.String() == "" {
			return errors.New("Element does not exist in the struct")
		}
	}

	if Field.Type().Kind() == reflect.Ptr {
		if Field.IsNil() {
			// fmt.Printf("NAME: %v   --   Type: %v\n", Field, Field.Type().Kind())
			return errors.New("Element does not exist in the struct")
		}
	}

	return nil
}

func isRelegionCategory(title string) bool {

	tags := []string{"church", "gospel", "igreja", "missa", "santa", "evangelho", "paróquia", "paroquia", "adoracao", "adoracao", "orando", "louvor", "oração", "oracao", "culto", "evangeliza", "padre", "pastor", "cristão", "cristao", "benção", "bencao", "bíblia", "biblia"}

	_, i := Find(tags, title)

	return i
}

func identifyLanguage(text string) (lang string, isReliable bool) {

	info := whatlanggo.Detect(text)
	//	c := fmt.Sprintf("%f", info.Confidence)

	if info.Lang.String() == "" {
		return "", false
	}

	return info.Lang.String(), info.IsReliable()

}

// Find - return if an element contain one of the array items
func Find(slice []string, val string) (int, bool) {

	for i, item := range slice {

		if strings.Contains(strings.ToUpper(val), strings.ToUpper(item)) {
			return i, true
		}

	}

	return -1, false
}
