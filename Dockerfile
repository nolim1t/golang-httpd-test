# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
# EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
# MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
# IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
# OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
# ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
# OTHER DEALINGS IN THE SOFTWARE.

# A lot of this borrows on https://github.com/lncm/invoicer/blob/master/Dockerfile and its concepts
# of not trusting binaries. 
# However I'm going to taylor this more for buildx as it seems more versatile.

# Version to be built
ARG VERSION=0.0.1

ARG VER_GO=1.15
ARG VER_ALPINE=3.12

ARG USER=httpapi
ARG DIR=/data/
ARG TAGS="static_build"

FROM golang:${VER_GO}-alpine${VER_ALPINE} AS alpine-builder

# Capture version and tags
ARG VERSION
ARG TAGS

ENV BINARY /go/bin/httpd
ENV LDFLAGS "-s -w -buildid= -X main.version=${VERSION}"
RUN apk add --no-cache  musl-dev  file  git  gcc

RUN mkdir -p /go/src/
COPY ./ /go/src/
WORKDIR /go/src/

RUN export GIT_HASH="$(git rev-parse HEAD)"; \
    echo "Building git tag: ${GIT_HASH}"; \
    go build  -x  -v  -trimpath  -mod=readonly  -tags="${TAGS}" \
        -ldflags="${LDFLAGS} -X main.gitHash=${GIT_HASH}" \
        -o "${BINARY}"


# Print rudimentary info about the built binary
RUN sha256sum   "${BINARY}"
RUN file -b     "${BINARY}"
RUN du          "${BINARY}"

FROM alpine:${VER_ALPINE} AS perms

ARG USER
ARG DIR

# NOTE: Default GID == UID == 1000
RUN adduser --disabled-password \
            --home ${DIR} \
            --gecos "" \
            ${USER}

# Needed to prevent `VOLUME ${DIR}` creating it with `root` as owner
USER ${USER}
RUN mkdir -p ${DIR}



#
## This is the final image that gets shipped to Docker Hub
#
FROM ${ARCH:+${ARCH}/}alpine:${VER_ALPINE} AS final

ARG USER
ARG DIR

LABEL maintainer="nolim1t (@nolim1t)"

# Copy only the relevant parts from the `perms` image
COPY  --from=perms /etc/group /etc/passwd /etc/shadow  /etc/

# From `perms`, copy *the contents* of `${DIR}` (ie. nothing), and set correct owner for destination `${DIR}`
COPY  --from=perms --chown=${USER}:${USER} ${DIR}  ${DIR}

# Copy the binary from the cross-check stage
COPY  --from=alpine-builder  /go/bin/httpd  /usr/local/bin/

USER ${USER}

# Expose the volume to communicate config, log, etc through (default: /data/)
VOLUME ${DIR}

# Expose port the servicer listens on
EXPOSE 8080

# Specify the start command and entrypoint as the httpd daemon.
ENTRYPOINT ["httpd"]

CMD ["-config", "/data/httpd.conf"]
