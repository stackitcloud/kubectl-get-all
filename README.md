# kubectl-get-all
[![Go Report Card](https://goreportcard.com/badge/stackitcloud/kubectl-get-all)](https://goreportcard.com/report/stackitcloud/kubectl-get-all)
[![LICENSE](https://img.shields.io/github/license/stackitcloud/kubectl-get-all.svg)](https://github.com/stackitcloud/kubectl-get-all/blob/main/LICENSE)
[![Releases](https://img.shields.io/github/release-pre/stackitcloud/kubectl-get-all.svg)](https://github.com/stackitcloud/kubectl-get-all/releases)

Kubectl plugin to show **really** all kubernetes resources

This is a fork of [ketall](https://github.com/corneliusweig/ketall)

## Intro
For a complete overview of all resources in a kubernetes cluster, `kubectl get all --all-namespaces` is not enough, because it simply does not show everything.
This helper lists **really** all resources the cluster has to offer.

## Demo
![kubectl-get-all demo](docs/demo.gif "kubectl-get-all demo")

## Examples
Get all resources...
- ... excluding events (this is hardly ever useful)
  ```bash
  kubectl-get-all
  ```

- ... including events
  ```bash
  kubectl-get-all --exclude=
  ```

- ... created in the last minute
  ```bash
  kubectl-get-all --since 1m
  ```
  This flag understands typical human-readable durations such as `1m` or `1y1d1h1m1s`.

- ... in the default namespace
  ```bash
  kubectl-get-all --namespace=default
  ```

- ... at cluster level
  ```bash
  kubectl-get-all --only-scope=cluster
  ```

- ... using list of cached server resources
  ```bash
  kubectl-get-all --use-cache
  ```
  Note that this may fail to show **really** everything, if the http cache is stale.

- ... and combine with common `kubectl` options
  ```bash
  KUBECONFIG=otherconfig kubectl-get-all -o name --context some --namespace kube-system --selector run=skaffold
  ```

Also see [Usage](docs/USAGE.md).

## Installation
There are several ways to install `kubectl-get-all`. The recommended installation method is via `krew`.

### Via krew
Krew is a `kubectl` plugin manager. If you have not yet installed `krew`, check out the [krew install guide.](https://krew.sigs.k8s.io/docs/user-guide/setup/install/)

Then installation is as simple as
```bash
kubectl krew install get-all
```
The plugin will be available as `kubectl get-all`.

### Binaries

Download the desired version for your operating system and processor architecture from the [`kubectl-get-all` releases page](https://github.com/stackitcloud/kubectl-get-all/releases).
Make the file executable and place it in a directory available in your `$PATH`.

### From source

#### Build on host

Requirements:
 - go
 - git

```bash
go build
```
