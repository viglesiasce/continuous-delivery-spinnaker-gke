#!/bin/bash -xe
export PROJECT=$(gcloud info --format='value(config.project)')
export CONFIG_BUCKET=$PROJECT-spinnaker-config
export MANIFEST_BUCKET=$PROJECT-kubernetes-manifests

for SA_EMAIL in `gcloud iam service-accounts list --filter="displayName:spinnaker-account" --format='value(email)'`;do
  gcloud projects remove-iam-policy-binding $PROJECT --role roles/storage.admin --member serviceAccount:$SA_EMAIL || true
  gcloud beta pubsub subscriptions remove-iam-policy-binding gcr-triggers --role roles/pubsub.subscriber --member serviceAccount:$SA_EMAIL || true
  echo y | gcloud iam service-accounts delete $SA_EMAIL || true
done

echo y | gcloud container clusters delete spinnaker-tutorial --zone us-central1-f || true
echo y | gcloud source repos delete sample-app || true
gcloud beta pubsub subscriptions delete gcr-triggers || true
gsutil -m rm -r gs://$CONFIG_BUCKET || true
gsutil -m rm -r gs://$MANIFEST_BUCKET || true
rm -f spinnaker-sa.json
rm -f spinnaker-config.yaml
rm -f helm-v2.10.0-linux-amd64.tar.gz
rm -rf linux-amd64
rm -f helm
