# Mail service for Containerum Platform

## Installation

### Using Helm

```
  helm repo add containerum https://containerum.github.io/mail
  helm repo update
  helm install helm containerum/mail
```
By default it uses **emptyDir** storage!

To use GlusterFS run:
```
  helm install helm containerum/mail --set volume.empty="" --set volume.gluster.glusterfs.endpoints=$GLUSTER_ENDPOINT --set volume.gluster.glusterfs.path=$GLUSTER_PATH
```
