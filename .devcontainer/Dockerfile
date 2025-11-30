FROM mcr.microsoft.com/devcontainers/go:1.25-trixie

# Install additional tools including zsh
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    protobuf-compiler \
    libprotobuf-dev \
    zsh \
    && apt-get clean -y && rm -rf /var/lib/apt/lists/*

# Install go-task
RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin

# Set up Go environment
ENV GOPATH=/go
ENV PATH=$PATH:$GOPATH/bin

# Create workspace directory
RUN mkdir -p /workspace
WORKDIR /workspace

# Change default shell to zsh for vscode user
RUN chsh -s /bin/zsh vscode
# Ensure vscode user owns the workspace directory
RUN chown -R vscode:vscode /workspace
RUN chown -R vscode:vscode /home/vscode/.oh-my-zsh/custom
# Ensure vscode user has permission to go modules directory
RUN chown -R vscode:vscode ${GOPATH}

# Copy custom oh-my-zsh configuration
COPY ./.oh-my-zsh/custom /home/vscode/.oh-my-zsh/custom
RUN chown -R vscode:vscode /home/vscode/.oh-my-zsh/custom
