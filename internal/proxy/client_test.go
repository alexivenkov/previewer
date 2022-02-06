package proxy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientOK(t *testing.T) {
	testCases := []struct {
		name        string
		code        int
		path        string
		contentType []string
	}{
		{
			name:        "test png",
			path:        "testdata/gopher.png",
			contentType: []string{"image/png"},
		},
		{
			name:        "test jpg",
			path:        "testdata/gopher.jpg",
			contentType: []string{"image/jpeg"},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			img, err := os.Open(tc.path)
			require.NoError(t, err)
			defer img.Close()

			stat, err := img.Stat()
			require.NoError(t, err)

			imageBuf := make([]byte, stat.Size())
			responseImageBuf := make([]byte, stat.Size())
			client := NewImageReceiver()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				n, err := img.Read(imageBuf)
				require.NoError(t, err)

				w.Header().Add("Content-Length", strconv.Itoa(n))

				_, err = w.Write(imageBuf)
				require.NoError(t, err)
			}))
			defer server.Close()

			res, err := client.Receive(server.URL, nil)
			require.NoError(t, err)

			responseImageBuf, err = io.ReadAll(res.Body)
			defer res.Body.Close()

			require.NoError(t, err)
			require.Equal(t, res.StatusCode, http.StatusOK)
			require.Equal(t, res.Header["Content-Type"], tc.contentType)
			require.Equal(t, imageBuf, responseImageBuf)
		})
	}
}

func TestWrongStatus(t *testing.T) {
	client := NewImageReceiver()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	res, err := client.Receive(server.URL, nil)
	require.Equal(t, (*http.Response)(nil), res)
	require.ErrorIs(t, ErrBadResponse, err.(*StatusError).Err)
	require.Equal(t, http.StatusServiceUnavailable, err.(*StatusError).Code)
}

func TestImageSizeExceed(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Length", strconv.Itoa(maxImageSize+1000))
	}))
	defer server.Close()

	client := NewImageReceiver()
	res, err := client.Receive(server.URL, nil)
	require.Equal(t, (*http.Response)(nil), res)
	require.ErrorIs(t, ErrMaxSizeExceed, err.(*StatusError).Err)
}

func TestResponseTypeNotImage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
	}))
	defer server.Close()

	client := NewImageReceiver()
	res, err := client.Receive(server.URL, nil)
	require.Equal(t, (*http.Response)(nil), res)
	require.ErrorIs(t, ErrUnsupportedFormat, err.(*StatusError).Err)
}
