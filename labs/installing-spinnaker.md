# Installing Spinnnaker

Spinnkaer has a lot of pieces and parts.  Below is a table listing everything.  You don't need to know *any* of this for the workshop, but it's here for completeness.

| Servivces | Port | Description |
| --- | --- | --- |
| Deck	| 9000 | Deck is a static AngularJS-based UI. |
| Clouddriver	| 7002 | Cloud Driver integrates with each cloud provider (AWS, GCP, Azure, etc.). It is responsible for all cloud provider-specific read and write operations. |
| Echo	| 8089 | Echo provides Spinnaker’s notification support, including integrations with Slack, Hipchat, SMS (via Twilio) and Email. |
| Front50	| 8080 | Front50 stores all application, pipeline and notification metadata. |
| Gate	| 8084 | Gate exposes APIs for all external consumers of Spinnaker (including deck). It is the front door to Spinnaker. |
| Igor	| 8088 | Igor facilitates the use of Jenkins in Spinnaker pipelines (a pipeline can be triggered by a Jenkins job or invoke a Jenkins job) |
| Orca	| 8083 | Orca handles pipeline and task orchestration (ie. starting a cloud driver operation and waiting until it completes). |
| Rosco	| 8087 | Rosco is a packer-based bakery. We believe in immutable infrastructure and rosco provides a means to take a Debian or Red Hat package and turn it into an Amazon Machine Image. Don’t worry, it also supports Google Compute Engine and Azure images. |
| Fiat	| 7003 | Fiat is the authorization server for the Spinnaker system.  It exposes a RESTful interface for querying the access permissions for a particular user. |

## Install Helm 
Download and install helm binary
From https://github.com/kubernetes/helm/blob/master/docs/quickstart.md
```shell
$ wget https://storage.googleapis.com/kubernetes-helm/helm-v2.4.2-linux-amd64.tar.gz
$ sudo tar -C /usr/local -xzf helm-v2.4.2-linux-amd64.tar.gz
```

Add Helm to PATH
```shell
$ export PATH=$PATH:/usr/local/linux-amd64
```
 
# Initialize local CLI
NOTE:  Does the directory matter for anything?
```shell
$ helm init
```

## Configure Spinnaker

First, update the values file with your project id.
```shell
$ sed -i.bak s/REPLACE_ME/$(gcloud info --format='value(config.project)')/g ./config/values.yaml
```

Next, copy your service account credentials into values.yaml.

Copy the output form the following command into values.yaml, replacing ```<SERVICE_ACCOUNT_JSON>```

```shell
$ SERVICE_ACCOUNT_JSON=$(cat account.json) && echo $SERVICE_ACCOUNT_JSON
```
TODO: Make this a sed operation
```shell
$ nano ./config/values.yaml 
```
 
```shell
# Disable minio the default
minio:
  enabled: false
 
# Enable gcs
gcs:
  enabled: true
  project: <my-project-name>
  jsonKey: '<SERVICE_ACCOUNT_JSON>'
 
# Name has to be unique in GCS.
storageBucket: <my-project-name>-spinnaker 
 
# Configure your Docker registries here
accounts:
- name: gcr
  address: https://gcr.io
  username: _json_key
  password: '<SERVICE_ACCOUNT_JSON>'
  email: 1234@5678.com
```

## Deploy Spinnaker Chart

### Temporary get the updated chart
```shell
$ git clone https://github.com/kubernetes/charts && cd charts

$ git fetch origin pull/1338/head:test-chart
$ git checkout test-chart

# update the dependencies
$ cd stable/spinnaker/
$ helm dep up 

$ cd ../../../
```

## Install spinnaker
NOTE: This is going to take a while. 
```shell
$ helm install ./charts/stable/spinnaker --name cd -f ./config/values.yaml --timeout 600
```
In another tab, you can monitor the progress of the installation. 
Errors will happen this is to be expected while the pods sync up.
```shell
$ kubectl get pods -w
```

## Access the spinnaker UI

```shell
$ DECK_POD=$(kubectl get pods -l "component=deck,app=cd-spinnaker"  \
    -o jsonpath="{.items[0].metadata.name}")
$ kubectl port-forward $DECK_POD 9000 >>/dev/null &
```
 
Visit the Spinnaker UI by opening your browser to: http://127.0.0.1:9000


### Misc / Monitor Progresss / Troubleshoot / Debug Installation
Everything in this helm chart will be labeled cd-spinnaker, so you can search for things like: 
```shell
$ kubectl get deployment -l app=cd-spinnaker
```

NOTE: If you want to make a quick change.  Helm can do a blue green deployment via upgrade.
```shell
$ helm upgrade cd ./charts/stable/spinnaker -f updated-values.yaml
```

Debugging can be done with
```shell
$ kubectl logs <pod name>
```

TODO:  Move this somewhere else.
To delete everything
```shell
$ helm delete cd --purge
```
 
get the latest list of charts
```shell
$ helm repo update
```

Common problems
Front50 will fail if the GCS bucket name is not unique.
