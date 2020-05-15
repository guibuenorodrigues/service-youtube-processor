package main

import (
	"errors"
	"log"
	"net/url"
	"strconv"

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

	// define variable
	// var message SanitizedMessage

	// sanitize headers
	message.Headers.CorrelationID = m.sanitizeUUID()
	message.Headers.AppID = m.sanitizeAppID()

	// sanitize content

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

	// check if there are any critical error during the process
	if hasCriticalError {
		return message, errors.New("An undefined error has happened during the sanitization, check log files")
	}

	return message, nil

}

func (m MessageResponse) sanitizeLikes() string {

	if string(m.Videos.Items[0].Statistics.LikeCount) == "" {
		// 	// create function to save log
		return ""
	}

	l := strconv.FormatUint(m.Videos.Items[0].Statistics.LikeCount, 10)

	return l
}

// Contains critical information
func (m MessageResponse) sanitizeTituloLive() string {

	if m.Videos.Items[0].Snippet.Title == "" {
		// create function to save log
		hasCriticalError = true
		log.Println(" [-] Error to sanitize Titulo live ")
		return ""
	}

	return m.Videos.Items[0].Snippet.Title
}

func (m MessageResponse) sanitizeThumbStandard() string {
	if m.Videos.Items[0].Snippet.Thumbnails.Standard.Url == "" {
		// create function to save log
		return ""
	}

	u := url.PathEscape(m.Videos.Items[0].Snippet.Thumbnails.Standard.Url)
	return u
}

func (m MessageResponse) sanitizeThumbMedium() string {

	if m.Videos.Items[0].Snippet.Thumbnails.Medium.Url == "" {
		// create function to save log
		return ""
	}
	u := url.PathEscape(m.Videos.Items[0].Snippet.Thumbnails.Medium.Url)
	return u
}

func (m MessageResponse) sanitizeThumbMaxRes() string {

	if m.Videos.Items[0].Snippet.Thumbnails.Maxres.Url == "" {
		// create function to save log
		return ""
	}
	u := url.PathEscape(m.Videos.Items[0].Snippet.Thumbnails.Maxres.Url)
	return u
}

func (m MessageResponse) sanitizeThumbHigh() string {

	if m.Videos.Items[0].Snippet.Thumbnails.High.Url == "" {
		// create function to save log
		return ""
	}

	u := url.PathEscape(m.Videos.Items[0].Snippet.Thumbnails.High.Url)
	return u
}

func (m MessageResponse) sanitizeThumbDefault() string {

	if m.Videos.Items[0].Snippet.Thumbnails.Default.Url == "" {
		// create function to save log
		return ""
	}

	u := url.PathEscape(m.Videos.Items[0].Snippet.Thumbnails.Default.Url)
	return u
}

func (m MessageResponse) sanitizeDescricaoLive() string {

	if m.Videos.Items[0].Snippet.Description == "" {
		// create function to save log
		return ""
	}

	return m.Videos.Items[0].Snippet.Description
}

func (m MessageResponse) sanitizeTituloCanal() string {

	if m.Videos.Items[0].Snippet.ChannelTitle == "" {
		// create function to save log
		return ""
	}

	return m.Videos.Items[0].Snippet.ChannelTitle
}

func (m MessageResponse) sanitizeIDCanal() string {

	if m.Videos.Items[0].Snippet.ChannelId == "" {
		// create function to save log
		return ""
	}

	return m.Videos.Items[0].Snippet.ChannelId
}

func (m MessageResponse) sanitizeIDCategoria() string {

	if m.Videos.Items[0].Snippet.CategoryId == "" {
		// create function to save log
		return ""
	}

	return m.Videos.Items[0].Snippet.CategoryId
}

func (m MessageResponse) sanitizeEmbedHTML() string {

	if m.Videos.Items[0].Player.EmbedHtml == "" {
		// create function to save log
		return ""
	}

	return m.Videos.Items[0].Player.EmbedHtml
}

// Contains critical information
func (m MessageResponse) sanitizeDataPublicacao() string {

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
