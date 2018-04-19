# Event Sourcing + CQRS using Go Workshop

## Background

!["Fairflax Media logo"](https://github.com/pavelnikolov/eventsourcing-go/blob/master/static/fairflax-media-logo-small.png?raw=true "Fairflax Media logo")

Imagine that you are working for a large news publisher called **Fairflax Media**, inspired by [flax flower](https://www.google.com.au/search?q=flax+flower&tbm=isch) (not to be mistaken with Fairfax Media, of course).

You have been tasked with breaking the huge monolith system into microservices and moving everything to Kuberentes.
You prefer domain driven design, and believe that each database should only ever be owned by a single application. But as time goes by, you notice that many services need common data. You need to make sure that the data is consistent across all services and the query performance is reasonable.

## Prerequisites

- Go 1.6 or later
- `GOPATH` [is set](https://golang.org/doc/code.html#GOPATH)
- `dep` [for managing dependencies](https://github.com/golang/dep#installation)

  On Mac OSX using `brew`...
  ```
  $ brew install dep
  $ brew upgrade dep
  ```
  On other platforms you can use the install.sh script:
  ```
  $ curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
  ```

- Add ensure the PATH contains $GOPATH/bin
  ```
  $ go help gopath
  $ export PATH=$PATH:$GOPATH/bin
  ```


## Quick install (Untested)

```
$ go get -u github.com/pavelnikolov/eventsourcing-go/demo-articles
$ go get -u github.com/pavelnikolov/eventsourcing-go/demo-graph
$ go get -u github.com/pavelnikolov/eventsourcing-go/demo-sitemap
$ go get -u github.com/pavelnikolov/eventsourcing-go/demo-rss
```

## Clone the repository

Either (I haven't tried it):
```
$ go get github.com/pavelnikolov/eventsourcing-go/...

```
Or:
```
$ mkdir -p $GOPATH/github.com/pavelnikolov
$ cd $GOPATH/github.com/pavelnikolov
$ git clone github.com/pavelnikolov/eventsourcing-go
```

Then install the dependencies:
```
$ dep ensure
```

## Try it!

- Run the articles API

  ```
  $ demo-articles &
  ```

- Run the GraphQL gateway API

  ```
  $ demo-graph
  ```

- Run the RSS generator

  ```
  $ demo-rss
  ```

- Run the Sitemap generator

  ```
  $ demo-sitemap
  ```

Navigate to the apps in your browser:
- GraphiQL UI - http://localhost:4001/
- Latest news RSS feed - http://localhost:4002/feed
- Latest business news RSS feed - http://localhost:4002/feed/business
- Latest political news RSS feed - http://localhost:4002/feed/politics
- (Naive and useless) Sitemap - http://localhost:4003/sitemap


## Optional tasks

### Rebuild the generated code

1. Install [protobuf compiler](https://github.com/google/protobuf/blob/master/README.md#protocol-compiler-installation)

1. Install the protoc Go plugin

   ```
   $ go get -u github.com/golang/protobuf/protoc-gen-go
   ```

1. Rebuild the generated Go code

   ```
   $ go generate github.com/pavelnikolov/eventsourcing-go/...
   ```
   
   Or run `protoc` command (with the grpc plugin)
   
   ```
   $ protoc -I publishing/ publishing/publishing.proto --go_out=plugins=grpc:publishing
   ```

### Rebuild  and run the apps 

```
go install ./cmd/demo-articles && demo-articles
go install ./cmd/demo-graph && demo-graph
go install ./cmd/demo-rss && demo-rss
go install ./cmd/demo-sitemap && demo-sitemap
```
