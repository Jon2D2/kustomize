# Demo: allinone

In this tutorial, you will learn how to:
- use the kustomize layering
- use variables (simple as well as complexes)
- use inlines (entire subtrees).
- add support for crd

## base

This tutorial is inspired from the wordpress tutorial with
and brings together additional features.

### base/kustomizeconfig content

The fields that are supporting variable replacement and inline
are described in the base/kustomizeconfig/*.yaml files

### base/crds content

Create simple CRDs into your K8s cluster. It is usefull to check
the CRD built by kustomize are syntaxically correct.

```bash
kubectl apply -f base/crds
```

### base/wordpress content

Service and Deployment for wordpress. Pay attentiion to the variable
and inlines

### base/mysql content

Service and Deployment for wordpress. Pay attentiion to the variable
and inlines

## dev environment

This leverage the patching and merging capabilities of kustomize.
The versions.yaml acts as a catalogue of software used in the development
environments.

```bash
kustomize build dev | kubectl apply -f -
```

## production environments

#### common content

This leverage the patching and merging capabilities of kustomize.
The versions.yaml acts as a catalogue of software used in the production
environments.

#### prod1 environment

```bash
kustomize build production/prod1 | kubectl apply -f -
```

#### prod2 environment

```bash
kustomize build production/prod2 | kubectl apply -f -
```
