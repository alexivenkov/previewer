package proxy

import (
	"context"
	"errors"
	"io"
	"net/http"
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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
		return nil, &StatusError{
			Err:  err,
			Code: http.StatusServiceUnavailable,
		}
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

	// checking size
	body := response.Body
	defer body.Close()

	buf := make([]byte, maxImageSize)

	_, err := io.ReadFull(body, buf)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return &StatusError{
			Err:  err,
			Code: http.StatusServiceUnavailable,
		}
	}

	if err == nil {
		return &StatusError{
			Err:  ErrMaxSizeExceed,
			Code: http.StatusUnprocessableEntity,
		}
	}

	// checking type
	if _, ok := imageTypes[http.DetectContentType(buf)]; !ok {
		return &StatusError{
			Err:  ErrUnsupportedFormat,
			Code: http.StatusUnprocessableEntity,
		}
	}

	return nil
}
