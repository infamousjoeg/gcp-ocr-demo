#!/bin/bash

# Create Cloud Storage bucket to upload images to
gsutil mb gs://gcp-ocr-demo
# Create Cloud Storage bucket to save text translations to
gsutil mb gs://gcp-ocr-demo-results
# Create a Pub/Sub topic to publish translation requests to
gcloud pubsub topics create gcp-ocr-demo
# Create a Pub/Sob topic to publish finished translation
# results to
gcloud pubsub topics create gcp-ocr-demo-results

echo "Finished preparing project."