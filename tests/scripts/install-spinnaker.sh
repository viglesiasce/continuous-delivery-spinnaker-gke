#!/bin/bash -xe
gcloud config set compute/zone us-central1-f
gcloud container clusters create spinnaker-tutorial --cluster-version 1.6.4 --machine-type=n1-standard-2
gcloud iam service-accounts create  spinnaker-storage-account  --display-name spinnaker-storage-account

# Needed to ensure SA gets role properly?

sleep 60
export SA_EMAIL=$(gcloud iam service-accounts list --filter="displayName:spinnaker-storage-account" --format='value(email)')
export PROJECT=$(gcloud info --format='value(config.project)')
gcloud projects add-iam-policy-binding $PROJECT --role roles/storage.admin --member serviceAccount:$SA_EMAIL
gcloud iam service-accounts keys create spinnaker-sa.json --iam-account $SA_EMAIL

wget https://storage.googleapis.com/kubernetes-helm/helm-v2.5.0-linux-amd64.tar.gz
tar zxfv helm-v2.5.0-linux-amd64.tar.gz
cp linux-amd64/helm .
./helm init
./helm update
# Give tiller a chance to start up
# TODO: Change this to polling
sleep 180
./helm version | grep 2.5.0

export PROJECT=$(gcloud info --format='value(config.project)')
export BUCKET=$PROJECT-spinnaker-config
gsutil mb -c regional -l us-central1 gs://$BUCKET

export SA_JSON=`cat spinnaker-sa.json`
export PROJECT=$(gcloud info --format='value(config.project)')
export BUCKET=$PROJECT-spinnaker-config
cat > spinnaker-config.yaml <<EOF
storageBucket: $BUCKET
gcs:
  enabled: true
  project: $PROJECT
  jsonKey: '$SA_JSON'

# Disable minio the default
minio:
  enabled: false

# Configure your Docker registries here
accounts:
- name: gcr
  address: https://gcr.io
  username: _json_key
  password: '$SA_JSON'
  email: 1234@5678.com
images:
  deck: quay.io/spinnaker/deck:v2.1095.0
EOF

./helm install -n cd stable/spinnaker -f spinnaker-config.yaml --version 0.3.1 --timeout 600

export DECK_POD=$(kubectl get pods --namespace default -l "component=deck" -o jsonpath="{.items[0].metadata.name}")
kubectl port-forward --namespace default $DECK_POD 8080:9000 &
sleep 10
curl localhost:8080

rm -rf .git
git config --global user.email "person@example.org"
git config --global user.name "person"
git init
git add .
git commit -m "Initial commit"
gcloud source repos create sample-app
git config credential.helper gcloud.sh
export PROJECT=$(gcloud info --format='value(config.project)')
git remote add origin https://source.developers.google.com/p/$PROJECT/r/sample-app

git push origin master
