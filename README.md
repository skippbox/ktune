# ktune
Kubernetes Applications Auto-tuning

For example, ktune watches image repository for new images, and triggers canary deployments in k8s.
It makes sure you are always running the latest images in your deployments.

## Build

You need glide

You will first need a working go build environment set up. [See here for advice](https://golang.org/doc/install)

Clone the repo and ```make```

```
$ git clone https://github.com/skippbox/ktune
$ glide install
$ make build
```

## Example

Create a deployment that you consider _stable_ and give it a label to identify it, for example in a single command:

```
$ kubectl run hostname --image=runseb/hostname:0.1.0 --labels="run=hostname,track=stable"
```

Start `ktune`, it will check the deployment and realize that there is an image on the Docker hub with the tag _latest_, it will generate a canary deployment.

```
$ ./ktune --debug
INFO[0000] Starting ktune
INFO[0000] Worker iteration at 2016-10-17 15:19:58.655875145 +0200 CEST
INFO[0010] Checking Deployment 'hostname'
INFO[0011] Latest tag found to be ''
INFO[0011] image '[runseb/hostname 0.1.0]' is outdated. New canary will update to 'runseb/hostname:latest'
INFO[0011] Creating new deployment 'hostname-ktune'
```

Check your deployments and resulting pods:

```
$ kubectl get deployments --show-labels
NAME             DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE       LABELS
hostname         1         1         1            1           26m       run=hostname,track=stable
hostname-ktune   1         1         1            1           5s        run=hostname,track=canary
$ kubectl get pods
NAME                             READY     STATUS    RESTARTS   AGE
hostname-3656820496-fztby        1/1       Running   0          31m
hostname-ktune-969058483-k16wj   1/1       Running   0          4m
```

## WIP

For now, ktune will:
- Reads kubeconfig file and checks RCs images
- Looks for images on docker hub (other repos not yet supported)
- Uses basic tag comparison to detect whether the tag being used is the latest
- Read images used by Deployments
- Creates canary deployment

Still doesn't ...
- Use patterns to indicate which images should be updated
- Use patterns to indicate which tags are latest
- Support gcr.io or quay.io

## Usage

Use ```./ktune -h``` to get usage

ktune expect a container version tag that can be machine compared, following the pattern
Major.Minor.Build-Revision

To get the latest image version, grisou:

1. Divides the tag string in dots
2. Keep the first 3 groups, discarding the rest
3. In each group get the first number, dicarding letters or symbols 
4. In the last group, looks for a minus symbol, and tries to read a number to the right

Examples:
- ```1.2```
- ```v1.2```
- ```1.2.3```
- ```1.2.3-4```
- ```v1.b.2.g3-4``` equals - ```1.2.3-4```
- ```latest``` is considered latest version.
 

