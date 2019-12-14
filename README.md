# kustomize++

This is a fork of [kustomize](https://github.com/keleustes/kustomize/tree/allinone/examples/issues) based on the api/v0.3.0 (20191212)

The kustomization.yaml syntax is identical to the upstream repo.

The allinone branch of this fork supports brings 5 main additional features.

Those features were build overtime to fullfill the needs of users
struggling to use kustomize at scale.
Very often the issue has been reproduced and a feature test added
to validate that there was a simple solution available:
[Feature Test](https://github.com/keleustes/kustomize/tree/allinone/examples/issues)

The intend is to remove the fork once those features are
available (even another form) as long as it is manageable by the
user and he does not end up with managing 1000 lines kustomization.yaml 

## Feature 1: Support for CRD Strategic Merge Patch

CRDs are now everywhere in kubernetes based projects..OpenShift, Istio, Argo, ...

StrategicMergePatch for CRDs is supported (as opposed to JMP and JsonPatch) 
like Kubernetes 1.16 does as long you can register your Scheme in 
the kustomize global Scheme.

[Registration](https://github.com/keleustes/kustomize/blob/airshipctl/kustomize/register/RolloutCRDRegister.go#L20)
[Examples of usage](https://github.com/keleustes/kustomize/tree/airshipctl/examples/crds)

The upstream version of kustomize is silently using different merge strategy
depending on the kind of your resource.

Moreover the current attempt to replace the kubernetes code by kyaml currently
prevent any consumer of the kustomize API library from addressing the issue.

## Feature 2: Support for Default Transformers Configuration adjustments

Let's you adjust default transformers behavior. For instance,
You can add, update, or remove/skip transformations performed
by the commonLabels transformer without having to copy/paste the entire
default configuration of the transformer.

## Feature 3: Let you control the order of the K8s resources in the output

If you can sort and filter your resources at build time. THis allows you
to add your ClusterWide CRD at the right place.

To ease the creation:

```bash
kustomize build <xxx> --reorder=kubectlapply | kubectl apply -f -
```

To ease the deletion:

```bash
kustomize build <xxx> --reorder=kubectldelete | kubectl delete -f -
```

## Feature 4: Let you refer/inline field from on K8s Resource to another.

You can refer the field value from one Kubernetes in another in one simple line.
No need for error prone var and varReference declaration, 
nor to create an unwanted ConfigMap just to be able to share a value accross Resources.
This is also verify useful to need end up doing something like

```bash
kustomize build <xxx> | envsubst
```
```bash
kustomize build <xxx> | <somepythonscript.py
```

For instance: 

```yaml
HorizontalPodScaler:
metadata:
  name: my-hpa
spec:
  minReplicas: $(Deployment.my-deployment.spec.replicas) 
```

## Feature 5: Diamond based import and Composition.

Also not perfect yet, this feature lets you create common resources that you
can import in the different components of your application (for instance log, db, prometheus, ingres) before
composing your application from those components.

