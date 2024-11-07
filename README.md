# Dockershrink CLI

![Typical interaction with dockershrink CLI](./static/dockershrink-how-it-works.gif)

Official commandline application to interact with the [Dockershrink](https://dockershrink.com) platform.

**Dockershrink is a SaaS platform that helps you reduce the size of your Nodejs Docker images by applying code optimizations.**

It combines the power of traditional rule-based analysis with Generative AI to yield state-of-the-art optimizations for your images :brain:

### NOTE
Dockershrink is in **BETA**

We would love to hear what you think! You can provide your feedback from your [dockershrink dashboard](https://dockershrink.com/dashboard) or [Create an Issue](https://github.com/duaraghav8/dockershrink-cli/issues) in this repository.

## Why does dockershrink exist?
Every org using containers in development or production environments understands the pain of managing hundreds or even thousands of BLOATED Docker images in their infrastructure.

But not everyone realizes that by just implementing some basic techniques, they can reduce the size of a 1GB Docker image down to **as little as 100 MB**!

Imagine the costs saved in storage & data transfer, decrease in build times AND the productivity gains of developers :exploding_head:

Dockershrink aims to auomatically apply advanced optimization techniques such as Multistage builds, Light base images, removing unused dependencies, etc. so that developers & devops engineers don't have to waste time and everybody reaps the benefits.

You're welcome :wink:

## How it works
The CLI is the primary user-facing application of dockershrink.

When you invoke it on a project, it communicates with the Dockershrink platform to analyze some of your project files.

Currently, these files are `Dockerfile`, `package.json` and `.dockerignore`.

It then creates a new directory called `dockershrink.optimized` inside the project, which contains modified versions of your files that will result in a smaller Docker Image.

The CLI outputs a list of actions it took on your files.

It may also include suggestions on further improvements you could make.

## Setup
If you haven't already

1. Download the pre-built CLI for your platform from the Releases page.

Alternatively, you can also clone this repository and build the binary yourself by following the build instructions below.

2. Rename your binary to `dockershrink` and make sure that the binary is in your `PATH` environment variable.

For example, on MacOS and Linux, you could copy the binary to your local bin directory using
```bash
cp ./dockershrink /usr/local/bin/dockershrink
```

3. Check that everything is setup by running `dockershrink help`.

4. Initialize dockershrink with your API Key.

Copy your platform API Key from your account dashboard once you log into dockershrink.
Then run the following command in your terminal:

```bash
dockershrink init --api-key <paste your api key here>
```

This is a one-time step.

Congratulations! Your instalation is now complete!

Head over to Usage.

## Usage
Navigate into the root directory of one of your Node.js projects and run the simplest command:

```bash
dockershrink optimize
```

Dockershrink will create a new directory with the optimized files and output the actions taken and (maybe) some more suggestions.

For more information on the `optimize` command, run
```bash
dockershrink help optimize
```

### Using AI Features

**NOTE**: By default, dockershrink only runs rule-based analysis to optimize your image definition.
If you want to enable AI, you must supply your own [OpenAI API Key](https://openai.com/index/openai-api/).

```bash
dockershrink optimize --openai-api-key <your openai api key>

# Alternatively, you can supply the key as an environment variable
export OPENAI_API_KEY=<your openai api key>
dockershrink optimize
```

We **highly recommend** you enable AI for more intelligent & powerful optimizations.

### Default file paths

**NOTE**: By default, the CLI looks for the files to optimize in the current directory.

You can also specify the paths to all files using options (see `dockershrink help optimize` for the available options).

## Build from source

### Prerequisites

- Clone this repository on to your local machine.
- Make sure [Golang](https://golang.org/dl/) is installed on your system (at least version 1.23)
- Make sure [Docker](https://www.docker.com/get-started) installed on your system and the Docker daemon is running.

### How to build the binaries
To compile the code and build binaries for all supported platforms and architectures, run:
```bash
make build-all
```

Alternatively, build the binaries yourself by running the following commands inside a Docker container:

```bash
docker build -t dockershrink-builder .

docker run --rm -v "$PWD":/app -w /app dockershrink-builder sh -c "\
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/dockershrink-linux-amd64 main.go &&\
    CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/dockershrink-linux-arm64 main.go &&\
    CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/dockershrink-darwin-amd64 main.go &&\
    CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o bin/dockershrink-darwin-arm64 main.go &&\
    CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/dockershrink-windows-amd64.exe main.go &&\
    CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o bin/dockershrink-windows-arm64.exe main.go\
"
```

This will create all the binaries inside the `bin/` directory.
