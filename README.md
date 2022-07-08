# Kecp: Webrtc Video Streaming Tool

## Introduction

Kecp is a webrtc video streaming tool for allowing peers to watch the same video simultaneously. Moreover, a small chatroom is also supported! For now, Kecp is still under construction. API may experience breaking changes.

## Get started

If you are not a developer who intends to use this tool as a library, this tool is generally available for you to use.

### Screenshots

![](assets/AB2D8BF.png)
![](assets/CF93C70.png)

### Config

Create a config called `config.toml` at the project's root folder.

Non-production example:

```toml
[server]
debug = true
tls = false
host = "127.0.0.1:8090"
allowed_origins = [""]
```

Production example:

```toml
[server]
debug = false
tls = true
host = "example.com"
allowed_origins = [""]
```

### Build

```shell
make build
```

### Run

```shell
make run
```

## License

Licensed under the Apache License, Version 2.0
