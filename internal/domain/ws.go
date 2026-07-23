package domain

import (
	"strings"

	"realTime/crypto"
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
	RequestType string      `json:"request_type"`
	Mod         string      `json:"mod"`
	ID          crypto.UUID `json:"id"`
	Content     string      `json:"content"`
	Destination crypto.UUID `json:"destination"`
}

func (wsR *WsRequest) ValidRequest() (*WsParsedRequest, error) {
	var parsed WsParsedRequest
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

	}K
	parsed.Destination.Value, err = crypto.ParseUUID(wsR.Destination)
	return &parsed, nil
}
