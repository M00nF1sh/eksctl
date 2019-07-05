---
title: "Project Roadmap"
weight: 20
url: community/roadmap
---

## Project Roadmap Themes

### GitOps: Declarative Configuration Managent (aka Cluster API)

One should be able to make EKS cluster configuration through declarative config files (`eksctl apply`). Additionally, they should be able to manage a cluster via a git repo.

### Cluster Applications (aka Add-ons)

It should be easy to create a cluster with various applications pre-installed, e.g. Weave Flux, Helm 2 (Tiller), ALB Ingress controller, to name a few.

> Note: this will depend on how [add-ons spec by SIG Cluster Lifecycle](https://github.com/kubernetes/enhancements/pull/746) will evolve.
