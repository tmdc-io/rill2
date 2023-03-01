# syntax = docker/dockerfile:1.1-experimental
FROM golang:1.19

RUN curl -sL https://deb.nodesource.com/setup_18.x | bash - && \
    apt-get install -y nodejs

WORKDIR /go/src/

COPY . ./

RUN make cli

FROM ubuntu

RUN apt-get update && apt-get install -y ca-certificates

COPY --from=0 /go/src/rill /usr/local/bin
RUN chmod 777 /usr/local/bin/rill

RUN groupadd -g 1000 rill \
    && useradd -m -u 1000 -s /bin/sh -g rill rill

ENTRYPOINT ["rill"]
CMD ["start"]
