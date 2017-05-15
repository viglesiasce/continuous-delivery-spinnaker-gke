FROM golang:1.8 
ENV SOURCE=/go/src/github.com/viglesiasce/gke-info \
    PATH=/opt/google-cloud-sdk/bin:$PATH \
    GOOGLE_CLOUD_SDK_VERSION=154.0.1
RUN set -x \
  && cd /opt \
  && echo 'debconf debconf/frontend select Noninteractive' | debconf-set-selections \
  && apt-get update \
  && apt-get install --no-install-recommends -y jq wget python git localepurge ca-certificates \
  && wget -q https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-${GOOGLE_CLOUD_SDK_VERSION}-linux-x86_64.tar.gz \
  && tar zxfv google-cloud-sdk-${GOOGLE_CLOUD_SDK_VERSION}-linux-x86_64.tar.gz \
  && ./google-cloud-sdk/install.sh \
  && gcloud components install kubectl
COPY . $SOURCE
WORKDIR $SOURCE/cmd/gke-info
RUN go build -o gke-info
RUN wget -q -O go-cloud-debug https://storage.googleapis.com/cloud-debugger/compute-go/go-cloud-debug && \
    chmod 0755 go-cloud-debug
RUN gcloud debug source gen-repo-info-file
CMD ["./go-cloud-debug", "-v", "-sourcecontext=source-context.json", "-appmodule=gke-info", "-appversion=v6.1.5", "--", "./gke-info"]