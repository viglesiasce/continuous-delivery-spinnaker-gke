# Workshop setup

### Prerequisites
1. A Google Cloud Platform Account
1. After signing into your GCP Account, [Click here to enable the Google Compute Engine and Google Container Engine APIs](https://console.cloud.google.com/flows/enableapi?apiid=compute_component,container)
1. The Google Cloud SDK.  [Click here to install the Google Cloud SDK](https://cloud.google.com/sdk/downloads#interactive).

## Setup

1. Create a new Google Cloud Platform project: [https://console.cloud.google.com](https://console.cloud.google.com/)

1. In the shell, set your default compute zone:

  ```shell
  $ gcloud config set compute/zone us-east1-d
  ```

## Get the Code

1. Clone the lab repository in your shell, then `cd` into that dir:

  ```shell
  $ git clone https://github.com/askcarter/spinnaker-k8s-workshop.git
  Cloning into 'spinnaker-k8s-workshop'...
  ...

  $ cd spinnaker-k8s-workshop
  ```

## Cluster Setup

### Create Cluster

Spinnaker takes up a lot of resources.  Plus we need read write access to GCS.
```shell
$ gcloud container clusters create workshop --scopes=storage-rw --machine-type=n1-standard-2
```

### Create Service Account
We need RW access because we’re storing data in gcs instead of minio.
 
```shell
$ gcloud iam service-accounts create spinnaker-bootstrap-account --display-name spinnaker-bootstrap-account
$ SA_EMAIL=$(gcloud iam service-accounts list \
    --filter="displayName:spinnaker-bootstrap-account" --format="value(email.basename())")
$ PROJECT=$(gcloud info --format='value(config.project)')
```
 
```shell
$ gcloud projects add-iam-policy-binding $PROJECT --role roles/storage.admin --member serviceAccount:$SA_EMAIL
```
 
```shell
$ gcloud iam service-accounts keys create account.json --iam-account $SA_EMAIL
```
