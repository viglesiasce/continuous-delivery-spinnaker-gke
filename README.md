# GKE Info 

| Test          |   Result      |
| ------------- |---------------|
| Tutorial      | ![badge](https://concourse.dev.vicnastea.io/api/v1/teams/main/pipelines/gke-info-post-submit/jobs/test-tutorial/badge) |
| Build App     | ![badge](https://concourse.dev.vicnastea.io/api/v1/teams/main/pipelines/gke-info-post-submit/jobs/build-gke-info/badge)|

# Continuous Delivery with Spinnaker and Kubernetes (in 40 minutes or less)

You've got code. It probably compiles. Now what? 

It's time to push code into production, cross your fingers, and pray! Right? On second thought, we should probably test the code and ensure it works BEFORE releasing it to the rest of the world. Ideally, we'll do this using open source, multi-cloud tools that will work whether we're using Java or Go, on-premise or in the cloud.

And that's where [Kubernetes](https://www.kubernetes.io) and [Spinnaker](https://www.spinnaker.io), the continuous delivery platform come in. 
In this workshop we'll setup a CICD pipeline and explain the pros and cons, along the way.

[See a video of the pipleine you'll create here!](https://youtu.be/dpbWpzAs-RwD)

In this solution, you'll will:
* Setting up an example Continuous Delivery pipeline from scratch.
* Use your favorite tools (well, my favorite, at least) Kubernetes and Spinnaker
* Learn and overcome common Pitfalls and Obstacles of the above tools (because nothing's perfect)

## Labs

1. [Workshop Setup](labs/workshop-setup.md)
1. [Building Container Images](labs/building-container-images.md)
1. [Installing Spinnaker](labs/installing-spinnaker.md)
1. [Creating Your Pipeline](labs/creating-your-pipeline.md)
1. [Triggering Deployments](labs/triggering-deployments-v2.md)

## What Next?

* PLACEHOLDER
