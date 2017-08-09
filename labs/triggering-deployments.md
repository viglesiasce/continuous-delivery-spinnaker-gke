# Triggering Deployments

## Setting up Source Control
TODO: This should before here, possibly in the overview

NOTE:  This may not be necessary here.  
```shell
$ git clone https://github.com/askcarter/spinnaker-k8s-workshop
$ cd spinnaker-k8s-workshop
```

Pushing the image to gcr. 
Note: Add v to semver so that we can use the v* tag as a build trigger
```shell
$ gcloud container builds submit -t gcr.io/askcarter-production-gke/gceme:v1.0.0 .
```

View everything in cluster.
```shell
$ kubectl get pods,service,deployments -l app=demo
```
