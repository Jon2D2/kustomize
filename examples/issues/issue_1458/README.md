# Feature Test for Issue 1458


This folder contains files describing how to address [Issue 1458](https://github.com/kubernetes-sigs/kustomize/issues/1458)

## Setup the workspace

First, define a place to work:

<!-- @makeWorkplace @test -->
```bash
DEMO_HOME=$(mktemp -d)
```

## Preparation

<!-- @makeDirectories @test -->
```bash
mkdir -p ${DEMO_HOME}/base
mkdir -p ${DEMO_HOME}/base/common
mkdir -p ${DEMO_HOME}/base/common/kustomizeconfig
mkdir -p ${DEMO_HOME}/base/web
mkdir -p ${DEMO_HOME}/base/worker
mkdir -p ${DEMO_HOME}/overlays
mkdir -p ${DEMO_HOME}/overlays/a
mkdir -p ${DEMO_HOME}/overlays/b
mkdir -p ${DEMO_HOME}/overlays/c
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/common/kustomization.yaml
---
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

configurations:
- ./kustomizeconfig/servicemonitor.yaml

EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/web/kustomization.yaml
---
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

commonLabels:
  type: web

resources:
- ../common
- ./servicemonitor.yaml
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/worker/kustomization.yaml
---
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

commonLabels:
  type: worker

resources:
- ../common
- ./servicemonitor.yaml
EOF
```


### Preparation Step KustomizationFile3

<!-- @createKustomizationFile3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/a/kustomization.yaml
---
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

commonLabels:
  part-of: a

resources:
- ../../base/web
- ../../base/worker
EOF
```


### Preparation Step KustomizationFile4

<!-- @createKustomizationFile4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/b/kustomization.yaml
---
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

commonLabels:
  part-of: b

resources:
- ../../base/web
- ../../base/worker
EOF
```


### Preparation Step KustomizationFile5

<!-- @createKustomizationFile5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/c/kustomization.yaml
---
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

commonLabels:
  part-of: c

resources:
- ../../base/web
- ../../base/worker
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/common/kustomizeconfig/servicemonitor.yaml
commonLabels:
- path: spec/selector/matchLabels
  create: true
  kind: ServiceMonitor
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/web/servicemonitor.yaml
---
apiVersion: monitoring.coreos.com/v1alpha1
kind: ServiceMonitor
metadata:
  name: monitor-web
spec:
  endpoints:
  - port: web            # works for different port numbers as long as the name matches
    interval: 10s        # scrape the endpoint every 10 seconds
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/worker/servicemonitor.yaml
---
apiVersion: monitoring.coreos.com/v1alpha1
kind: ServiceMonitor
metadata:
  name: monitor-worker
spec:
  endpoints:
  - port: worker         # works for different port numbers as long as the name matches
    interval: 10s        # scrape the endpoint every 10 seconds
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/overlays/a -o ${DEMO_HOME}/actual/a.yaml
kustomize build ${DEMO_HOME}/overlays/b -o ${DEMO_HOME}/actual/b.yaml
kustomize build ${DEMO_HOME}/overlays/c -o ${DEMO_HOME}/actual/c.yaml
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/a.yaml
apiVersion: monitoring.coreos.com/v1alpha1
kind: ServiceMonitor
metadata:
  labels:
    part-of: a
    type: web
  name: monitor-web
spec:
  endpoints:
  - interval: 10s
    port: web
  selector:
    matchLabels:
      part-of: a
      type: web
---
apiVersion: monitoring.coreos.com/v1alpha1
kind: ServiceMonitor
metadata:
  labels:
    part-of: a
    type: worker
  name: monitor-worker
spec:
  endpoints:
  - interval: 10s
    port: worker
  selector:
    matchLabels:
      part-of: a
      type: worker
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/b.yaml
apiVersion: monitoring.coreos.com/v1alpha1
kind: ServiceMonitor
metadata:
  labels:
    part-of: b
    type: web
  name: monitor-web
spec:
  endpoints:
  - interval: 10s
    port: web
  selector:
    matchLabels:
      part-of: b
      type: web
---
apiVersion: monitoring.coreos.com/v1alpha1
kind: ServiceMonitor
metadata:
  labels:
    part-of: b
    type: worker
  name: monitor-worker
spec:
  endpoints:
  - interval: 10s
    port: worker
  selector:
    matchLabels:
      part-of: b
      type: worker
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/c.yaml
apiVersion: monitoring.coreos.com/v1alpha1
kind: ServiceMonitor
metadata:
  labels:
    part-of: c
    type: web
  name: monitor-web
spec:
  endpoints:
  - interval: 10s
    port: web
  selector:
    matchLabels:
      part-of: c
      type: web
---
apiVersion: monitoring.coreos.com/v1alpha1
kind: ServiceMonitor
metadata:
  labels:
    part-of: c
    type: worker
  name: monitor-worker
spec:
  endpoints:
  - interval: 10s
    port: worker
  selector:
    matchLabels:
      part-of: c
      type: worker
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

