package services

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os"
	"sync"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/rs/zerolog/log"
)

type IBlobService interface {
	Upload(fileName string, data []byte) (string, error)
	Delete(fileName string) error
}

type BlobService struct {
	containerUrl azblob.ContainerURL
}

var (
	blobOnce    sync.Once
	blobService IBlobService
)

func GetBlobService() IBlobService {
	blobOnce.Do(func() {

		endpoint := os.Getenv("AZURITE_ENDPOINT")
		accountName := os.Getenv("AZURITE_NAME")
		accountKey := os.Getenv("AZURITE_KEY")
		containerName := os.Getenv("AZURITE_CONTAINER")

		// Parse the connection string
		credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
		if err != nil {
			log.Fatal().Err(err).Msg("Invalid credentials with error")
		}
		pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})

		// Create a URL to the container
		URL, _ := url.Parse(fmt.Sprintf("%s/%s", endpoint, containerName))

		// Create a container URL
		containerURL := azblob.NewContainerURL(*URL, pipeline)

		blobService = &BlobService{
			containerUrl: containerURL,
		}
	})

	return blobService
}

// Upload uploads a new blob to the container and returns the URL of the uploaded file.
func (s *BlobService) Upload(fileName string, data []byte) (string, error) {
	// Create a URL to the blob
	blobURL := s.containerUrl.NewBlockBlobURL(fileName)

	// Create a buffer from the data
	buffer := bytes.NewReader(data)

	// Upload the blob
	_, err := azblob.UploadStreamToBlockBlob(context.Background(), buffer, blobURL, azblob.UploadStreamToBlockBlobOptions{
		BufferSize: 4 * 1024 * 1024,
		MaxBuffers: 3,
		BlobHTTPHeaders: azblob.BlobHTTPHeaders{
			ContentType: "application/octet-stream",
		},
	})
	if err != nil {
		return "", err
	}

	// Create a URL to the blob on the CDN
	cdnURL := azblob.NewBlobURL(blobURL.URL(), azblob.NewPipeline(azblob.NewAnonymousCredential(), azblob.PipelineOptions{})).URL()
	if err != nil {
		return "", err
	}

	return cdnURL.String(), nil
}

// Delete deletes a blob.
func (s *BlobService) Delete(fileName string) error {
	// Create a URL to the blob
	blobURL := s.containerUrl.NewBlockBlobURL(fileName)

	// Delete the blob
	_, err := blobURL.Delete(context.Background(), azblob.DeleteSnapshotsOptionInclude, azblob.BlobAccessConditions{})
	return err
}
