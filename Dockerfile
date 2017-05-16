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
RUN wget -q https://github.com/Masterminds/glide/releases/download/v0.12.3/glide-v0.12.3-linux-amd64.tar.gz \
    && tar zxfv glide-v0.12.3-linux-amd64.tar.gz \
    && mv linux-amd64/glide /usr/local/bin
RUN wget -q -O go-cloud-debug https://storage.googleapis.com/cloud-debugger/compute-go/go-cloud-debug \
    && chmod 0755 go-cloud-debug \
    && mv go-cloud-debug /usr/local/bin
COPY . $SOURCE
RUN cd $SOURCE && glide install
WORKDIR $SOURCE/cmd/gke-info
RUN go build -o gke-info
RUN gcloud debug source gen-repo-info-file
CMD ["bash", "-c", "$SOURCE/run-app.sh"]