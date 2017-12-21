# README #

### What is this repository for? ###

* *.proto files and generated source code using gRPC framework 
* generated swagger specifications

### How do I get set up? ###

* Install docker (recommended). If you want to build without docker, use `make USE_DOCKER=0 ...`
* Run `make`

### Including to other projects ###
* Golang projects: `go get bitbucket.org/exonch/ch-grpc/<service-name>`
* Other projects: as git submodule

### Documentation ###
See DOCUMENTATION.html