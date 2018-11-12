#!/bin/bash -xe
gcloud config set compute/zone us-central1-f
gcloud container clusters create spinnaker-tutorial \
        --machine-type=n1-standard-2
SA_NAME=spinnaker-account-$(date +%s)
gcloud iam service-accounts create $SA_NAME \
        --display-name $SA_NAME

export SA_EMAIL=$(gcloud iam service-accounts list \
    --filter="displayName:$SA_NAME" \
    --format='value(email)')
export PROJECT=$(gcloud info --format='value(config.project)')
gcloud projects add-iam-policy-binding \
        $PROJECT --role roles/storage.admin --member serviceAccount:$SA_EMAIL
gcloud iam service-accounts keys create spinnaker-sa.json --iam-account $SA_EMAIL

# Create PubSub trigger
gcloud beta pubsub topics create projects/$PROJECT/topics/gcr || true
gcloud beta pubsub subscriptions create gcr-triggers \
    --topic projects/${PROJECT}/topics/gcr
export SA_EMAIL=$(gcloud iam service-accounts list \
    --filter="displayName:$SA_NAME" \
    --format='value(email)')
gcloud beta pubsub subscriptions add-iam-policy-binding gcr-triggers \
        --role roles/pubsub.subscriber --member serviceAccount:$SA_EMAIL

# Deploy Spinnaker using Helm
wget https://storage.googleapis.com/kubernetes-helm/helm-v2.10.0-linux-amd64.tar.gz
tar zxfv helm-v2.10.0-linux-amd64.tar.gz
cp linux-amd64/helm .
kubectl create clusterrolebinding user-admin-binding --clusterrole=cluster-admin --user=$(gcloud config get-value account)
kubectl create serviceaccount tiller --namespace kube-system
kubectl create clusterrolebinding tiller-admin-binding --clusterrole=cluster-admin --serviceaccount=kube-system:tiller

kubectl create clusterrolebinding --clusterrole=cluster-admin --serviceaccount=default:default spinnaker-admin

./helm init --service-account=tiller
./helm update

until ./helm version; do sleep 10;done

export PROJECT=$(gcloud info \
    --format='value(config.project)')
export BUCKET=$PROJECT-spinnaker-config
gsutil mb -c regional -l us-central1 gs://$BUCKET

export SA_JSON=$(cat spinnaker-sa.json)
export PROJECT=$(gcloud info --format='value(config.project)')
export BUCKET=$PROJECT-spinnaker-config
cat > spinnaker-config.yaml <<EOF
gcs:
  enabled: true
  bucket: $BUCKET
  project: $PROJECT
  jsonKey: '$SA_JSON'

dockerRegistries:
- name: gcr
  address: https://gcr.io
  username: _json_key
  password: '$SA_JSON'
  email: 1234@5678.com


# Disable minio as the default storage backend
minio:
  enabled: false


# Configure Spinnaker to enable GCP services
halyard:
  spinnakerVersion: 1.10.2
  image:
    tag: 1.12.0
  additionalScripts:
    create: true
    data:
      enable_gcs_artifacts.sh: |-
        \$HAL_COMMAND config artifact gcs account add gcs-$PROJECT --json-path /opt/gcs/key.json
        \$HAL_COMMAND config artifact gcs enable
      enable_pubsub_triggers.sh: |-
        \$HAL_COMMAND config pubsub google enable
        \$HAL_COMMAND config pubsub google subscription add gcr-triggers \
          --subscription-name gcr-triggers \
          --json-path /opt/gcs/key.json \
          --project $PROJECT \
          --message-format GCR
EOF

./helm install -n cd stable/spinnaker -f spinnaker-config.yaml --timeout 600 \
    --version 1.1.6 --wait

export DECK_POD=$(kubectl get pods --namespace default -l "cluster=spin-deck" \
    -o jsonpath="{.items[0].metadata.name}")

# Wiat for pods to settle before doing port-forward
until kubectl get pods $DECK_POD | grep Running; do sleep 10;done
kubectl port-forward --namespace default $DECK_POD 8080:9000 &


until curl localhost:8080; do sleep 10;done

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

# Setup K8s manifests
export PROJECT=$(gcloud info --format='value(config.project)')
gsutil mb -l us-central1 gs://$PROJECT-kubernetes-manifests
gsutil versioning set on gs://$PROJECT-kubernetes-manifests
sed -i s/PROJECT/$PROJECT/g k8s/deployments/*

# Install spin CLI
curl -LO https://storage.googleapis.com/spinnaker-artifacts/spin/$(curl -s https://storage.googleapis.com/spinnaker-artifacts/spin/latest)/linux/amd64/spin
chmod +x spin

# Create application
until timeout 1m ./spin application save --application-name sample \
                        --owner-email example@example.com \
                        --cloud-providers kubernetes \
                        --gate-endpoint http://localhost:8080/gate; do sleep 10;done

# Create pipeline
export PROJECT=$(gcloud info --format='value(config.project)')
sed s/PROJECT/$PROJECT/g spinnaker/pipeline-deploy.json > pipeline.json
./spin pipeline save --gate-endpoint http://localhost:8080/gate -f pipeline.json

# Push a tag
git tag v1.0.0
git push --tags

kill %1
