# ============================
#     1. BUILDER STAGE
# ============================
FROM debian:stable-slim AS builder
ENV DEBIAN_FRONTEND=noninteractive

# Install build deps + docbook tools (fixes xmllint + a2x errors)
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
      build-essential \
      python3 \
      procps \
      wget \
      xz-utils \
      git \
      libcap-dev \
      libsystemd-dev \
      pkgconf \
      pkgconf-bin \
      asciidoc \
      libxml2-utils \
      xsltproc \
      docbook-xsl \
      docbook-xml \
      python3-docutils \
      ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Clone isolate (shallow for speed)
RUN git clone --depth 1 https://github.com/ioi/isolate.git /tmp/isolate-src

# Build isolate
WORKDIR /tmp/isolate-src
RUN make && make install DESTDIR=/tmp/isolate-install

# ============================
#   2. RUNTIME STAGE (CLEAN)
# ============================
FROM debian:stable-slim AS runtime
ENV DEBIAN_FRONTEND=noninteractive

# Only install runtime deps needed by isolate
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
      libcap2-bin \
      procps \
      ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy installed isolate files from builder
COPY --from=builder /tmp/isolate-install/ /

# Set correct privileges for isolate
# isolate MUST be setuid-root or it cannot create namespaces/mounts
RUN chown root:root /usr/local/bin/isolate && \
    chmod 4755 /usr/local/bin/isolate

# Create non-root user for sandbox execution
RUN useradd -m -s /usr/sbin/nologin judgeuser && \
    mkdir -p /app /var/local/lib/isolate && \
    chown -R judgeuser:judgeuser /app /var/local/lib/isolate

# Switch to restricted user
USER judgeuser
WORKDIR /app

# Default shell (just a placeholder)
CMD ["/bin/bash"]
