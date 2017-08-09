## Trigger a Deployments from the Command Line
Now, we can use git to tag a commit and trigger the build.

```shell
$ git tag -a v1.0.1 -m "my version 1.0.1"
$ git push google v1.0.1
```

Back in the Cloud Console, we should see our image being built.
After the image is in GCR, our build should kick off in Spinnaker.
This may take a few minutes.

## Rolling Back a Deployment from the Command Line
Use git to go back a commit, then push the image and bump the tag.
I workshop have user make two commits (the 2nd of which is bad).  When the user pushes them, have them follow this process.

```shell
# Change code
$ git add <change>
$ git commit -m "Working code fix."

# Change code
$ git add <change>
$ git commit -m "Broken code fix."

$ git push google master

$ git tag -a v3.0.0 -m "my version 3.0.0"
$ git push google v3.0.0
```

Roll back change
```shell
$ git <rollback one commit>
$ git commit -m "Rolled back buggy code."
$ git push google master
$ git tag -a v4.0.0 -m "my version 4.0.0"
$ git push google v4.0.0
```

## Trigger a Builds and Rolling Back changes from the UI
You can also use the UI if you need an escape hatch.

