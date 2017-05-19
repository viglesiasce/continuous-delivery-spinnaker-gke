FROM golang:1.8 
ENV SOURCE=/go/src/github.com/viglesiasce/gke-info \
    GLIDE_VERSION=v0.12.3
RUN wget -q https://github.com/Masterminds/glide/releases/download/${GLIDE_VERSION}/glide-${GLIDE_VERSION}-linux-amd64.tar.gz \
    && tar zxfv glide-${GLIDE_VERSION}-linux-amd64.tar.gz \
    && mv linux-amd64/glide /usr/local/bin
COPY ./glide* $SOURCE/
RUN cd $SOURCE && glide install
COPY . $SOURCE
WORKDIR $SOURCE/cmd/gke-info
RUN go build -o gke-info
CMD ./gke-info