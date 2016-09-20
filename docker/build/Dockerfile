# Truss
# truss-build Docker image
#
# NOTE: YOU WILL PROBABLY NEVER NEED TO CREATE A TRUSS-BUILD IMAGE. The primary use
# of the truss-build image will be to create updated truss-run images A pre-built
# truss-build image from the Docker registry can build the binaries that are
# required for a Truss runtime image. Or more simply, a pre-built truss-run image
# can be used, if it meets your needs.
#
# This Dockerfile builds an image that can itself build Truss. The container
# that runs this image must mount the project source code to an external
# volume, so the go compiler can build it should have a volume mounted, so it
# can read the Truss code and generate binaries.
#
# The base of this Docker image is a full Linux distribution, because the protoc
# binary is dynamically linked and requires libraries that are not present in a
# minimal distribution like Alpine. We pull the required libraries from this base
# image and package them up to be added to the runtime image.
#
# E.g., docker run --rm -v $PWD:/go/src/github.com/TuneLab/go-truss tunelab/truss-build
#
# The binaries are subsequently packaged into a Docker image for running Truss.
# The runtime images is much smaller because of this.
FROM golang:1.7.1
MAINTAINER lab@tune.com

ENV PROTOBUF_VERSION 3.0.2
ENV PROTOBUF_ZIP protoc-$PROTOBUF_VERSION-linux-x86_64.zip
ENV TRUSS_REPO github.com/TuneLab/go-truss
ENV TRUSS_PATH $GOPATH/src/$TRUSS_REPO
ENV PROTO_GEN_TRUSS protoc-gen-truss-protocast

WORKDIR $TRUSS_PATH

# Get all the supporting libraries and programs necessary to build truss.
RUN apt-get update \
    && apt-get install -y ca-certificates wget unzip git \
    && update-ca-certificates \
    && go version \
    && go get -u -v github.com/golang/protobuf/proto \
        github.com/golang/protobuf/protoc-gen-go \
        github.com/go-kit/kit/...

# Get the Protobuf compiler and unzip it into the build output directory.
RUN wget -O /tmp/$PROTOBUF_ZIP https://github.com/google/protobuf/releases/download/v$PROTOBUF_VERSION/$PROTOBUF_ZIP \
    && (cd / && unzip /tmp/$PROTOBUF_ZIP bin/protoc) \
    && ldd /bin/protoc | egrep -o '/.+ ' | sed -e 's/^\///' > /tmp/protoc-libs.txt \
    && rm -rf /tmp/$PROTOBUF_ZIP

COPY docker/build/build_truss.sh /

ENTRYPOINT ["bash"]
CMD ["-c", "/build_truss.sh"]

# After execution of this image in container, $GOPATH/build contents will be on
# locally mounted volume, and can be placed into the Truss runtime image.
