# Use a base Node.js image
FROM node:20 AS builder

# Install Go
RUN wget -O go.tar.gz https://golang.org/dl/go1.22.1.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go.tar.gz && \
    rm go.tar.gz

# Set Go environment variables
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/go"
ENV PATH="${GOPATH}/bin:${PATH}"

# Copy the Go application files
WORKDIR /app
COPY . .

# Build the Go application
RUN make cli

# Use a base ubuntu for the final image
FROM python:3.9-slim

RUN apt-get update && apt-get install -y ca-certificates

# Copy the built Go application from the builder image
COPY --from=builder /app/rill /usr/local/bin

RUN chmod 777 /usr/local/bin/rill

RUN groupadd -g 1001 rill \
    && useradd -m -u 1001 -s /bin/sh -g rill rill

RUN rill runtime install-duckdb-extensions

WORKDIR /etc/dataos/work
COPY lens2 /lens2
RUN python3 -m pip install -r /lens2/requirements.txt

# Start the application
# ENTRYPOINT ["rill"]
# CMD ["start"]