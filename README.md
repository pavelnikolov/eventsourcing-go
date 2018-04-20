# Event Sourcing + CQRS using Go Tutorial

## Background

!["Fairflax Media logo"](https://github.com/pavelnikolov/eventsourcing-go/blob/master/static/fairflax-media-logo-small.png?raw=true "Fairflax Media logo")

Imagine that you are working for a large news publisher called **Fairflax Media**, inspired by [flax flower](https://www.google.com.au/search?q=flax+flower&tbm=isch) (not to be mistaken with Fairfax Media, of course).

You have been tasked with breaking the huge monolith system into microservices and moving everything to Kuberentes.
You prefer domain driven design, and believe that each database should only ever be owned by a single application. But as time goes by, you notice that many services need common data. You need to make sure that the data is consistent across all services and the query performance is reasonable.

## Project Goals
* Create a step-by-step tutorial for (opinionated) asynchronous microservices using Golang.
* Each step of the tutorial should be in a separate git branch.
* The `master` branch contains the initial state - tightly coupled services.
* The `event-sourcing` branch should contain the Event Sourcing + CQRS implementation, containing asynchronous services, which are independent and easy to test.
* Present actual problems and solve them using Event Sourcing + CQRS in Golang (with a few dependencies)
* Use minimalistic approach and try to use fewer moving parts. (For example, the services are intentionally kept basic and do not include proper logging, tracing, metrics, feature flags, rate-limitting, docker containers, authentication, authorisation, K8s helm charts, docker-compose etc.).
* In the initial branch DB is mocked up intentionally.
* Demonstrate the use of in-memory data stores (e.g. BoltDB and Bleve).
* Be more detailed and realistic than the regular _Hellow World_ example, but remain minimalistic.
* Demonstrate the use of channels for the basic event broker.

## Non-Goals
* Try to be everything for everyone.
* Demonstrate how to use Kafka/Kafka Streams.
* Demonstrate how to use DistributedLog.
* Guarantee exactly-once delivery.
* Show how to save/restore snapshots and use pub/sub systems with limited retention periods (e.g. AWS Kinesis and Google Cloud Pub/Sub).
* Implement server-less version, for example, using AWS Lambda or similar.
* Use event sourcing with a decentralized events ledger/blockchain.

## Target Audience
This project is intended for anyone willing to switch to a asynchronous microservice architecture using Golang.

## Prerequisites

- Go 1.9 or later
- [`$GOPATH`](https://golang.org/doc/code.html#GOPATH) is set
- [`dep`](https://github.com/golang/dep#installation) for managing dependencies

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


## Quick install

```
$ go get -u github.com/pavelnikolov/eventsourcing-go/demo-articles
$ go get -u github.com/pavelnikolov/eventsourcing-go/demo-graph
$ go get -u github.com/pavelnikolov/eventsourcing-go/demo-sitemap
$ go get -u github.com/pavelnikolov/eventsourcing-go/demo-rss
```

## Clone the repository

Either:
```
$ go get -u github.com/pavelnikolov/eventsourcing-go/...
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
Run in four different terminal windows:

```
go install ./cmd/demo-articles && demo-articles
go install ./cmd/demo-graph && demo-graph
go install ./cmd/demo-rss && demo-rss
go install ./cmd/demo-sitemap && demo-sitemap
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
