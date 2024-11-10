<h1 align="center">
  Katalog Agent
</h1>
<h2 align="center">
  an application to scrape deployment information from a kubernetes cluster and forward it via gRPC to a central collection server
</h2>

<div align="center">

[![Made With Go][made-with-go-badge]][for-the-badge-link] [![Made With gRPC][made-with-grpc-badge]][for-the-badge-link]
</div>

---
## Development
### Prerequisites
* It is strongly recommended that you run `make` which sets up pre-commit hooks and installs the recommended tools needed for working in this repository (except `kind`).
* Kind
  * https://kind.sigs.k8s.io/docs/user/quick-start#installation
### Building
* Please use `kind` in combination with the `make` targets to get the full experience
  * `docker compose` is supported, however will not provide a good experience because there's no kubernetes api to interface with

1. `make kind`
   1. Sets up Kind Cluster with one control-plane and one worker node
2. `make kind-load`
   1. Builds image with Docker and loads it onto the cluster
3. `make apply`
   1. Deploys the Kubernetes files in `k8s/` to the cluster
4. Profit!

<!--

Reference Variables

-->

<!-- Badges -->
[made-with-go-badge]: .github/images/made-with-go.svg
[made-with-grpc-badge]: .github/images/made-with-grpc.svg

<!-- Links -->
[blank-reference-link]: #
[for-the-badge-link]: https://forthebadge.com