# Workshop setup
In this lab you'll be setting up a Google Container Engine (GKE) cluster and a Cloud Identity and Access Management (Cloud IAM) Servicve Account so that you'll have enough resources and permissions to run Spinnaker and your sample application.

## Before you begin
Enable the Container Engine API by following [this link](https://console.cloud.google.com/flows/enableapi?apiid=container).

## Start a Cloud Shell instance

1. Open the cloud shell by navigating to [this link](https://console.cloud.google.com/?cloudshell=true).
1. The rest of the tutorial should be executed from inside the Cloud Shell.


## Create Google Container Engine Cluster

You will need a Container Engine cluster to deploy Spinnaker and the sample application. The cluster requires the storage-rw authentication scope for Spinnaker to store its pipeline data in Cloud Storage:
1. Enable the Container Engine API by following this link.
1. Then run the following comamnds in the cloud shell console:
```shell
gcloud config set compute/zone us-central1-f
gcloud container clusters create spinnaker-tutorial --cluster-version 1.6.10 --machine-type=n1-standard-2
```

## Configuring IAM
You will need to create a service account to delegate permissions to Spinnaker to be able to store its data in Google Cloud Storage. Spinnaker stores its pipeline data in Google Cloud Storage to ensure reliability and resiliency. If for some reason your Spinnaker deployment is torn down, you will be able to bring up an identical one in minutes with the same pipeline data.

### Create Service Account

1. First create the service account itself:
```shell
gcloud iam service-accounts create  spinnaker-storage-account  --display-name spinnaker-storage-account
```
1. Store the service account email address and your current project ID in an environment variables for use in later commands:
```shell
export SA_EMAIL=$(gcloud iam service-accounts list --filter="displayName:spinnaker-storage-account" --format='value(email)')
export PROJECT=$(gcloud info --format='value(config.project)')
```
1. Next bind the storage.admin role to your service account:
```shell
gcloud projects add-iam-policy-binding $PROJECT --role roles/storage.admin --member serviceAccount:$SA_EMAIL
```
1. Now that the service account has the appropriate permissions granted to it, you can download its key. The key will later be uploaded to Container Engine during the Spinnaker install:
```shell
gcloud iam service-accounts keys create spinnaker-sa.json --iam-account $SA_EMAIL
```

## What's Next
Now that you have a working cluster and a Cloud IAM service account, it's time to install Spinnaker!
