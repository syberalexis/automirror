# Automirror

[![Build Status](https://travis-ci.com/syberalexis/automirror.svg)][travis]
[![Go Report Card](https://goreportcard.com/badge/github.com/syberalexis/automirror)](https://goreportcard.com/report/github.com/syberalexis/automirror)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/3394/badge)](https://bestpractices.coreinfrastructure.org/projects/3394/badge)

This Project is a POC to learn Golang. It's used to mirror different repositories
in local for a disconnected environment or embedded systems.

## Concept

This software is based on two engines : pullers and pushers.

### Pullers

This engine is configured with a source to download binaries, sources,
or something else into a destination and a special file configuration.

### Pushers

This engine is configured with a destination and used to push something into a platform, 
like Artifactory, Storage server, GitLab, or something else.

## Architecture



## Download

There ar two ways to get Automirror.

### Precompiled binaries

Precompiled binaries for released versions are available in the
[*github releases* section](https://github.com/syberalexis/automirror/releases).

### Building from source

To build Automirror from the source code yourself you need to have a working
Go environment with [version 1.13 or greater installed](https://golang.org/doc/install).

You can clone the repository yourself and build using `make` to compile it in
binary:

    $ mkdir -p $GOPATH/src/github.com/syberalexis
    $ cd $GOPATH/src/github.com/syberalexis
    $ git clone https://github.com/syberalexis/automirror.git
    $ cd automirror
    $ make
    $ mkdir -p /etc/automirror
    $ cp examples/config.toml /etc/automirror/config.toml
    $ ./dist/automirror

The Makefile provides several targets:

  * *build*: build the `automirror` binary (default)
  * *test*: run the tests
  * *clean*: clean buildings

## Usage

### Standalone

    $ ./automirror

### Systemd Service

    [Service]
    

## How to configure

Configuration is write in TOML. To accept more types, it's might be developed.

### Automirror

 * log_file : the location file to write logs
 * log_format : text (default) or json
 * log_level : the log level defined by logrus
 * mirrors : define activated mirrors in a map
 
#### Mirror

 * timer : Schedule timer to pull and push
 * puller : Activated puller configuration
    * name : Puller name to print in logs
    * source : Mirror source
    * destination : Mirror destination
    * config : Specific config file
 * pusher : Activated pusher configuration
    * name : Pusher name to print in logs
    * config : Specific config file

#### Example

    log_file = "/var/log/automirror.log"
    log_format = "json"
    [mirrors]
        [mirrors.MIRROR_NAME]
            timer = "24h"
            [mirrors.MIRROR_NAME.puller]
                name = "NAME"
                source = "URL"
                destination = "/tmp/NAME"
                config = "PULLER.toml"
            [mirrors.MIRROR_NAME.pusher]
                name = "NAME"
                config = "PUSHER.toml"

### Pullers

#### Deb

    dist = "stretch,stretch-updates,stretch-backports"
    arch = "amd64"
    section = "main,contrib,non-free,main/debian-installer"
    root = "/debian"
    method = "http"
    keyring = "/usr/share/keyrings/debian-archive-keyring.gpg"
    cleanup = false
    source = true
    i18n = true
    options = "--progress"

#### Docker

    auth_uri = "https://auth.docker.io/token?service=registry.docker.io"
    [[images]]
        name = "alpine"
    [[images]]
        name = "centos"
    [[images]]
        name = "ubuntu"

#### Git

    options = "--progress"

#### Maven

    metadata_file_name = "maven-metadata.xml"
    pom_file = "/tmp/pom.xml"
    database_file = "maven.db"
    
    [[artifact]]
        group = "org.springframework.boot"
        id = "spring-boot"
        minimum_version = "2.2.0.RELEASE"

#### Python

    database_file = "pip.db"
    file_extensions = "tar.gz|whl|zip|bz2|tar.bz2"
    sleep_timer = "500ms"

#### Rsync

    options = "-a -v"

#### Wget

    options = "-nH -nd -N"

### Pushers

#### JFrog

    url = "http://localhost:8081/artifactory/test"
    api_key = "AKCp5e2qXnFDWrtw7hJHjjWxR6ei5tCQ3HCvdnSYop6Y8w1vK1GQeUEKeFqSePJXmpCHexcac"
    source = "/tmp/maven"
    exclude_regexp = "(_remote.repositories)|.*\\.(sha1|md5)$"

## Contributing

Refer to [CONTRIBUTING.md](https://github.com/syberalexis/automirror/blob/master/CONTRIBUTING.md)

## License

Apache License 2.0, see [LICENSE](https://github.com/syberalexis/automirror/blob/master/LICENSE).

[travis]: https://travis-ci.com/syberalexis/automirror