package client

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

// FileRequester is an interface for getting a file from filestore.
type FileRequester interface {
	GetFile(uuid.UUID) (io.Reader, error)
}

// FileWriter is an interface for storing files.
type FileWriter interface {
	PutFile(uuid.UUID, io.Reader) error
	StoreFile(io.Reader) error
}

// Filestore is both a FileWriter and a FileRequester.
type FileStore interface {
	FileRequester
	FileWriter
}

type FileClient struct {
	baseURL string
	logger  *slog.Logger
}

func NewFileClient(baseURL string, logger *slog.Logger) *FileClient {
	return &FileClient{
		baseURL: baseURL,
		logger:  logger,
	}
}

func (fr *FileClient) GetFile(id uuid.UUID) (io.Reader, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s", fr.baseURL, id))
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (fr *FileClient) PutFile(id uuid.UUID, reader io.Reader) error {
	resp, err := http.Post(
		fmt.Sprintf("%s/%s", fr.baseURL, id),
		"application/octet-stream",
		reader,
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("POST request failed with status: %d", resp.StatusCode)
	}

	return nil
}

func (fr *FileClient) StoreFile(reader io.Reader) error {
	resp, err := http.Post(fmt.Sprintf("%s", fr.baseURL), "application/octet-stream", reader)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("POST request failed with status: %d", resp.StatusCode)
	}

	return nil
}
