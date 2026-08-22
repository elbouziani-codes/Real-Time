package domain

import (
	"strings"

	"uuid"
)

type WsRequest struct {
	// message or typing...
	RequestType string `json:"request_type"`
	Mod         string `json:"mod"`
	ID          string `json:"id"`
	Content     string `json:"content"`
	Destination string `json:"destination"`
}

type WsParsedRequest struct {
	// message or typing...
	RequestType string `json:"request_type"`
	Mod         string `json:"mod"`
	//ID          uuid.UUID `json:"id"`
	Content     string    `json:"content"`
	Destination uuid.UUID `json:"destination"`
}

func (wsR *WsRequest) ValidRequest() (*WsParsedRequest, error) {
	if wsR.RequestType != "message" && wsR.RequestType != "typing" {
		return nil, Error{Message: "Error in Type Request", Code: BadFormatCode}
	}
	if wsR.RequestType != "typing" {
		wsR.Content = strings.TrimSpace(wsR.Content)
		if len(wsR.Content) > 2048 || len(wsR.Content) == 0 {
			return nil, Error{Message: "length message must be between 1 and 2048  chars", Code: BadFormatCode}
		}
		if wsR.Mod != "Create" && wsR.Mod != "Edit" && wsR.Mod != "Delete" {
			return nil, Error{Message: "error in data Mod", Code: BadFormatCode}
		}
		if wsR.Mod != "Create" && wsR.ID == "" {
			return nil, Error{Message: "error in ID Message", Code: BadFormatCode}
		}
	}
	Destination, err := uuid.Parse(wsR.Destination)
	if err != nil {
		return nil, Error{Message: "error in ID Message", Code: BadFormatCode}
	}
	var parsed WsParsedRequest
	parsed.RequestType = wsR.RequestType
	parsed.Mod = wsR.Mod
	parsed.Content = wsR.Content
	parsed.Destination = Destination
	return &parsed, nil
}
