package ocr

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"cloud.google.com/go/translate"
	vision "cloud.google.com/go/vision/apiv1"
	"golang.org/x/text/language"
)

type ocrMessage struct {
	Text     string       `json:"text"`
	FileName string       `json:"fileName"`
	Lang     language.Tag `json:"lang"`
	SrcLang  language.Tag `json:"srcLang"`
}

// GCSEvent is the payload of a GCS event.
type GCSEvent struct {
	Bucket         string    `json:"bucket"`
	Name           string    `json:"name"`
	Metageneration string    `json:"metageneration"`
	ResourceState  string    `json:"resourceState"`
	TimeCreated    time.Time `json:"timeCreated"`
	Updated        time.Time `json:"updated"`
}

// PubSubMessage is the payload of a Pub/Sub event.
// See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Data []byte `json:"data"`
}

var (
	visionClient    *vision.ImageAnnotatorClient
	translateClient *translate.Client
	pubsubClient    *pubsub.Client
	storageClient   *storage.Client

	projectID      string
	resultBucket   string
	resultTopic    string
	toLang         []string
	translateTopic string
)

func setup(ctx context.Context) error {
	projectID = os.Getenv("GCP_PROJECT")
	resultBucket = os.Getenv("RESULT_BUCKET")
	resultTopic = os.Getenv("RESULT_TOPIC")
	toLang = strings.Split(os.Getenv("TO_LANG"), ",")
	translateTopic = os.Getenv("TRANSLATE_TOPIC")

	var err error // Prevent shadowing clients with :=.

	if visionClient == nil {
		visionClient, err = vision.NewImageAnnotatorClient(ctx)
		if err != nil {
			return fmt.Errorf("vision.NewImageAnnotatorClient: %v ", err)
		}
	}

	if translateClient == nil {
		translateClient, err = translate.NewClient(ctx)
		if err != nil {
			return fmt.Errorf("translate.NewClient: %v", err)
		}
	}

	if pubsubClient == nil {
		pubsubClient, err = pubsub.NewClient(ctx, projectID)
		if err != nil {
			return fmt.Errorf("translate.NewClient: %v", err)
		}
	}

	if storageClient == nil {
		storageClient, err = storage.NewClient(ctx)
		if err != nil {
			return fmt.Errorf("storage.NewClient: %v", err)
		}
	}
	return nil
}
