FROM ubuntu:24.04 AS builder

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    build-essential \
    cmake \
    protobuf-compiler-grpc \
    libgrpc++-dev \
    protobuf-compiler \
    libprotobuf-dev \
    pkg-config

WORKDIR /app
COPY . .
RUN rm -rf build && mkdir build && cd build && cmake .. && make

FROM ubuntu:24.04
RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    libgrpc++-dev \
    libprotobuf-dev \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*
    
WORKDIR /app
COPY --from=builder /app/build/spm_agent .

CMD ["./spm_agent", "spm_api:50051"]
