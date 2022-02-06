package proxy

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var imageTypes = map[string]interface{}{
	"image/jpeg": nil,
	"image/png":  nil,
}

// 5 MB
const (
	MB           = 1 << (10 * 2)
	maxImageSize = 5 * MB
)

type ImageReceiver struct {
	client http.Client
}

func NewImageReceiver() ImageReceiver {
	tr := &http.Transport{
		MaxIdleConns:    100,
		IdleConnTimeout: 90 * time.Second,
	}

	receiver := ImageReceiver{
		client: http.Client{
			Transport: tr,
		},
	}

	return receiver
}

func (ir *ImageReceiver) Receive(imageUrl string, headers http.Header) (*http.Response, error) {
	// building request
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	imageRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, imageUrl, nil)
	if err != nil {
		return nil, &StatusError{
			Err:  err,
			Code: http.StatusServiceUnavailable,
		}
	}
	var builder strings.Builder

	for header, values := range headers {
		builder.Reset()

		for _, value := range values {
			_, err := builder.WriteString(value + " ")
			if err != nil {
				return nil, &StatusError{
					Err:  err,
					Code: http.StatusServiceUnavailable,
				}
			}
		}

		imageRequest.Header.Add(header, strings.TrimRight(builder.String(), " "))
	}

	// doing request
	response, err := ir.client.Do(imageRequest)
	if err != nil {
		return nil, &StatusError{
			Err:  err,
			Code: http.StatusServiceUnavailable,
		}
	}

	err = checkResponse(response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func checkResponse(response *http.Response) error {
	// checking status
	if response.StatusCode != http.StatusOK {
		return &StatusError{
			Err:  ErrBadResponse,
			Code: response.StatusCode,
		}
	}

	// checking Content-Length
	if headerContentLength, ok := response.Header["Content-Length"]; ok {
		contentLength, err := strconv.Atoi(headerContentLength[0])
		if err != nil {
			return &StatusError{
				Err:  err,
				Code: http.StatusServiceUnavailable,
			}
		}

		if contentLength > maxImageSize {
			return &StatusError{
				Err:  ErrMaxSizeExceed,
				Code: http.StatusUnprocessableEntity,
			}
		}
	}

	// checking type
	if contentType, ok := response.Header["Content-Type"]; ok {
		responseType := contentType[0]

		if _, ok := imageTypes[responseType]; !ok {
			return &StatusError{
				Err:  ErrUnsupportedFormat,
				Code: http.StatusUnprocessableEntity,
			}
		}
	}

	return nil
}
