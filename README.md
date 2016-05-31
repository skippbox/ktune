# Grisou
Kubernetes canary deployments based on image changes

## WIP

For now
- Reads kubeconfig file and checks RCs images
- Looks for images on docker hub (not other repos supported)
- Uses basic tag comparison to detect whether the tag being used is the latest


Still doesn't ...
- Reads images used by Replica Sets / Deployments
- Automatically update images
- Use patterns to indicate which images should be updated
- Use patterns to indicate which tags are latest
