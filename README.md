# Dockershrink CLI

![Typical interaction with dockershrink CLI](./static/dockershrink-how-it-works.gif)

Official commandline application to interact with the [Dockershrink](https://dockershrink.com) backend.

**Dockershrink is a SaaS platform that helps you reduce the size of your Nodejs Docker images by applying code optimizations.**

It combines the power of traditional rule-based analysis with Generative AI to yield state-of-the-art optimizations for your images :brain:

> [!IMPORTANT]
> Dockershrink is in **BETA** and is **FREE** to use.
> 
> We would love to hear what you think! You can provide your feedback from your [dockershrink dashboard](https://dockershrink.com/dashboard) or [Create an Issue](https://github.com/duaraghav8/dockershrink-cli/issues) in this repository.

## Why does dockershrink exist?
Every org using containers in development or production environments understands the pain of managing hundreds or even thousands of BLOATED Docker images in their infrastructure.

But not everyone realizes that by just implementing some basic techniques, they can reduce the size of a 1GB Docker image down to **as little as 100 MB**!

([I also made a video on how to do this.](https://youtu.be/vHBHxQfK6cM))

Imagine the costs saved in storage & data transfer, decrease in build times AND the productivity gains for developers :exploding_head:

Dockershrink aims to auomatically apply advanced optimization techniques such as Multistage builds, Light base images, removing unused dependencies, etc. so that developers & devops engineers don't have to waste time doing so and everybody still reaps the benefits!

You're welcome :wink:

## How it works
The CLI is the primary user-facing application of dockershrink.

When you invoke it on a project, it communicates with the Dockershrink platform to analyze some of your project files.

Currently, these files are:

:point_right: `Dockerfile`

:point_right: `package.json`

:point_right: `.dockerignore`

It then creates a new directory called `dockershrink.optimized` inside the project, which contains modified versions of your files that will result in a smaller Docker Image.

The CLI outputs a list of actions it took on your files.

It may also include suggestions on further improvements you could make.

## Setup
If you haven't already, [signup for an account](https://dockershrink.com) and get your API Key from your dashboard. Then proceed with the below steps.

1. Download the pre-built CLI for your platform from the [Releases](https://github.com/duaraghav8/dockershrink-cli/releases) page.

On MacOS, use brew to install the CLI:
```bash
brew install duaraghav8/tap/dockershrink
```

Alternatively, you can clone this repository and build the binary yourself by following the build instructions below.

2. Rename your binary to `dockershrink` and make sure that the binary is in your `PATH` environment variable.

For example, on MacOS and Linux, you could copy the binary to your local bin directory using
```bash
cp ./dockershrink /usr/local/bin/dockershrink
```

3. Check that everything is setup by running `dockershrink help`.

4. Initialize dockershrink with your API Key.

Copy your API Key from your dashboard once you log into dockershrink.
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

> [!NOTE]
> Using AI features is optional, but **highly recommended** for more customized and powerful optimizations.
>
> Currently, you need to supply your own openai api key, so even though Dockershrink itself is free, openai usage might incur some cost for you.

By default, dockershrink only runs rule-based analysis to optimize your image definition.

If you want to enable AI, you must supply your own [OpenAI API Key](https://openai.com/index/openai-api/).

```bash
dockershrink optimize --openai-api-key <your openai api key>

# Alternatively, you can supply the key as an environment variable
export OPENAI_API_KEY=<your openai api key>
dockershrink optimize
```

> [!NOTE]
> Dockershrink does not store your OpenAI API Key.
>
> So you must provide your key every time you want "optimize" to enable AI features.

### Default file paths
By default, the CLI looks for the files to optimize in the current directory.

You can also specify the paths to all files using options (see `dockershrink help optimize` for the available options).

## Build from source

### Prerequisites

- Clone this repository on to your local machine.
- Make sure [Golang](https://golang.org/dl/) is installed on your system (at least version 1.23)
- Make sure [Docker](https://www.docker.com/get-started) installed on your system and the Docker daemon is running.
- Install [GoReleaser](https://goreleaser.com/) (at least version 2.4)

### Build for local testing
```bash
# Single binary
goreleaser build --single-target --clean --snapshot

# All binaries
goreleaser release --snapshot --clean
```

### Create a new release
1. Create a Git Tag with the new version

```bash
git tag -a v0.1.0 -m "Release version 0.1.0"
git push origin v0.1.0
```

2. Release
```bash
# Make sure GPG is present on your system and you have a default key which is added to Github.

# set your github access token
export GITHUB_TOKEN="<your GH token>"

goreleaser release --clean
```
