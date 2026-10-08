FROM debian:12-slim

# ─── System dependencies ──────────────────────────────────────────────────
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    wget \
    git \
    make \
    gcc \
    g++ \
    ca-certificates \
    gnupg \
    lsb-release \
    unzip \
    sudo \
    bash \
    pkg-config \
    libgtk-3-dev \
    libwebkit2gtk-4.0-dev \
    && rm -rf /var/lib/apt/lists/*

# ─── Non-root user (mirrors the vscode user convention) ───────────────────
ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=1000
RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -m -s /bin/bash $USERNAME \
    && echo "$USERNAME ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers.d/$USERNAME \
    && chmod 0440 /etc/sudoers.d/$USERNAME

# ─── Go ───────────────────────────────────────────────────────────────────
ARG GO_VERSION=1.23.4
RUN curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" \
    | tar -xz -C /usr/local
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH=/go
ENV PATH="${GOPATH}/bin:${PATH}"
# Give the vscode user ownership of the Go workspace
RUN mkdir -p /go && chown -R $USERNAME:$USERNAME /go

# Go tools — installed as root so they land in /go/bin on PATH for all users
RUN go install golang.org/x/tools/gopls@latest \
    && go install github.com/go-delve/delve/cmd/dlv@latest \
    && go install golang.org/x/tools/cmd/goimports@latest \
    && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest \
    && go install gotest.tools/gotestsum@latest \
    && go install github.com/wailsapp/wails/v2/cmd/wails@latest

# ─── Node.js (via nvm) ────────────────────────────────────────────────────
ARG NODE_VERSION=22
ENV NVM_DIR=/usr/local/nvm
RUN mkdir -p $NVM_DIR \
    && curl -fsSL https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.1/install.sh | bash \
    && . "$NVM_DIR/nvm.sh" \
    && nvm install ${NODE_VERSION} \
    && nvm use ${NODE_VERSION} \
    && nvm alias default ${NODE_VERSION} \
    && npm install -g pnpm \
    # Make nvm dir accessible to the vscode user
    && chown -R $USERNAME:$USERNAME $NVM_DIR

ENV PATH="$NVM_DIR/versions/node/v${NODE_VERSION}.0/bin:${PATH}"

# ─── Final setup ──────────────────────────────────────────────────────────
WORKDIR /workspace
RUN chown $USERNAME:$USERNAME /workspace

# Pre-create the full vscode-server + Cline globalStorage path as vscode-owned.
# This ensures the named Docker volume (cline-config) that is mounted at the
# Cline globalStorage path inherits vscode ownership instead of being created
# as root by Docker, which would block VS Code Server from writing anywhere
# under ~/.vscode-server at container start-up.
RUN mkdir -p /home/$USERNAME/.vscode-server/data/User/globalStorage/saoudrizwan.claude-dev \
    && chown -R $USERNAME:$USERNAME /home/$USERNAME/.vscode-server

USER vscode

