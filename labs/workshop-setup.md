# Workshop setup

### Prerequisites
1. A Google Cloud Platform Account
1. After signing into your GCP Account, [Click here to enable the Google Compute Engine and Google Container Engine APIs](https://console.cloud.google.com/flows/enableapi?apiid=compute_component,container)

## Start a Cloud Shell instance

1. Open the cloud shell by navigating to [this link](https://console.cloud.google.com/?cloudshell=true).
1. The rest of the tutorial should be executed from inside the Cloud Shell.


## Setup

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

You will need a Container Engine cluster to deploy Spinnaker and the sample application. The cluster requires the storage-rw authentication scope for Spinnaker to store its pipeline data in Cloud Storage:
1. Enable the Container Engine API by following this link.
1. Then run the following comamnds in the cloud shell console:
```shell
gcloud config set compute/zone us-central1-f
gcloud container clusters create spinnaker-tutorial --cluster-version 1.6.7 --machine-type=n1-standard-2
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

## Installing Spinnaker
You will use Helm to install Spinnaker from the Charts repository. Helm is a package manager that makes it easy to configure and deploy Kubernetes applications.

### Installing Helm
1. Download and install the helm binary
```shell
wget https://storage.googleapis.com/kubernetes-helm/helm-v2.5.0-linux-amd64.tar.gz
```

1. Unzip the file to your local system:
```shell
tar zxfv helm-v2.5.0-linux-amd64.tar.gz
cp linux-amd64/helm .
```

1. Initialize Helm. This ensures that the server side of Helm (Tiller) is properly installed in your cluster.
```shell
./helm init
./helm update
```

1. Ensure Helm is properly installed by running the following command. You should see versions appear for both the server and the client of ```v2.5.0```:
```shell
./helm version
Client: &version.Version{SemVer:"v2.5.0", GitCommit:"012cb0ac1a1b2f888144ef5a67b8dab6c2d45be6", GitTreeState:"clean"}Server: &version.Version{SemVer:"v2.5.0", GitCommit:"012cb0ac1a1b2f888144ef5a67b8dab6c2d45be6", GitTreeState:"clean"}
```
### Configure Spinnaker
1. Create a bucket for spinnaker to store its pipeline configuration:
```shell
export PROJECT=$(gcloud info --format='value(config.project)')
export BUCKET=$PROJECT-spinnaker-config
gsutil mb -c regional -l us-central1 gs://$BUCKET
```
1. Create your configuration file by pasting the following commands into your Cloud Shell:
```
export SA_JSON=`cat spinnaker-sa.json`
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
EOF
```

## Deploying the Spinnaker Chart
1. Use the Helm CLI to deploy the chart with your configuration set. 
> **Note:** This command typically takes 5-10 minutes to complete.

```shell
./helm install -n cd stable/spinnaker -f spinnaker-config.yaml --timeout 600 --version 0.3.1
```

1. Once that command completes, run the following command to setup port forwarding to the Spinnaker U from the Cloud Shell
```shell
export DECK_POD=$(kubectl get pods --namespace default -l "component=deck" -o jsonpath="{.items[0].metadata.name}")
kubectl port-forward --namespace default $DECK_POD 8080:9000 >> /dev/null &
```
To get to the Spinnaker user interface, click on the Web Preview button in cloud shell, then click “Preview on port 8080”
![](../docs/img/PLACEHOLDER.png)

You should now see the following screen:
![](../docs/img/PLACEHOLDER.png)

