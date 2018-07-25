# Mail

Mail is a mail server and template manager for [Containerum](https://github.com/containerum/containerum).

## Features
* Direct mailing and newsletters
* Instant or scheduled mailing
* Storing templates
* Template management (creating, upgrading and deleting)

## Prerequisites
* Kubernetes

## Installation

### Using Helm

```
  helm repo add containerum https://charts.containerum.io
  helm repo update
  helm install containerum/mail
```

## Contributions
Please submit all contributions concerning Mail component to this repository. Contributing guidelines are available [here](https://github.com/containerum/containerum/blob/master/CONTRIBUTING.md).

## License
Mail project is licensed under the terms of the Apache License Version 2.0. Please see LICENSE in this repository for more details.
