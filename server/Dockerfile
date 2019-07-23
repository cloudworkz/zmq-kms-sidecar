FROM golang:1.12.7-stretch

RUN apt-get update && apt-get install -y apt-utils

# install zeromq libraries
RUN echo "deb http://download.opensuse.org/repositories/network:/messaging:/zeromq:/release-stable/Debian_9.0/ ./" >> /etc/apt/sources.list
RUN wget https://download.opensuse.org/repositories/network:/messaging:/zeromq:/release-stable/Debian_9.0/Release.key -O- | apt-key add
RUN apt-get update && apt-get install -y gcc build-essential libzmq3-dev

WORKDIR /app
COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -installsuffix cgo -o /go/bin/zmq-kms

RUN apt-get purge -y build-essential gcc apt-utils && \
    apt-get autoremove -y && \
    rm -rf /var/lib/apt/lists/*

ENTRYPOINT ["/go/bin/zmq-kms"]