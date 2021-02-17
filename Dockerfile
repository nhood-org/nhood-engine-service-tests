FROM golang:1.16.0

ENV GO111MODULE=on

WORKDIR /usr/src/nhood-org/

COPY . .
RUN make install-dependencies && \
    make install-tools
    
CMD make run
