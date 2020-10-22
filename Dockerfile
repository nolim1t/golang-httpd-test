# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
# EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
# MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
# IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
# OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
# ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
# OTHER DEALINGS IN THE SOFTWARE.

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

RUN mkdir -p /go/bin/
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
