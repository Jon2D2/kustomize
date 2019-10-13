# Feature Test for Issue 1251


This folder contains files describing how to address [Issue 1251](https://github.com/kubernetes-sigs/kustomize/issues/1251)

## Setup the workspace

First, define a place to work:

<!-- @makeWorkplace @test -->
```bash
DEMO_HOME=$(mktemp -d)
```

## Preparation

<!-- @makeDirectories @test -->
```bash
mkdir -p ${DEMO_HOME}/
mkdir -p ${DEMO_HOME}/inheritance
mkdir -p ${DEMO_HOME}/inheritance/base
mkdir -p ${DEMO_HOME}/inheritance/composite
mkdir -p ${DEMO_HOME}/inheritance/dns
mkdir -p ${DEMO_HOME}/inheritance/probe
mkdir -p ${DEMO_HOME}/inheritance/restart
mkdir -p ${DEMO_HOME}/references
mkdir -p ${DEMO_HOME}/references/base
mkdir -p ${DEMO_HOME}/references/composite
mkdir -p ${DEMO_HOME}/references/dns
mkdir -p ${DEMO_HOME}/references/probe
mkdir -p ${DEMO_HOME}/references/restart
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/inheritance/base/kustomization.yaml
resources:
- deployment.yaml
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/inheritance/composite/kustomization.yaml
resources:
- ../probe
- ../dns
- ../restart
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/inheritance/dns/kustomization.yaml
resources:
- ../base

patchesStrategicMerge:
- dep-patch.yaml
EOF
```


### Preparation Step KustomizationFile3

<!-- @createKustomizationFile3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/inheritance/probe/kustomization.yaml
resources:
- ../base

patchesStrategicMerge:
- dep-patch.yaml
EOF
```


### Preparation Step KustomizationFile4

<!-- @createKustomizationFile4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/inheritance/restart/kustomization.yaml
resources:
- ../base

patchesStrategicMerge:
- dep-patch.yaml
EOF
```


### Preparation Step KustomizationFile5

<!-- @createKustomizationFile5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/references/base/kustomization.yaml
resources:
- deployment.yaml

# configurations:
# - kustomizeconfig.yaml
EOF
```


### Preparation Step KustomizationFile6

<!-- @createKustomizationFile6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/references/composite/kustomization.yaml
resources:
- ../base
- ../probe
- ../dns
- ../restart
EOF
```


### Preparation Step KustomizationFile7

<!-- @createKustomizationFile7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/references/dns/kustomization.yaml
resources:
- ../base
- dep-patch.yaml

# vars:
# - name: Deployment.dns.spec.template.spec.dnsPolicy
#   objref:
#     kind: Deployment
#     name: dns
#     apiVersion: apps/v1
#   fieldref:
#     fieldpath: spec.template.spec.dnsPolicy
EOF
```


### Preparation Step KustomizationFile8

<!-- @createKustomizationFile8 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/references/probe/kustomization.yaml
resources:
- ../base
- dep-patch.yaml

# vars: 
# - name: Deployment.probe.spec.template.spec.containers[0].livenessProbe
#   objref:
#     kind: Deployment
#     name: probe
#     apiVersion: apps/v1
#   fieldref:
#     fieldpath: spec.template.spec.containers[0].livenessProbe
EOF
```


### Preparation Step KustomizationFile9

<!-- @createKustomizationFile9 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/references/restart/kustomization.yaml
resources:
- ../base
- dep-patch.yaml

# vars:
# - name: Deployment.restart.spec.template.spec.restartPolicy
#   objref:
#     kind: Deployment
#     name: restart
#     apiVersion: apps/v1
#   fieldref:
#     fieldpath: spec.template.spec.restartPolicy
EOF
```


### Preparation Step KustomizeConfig0

<!-- @createKustomizeConfig0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/references/base/kustomizeconfig.yaml
varReference:
- path: spec/template/spec/containers[]/livenessProbe
  kind: Deployment
- path: spec/template/spec/dnsPolicy
  kind: Deployment
- path: spec/template/spec/restartPolicy
  kind: Deployment
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/inheritance/base/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  template:
    spec:
      containers:
      - name: my-deployment
        image: my-image
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/inheritance/dns/dep-patch.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  template:
    spec:
      dnsPolicy: ClusterFirst
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/inheritance/dns/values.yaml
EOF
```


### Preparation Step Resource3

<!-- @createResource3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/inheritance/probe/dep-patch.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  template:
    spec:
      containers:
      - name: my-deployment
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8080
EOF
```


### Preparation Step Resource4

<!-- @createResource4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/inheritance/probe/values.yaml
EOF
```


### Preparation Step Resource5

<!-- @createResource5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/inheritance/restart/dep-patch.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  template:
    spec:
      restartPolicy: Always
EOF
```


### Preparation Step Resource6

<!-- @createResource6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/references/base/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  template:
    spec:
      containers:
      - name: my-deployment
        image: my-image
        livenessProbe: $(Deployment.probe.spec.template.spec.containers[0].livenessProbe)
      dnsPolicy: $(Deployment.dns.spec.template.spec.dnsPolicy)
      restartPolicy: $(Deployment.restart.spec.template.spec.restartPolicy)
EOF
```


### Preparation Step Resource7

<!-- @createResource7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/references/dns/dep-patch.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: dns
  namespace: patch
spec:
  template:
    spec:
      dnsPolicy: ClusterFirst
EOF
```


### Preparation Step Resource8

<!-- @createResource8 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/references/probe/dep-patch.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: probe
  namespace: patch
spec:
  template:
    spec:
      containers:
      - livenessProbe:
          httpGet:
            path: /healthz
            port: 8080
EOF
```


### Preparation Step Resource9

<!-- @createResource9 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/references/restart/dep-patch.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: restart
  namespace: patch
spec:
  template:
    spec:
      restartPolicy: Always
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir -p ${DEMO_HOME}/actual
mkdir -p ${DEMO_HOME}/actual/inheritance
mkdir -p ${DEMO_HOME}/actual/references
kustomize build ${DEMO_HOME}/inheritance/composite -o ${DEMO_HOME}/actual/inheritance
kustomize build ${DEMO_HOME}/references/composite -o ${DEMO_HOME}/actual/references
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir -p ${DEMO_HOME}/expected
mkdir -p ${DEMO_HOME}/expected/inheritance
mkdir -p ${DEMO_HOME}/expected/references
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/inheritance/apps_v1_deployment_my-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  template:
    spec:
      containers:
      - image: my-image
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8080
        name: my-deployment
      dnsPolicy: ClusterFirst
      restartPolicy: Always
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/references/default_apps_v1_deployment_my-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  template:
    spec:
      containers:
      - image: my-image
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8080
        name: my-deployment
      dnsPolicy: ClusterFirst
      restartPolicy: Always
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/references/patch_apps_v1_deployment_dns.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: dns
  namespace: patch
spec:
  template:
    spec:
      dnsPolicy: ClusterFirst
EOF
```


### Verification Step Expected3

<!-- @createExpected3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/references/patch_apps_v1_deployment_probe.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: probe
  namespace: patch
spec:
  template:
    spec:
      containers:
      - livenessProbe:
          httpGet:
            path: /healthz
            port: 8080
EOF
```


### Verification Step Expected4

<!-- @createExpected4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/references/patch_apps_v1_deployment_restart.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: restart
  namespace: patch
spec:
  template:
    spec:
      restartPolicy: Always
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

