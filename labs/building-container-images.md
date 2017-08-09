# Building Container Images

In this lab, we'll use a build trigger to connect Google Source Repository to Google Container Registry. 
After this module, whenever we tag an image for release, we'll kick off a spinnaker deployment.

## Set up the Build Trigger

![](../docs/img/Setup-build-trigger.png)

In the Google Cloud Console, set the build trigger to: 
 'Changes pushed to "v.*" tag will trigger a build of "gcr.io/askcarter-production-gke/$REPO_NAME:$TAG_NAME"'


## Create a git repository for the code
```shell
$ cd sample-app
$ git init
$ git add .
$ git commit -m "Intial commit."
```

## Set up Cloud Source Repository
Back on the command line, set up the cloud repository.

```shell
$ gcloud beta source repos create gceme
$ git config credential.helper gcloud.sh
$ git remote add google https://source.developers.google.com/p/$PROJECT/r/gceme
$ git push --all google
```

## Tag an image, Trigger a build

Now, use git to tag a commit and trigger the build.

```shell
$ git tag -a v1.0.0 -m "my version 1.0.0"
$ git push google v1.0.0
```

Back in the Container Registry UI's Build history, our build should've kicked off.
