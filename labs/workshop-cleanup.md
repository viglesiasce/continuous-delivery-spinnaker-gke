# Workshop Cleanup

1. Delete the Spinnaker installation
../helm delete --purge cd
1. Delete the sample app services
```shell
kubectl delete -f k8s/services
```

1. Delete the service account:
```shell
export SA_EMAIL=$(gcloud iam service-accounts list --filter="displayName:spinnaker-storage-account" --format='value(email)')
gcloud iam service-accounts delete $SA_EMAIL
```

1. Delete the Container Engine cluster:
```shell
gcloud container clusters delete spinnaker-tutorial
```

1. Delete Cloud Source Repository:
```shell
gcloud source repos delete sample-app
```

1. Delete the bucket:
```shell
export PROJECT=$(gcloud info --format='value(config.project)')
export BUCKET=$PROJECT-spinnaker-config
gsutil -m rm -r gs://$BUCKET
```

1. Delete your container images:
```shell
export PROJECT=$(gcloud info --format='value(config.project)')
gcloud container images delete gcr.io/$PROJECT/sample-app:v1.0.0
gcloud container images delete gcr.io/$PROJECT/sample-app:v1.0.1
```
