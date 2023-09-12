FROM golang:1.18-bullseye

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    cmake \
    g++ \
    python3-pip \
    && apt-get clean -y  \
    && apt-get autoremove -y \
    && rm -rf /var/lib/apt/lists/*

ENV LIBRARY_PATH=/usr/local/lib${LIBRARY_PATH:+:$LIBRARY_PATH}
ENV LD_LIBRARY_PATH=/usr/local/lib${LD_LIBRARY_PATH:+:$LD_LIBRARY_PATH}

RUN cd /tmp/ && git clone https://github.com/dmlc/treelite.git -b 3.9.0 \
    && cd treelite \
    && mkdir build && cd build \
    && cmake .. \
    && make install -j $(nproc) \
    && rm -r /tmp/treelite

COPY . /src

WORKDIR /workspace

ENV GOPROXY='https://goproxy.io,direct'

RUN cd /src \
    && mkdir server \
    && go build -o /bin/server \
    && cp /bin/server /workspace/server \
    && rm -r /src

WORKDIR /workspaces

ENTRYPOINT ["/workspace/server"]

