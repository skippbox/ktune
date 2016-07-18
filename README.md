# Grisou
Kubernetes canary deployments based on image changes

## WIP

For now, Grisou will:
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

You will first need a working go build environment set up. [See here for advice](https://golang.org/doc/install)

Clone the repo and ```make```

Use ```./grisou -h``` to get usage

Grisou expect a container version tag that can be machine compared, following the pattern
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
 

