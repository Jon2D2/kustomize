# Demo: Using non-string values accross objects

Kustomize is able to use any value from a given kubernetes object as a
variable. The following example demonstrates how to keep the number of replicas
for the `bar` deployment synchronized with that of the `foo` deployment.

First, define a place to work:

<!-- @makeWorkplace @test -->
```
DEMO_HOME=$(mktemp -d)
```

## Create the deployments

First, create the `foo` deployment. We'll set it to use 5 replicas.

<!-- @createFoo @test -->
```
mkdir -p $DEMO_HOME/foo
cat <<'EOF' >$DEMO_HOME/foo/deployment.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: foo
  labels:
    app: foo
spec:
  replicas: 5
  selector:
    matchLabels:
      app: foo
  template:
    metadata:
      labels:
        app: foo
    spec:
      containers:
      - image: alpine
        name: foo
EOF

cat <<'EOF' >$DEMO_HOME/foo/kustomization.yaml
resources:
- deployment.yaml
EOF
```

Next, create the `bar` deployment. Note that this deployment will use a
variable rather than hard-coding the number of replicas.

<!-- @createBar @test -->
```
mkdir -p $DEMO_HOME/bar
cat <<'EOF' >$DEMO_HOME/bar/deployment.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: bar
  labels:
    app: bar
spec:
  replicas: $(REPLICAS)
  selector:
    matchLabels:
      app: bar
  template:
    metadata:
      labels:
        app: bar
    spec:
      containers:
      - image: alpine
        name: bar
EOF

cat <<'EOF' >$DEMO_HOME/bar/kustomization.yaml
resources:
- deployment.yaml
EOF
```

## Create the variable references

Now we just need to hook up the variable to the field it refers to.

<!-- @createVar @test -->
```
cat <<'EOF' >$DEMO_HOME/kustomization.yaml
bases:
  - foo
  - bar
namePrefix: generics-example-
vars:
  - name: REPLICAS
    objref:
      kind: Deployment
      name: foo
      apiVersion: apps/v1beta2
    fieldref:
      fieldpath: spec.replicas
configurations:
  - transformer.yaml
EOF

cat <<'EOF' >$DEMO_HOME/transformer.yaml
varReference:
- kind: Deployment
  path: spec/replicas
EOF
```

## Build and inspect the results

Build both deployments with the following. Results will be emitted to stdout.
Note that the `bar` deployment has its `replicas` field set to 5.

<!-- @build @test -->
```
kustomize build $DEMO_HOME
```
