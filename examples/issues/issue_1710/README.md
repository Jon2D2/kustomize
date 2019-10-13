# Feature Test for Issue 1710


This folder contains files describing how to address [Issue 1710](https://github.com/kubernetes-sigs/kustomize/issues/1710)

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
mkdir -p ${DEMO_HOME}/kustomizeconfig
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomization.yaml
resources:
- prometheus.yaml
- catalog-source.yaml
- service-monitor.yaml

namespace: kustomizedns

configurations:
- kustomizeconfig/prometheus.yaml
- kustomizeconfig/servicemonitor.yaml

patchesStrategicMerge:
- set-catalog-source-image.yaml
- alert-ns-patch.yaml

patchesJson6902:
- target:
    group: monitoring.coreos.com
    version: v1
    kind: ServiceMonitor
    name: prometheus-service-monitor
  path: set-svc-mon-ns.yaml

vars:
- name: NS
  objref:
    apiVersion: operators.coreos.com/v1alpha1
    kind: CatalogSource
    name: foocorp-operators-registry
  fieldref:
    fieldpath: metadata.namespace
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/alert-ns-patch.yaml
apiVersion: monitoring.coreos.com/v1
kind: Prometheus
metadata:
  name: prometheus
spec:
  alerting:
    alertmanagers:
    - name: alertmanager
      namespace: $(NS)
      port: alertmanager
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/catalog-source.yaml
apiVersion: operators.coreos.com/v1alpha1
kind: CatalogSource
metadata:
  name: foocorp-operators-registry
spec:
  displayName: Foo
  publisher: Foo Corp.
  sourceType: grpc
  image: docker-registry.default.svc:/foocorp/foocorp-operators-registry:0.0.1
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomizeconfig/prometheus.yaml
varReference:
- path: spec/alerting/alertmanagers/namespace
  kind: Prometheus
EOF
```


### Preparation Step Resource3

<!-- @createResource3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomizeconfig/servicemonitor.yaml
varReference:
- path: spec/selector/namespaceSelector/matchNames
  kind: ServiceMonitor
EOF
```


### Preparation Step Resource4

<!-- @createResource4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/prometheus.yaml
apiVersion: monitoring.coreos.com/v1
kind: Prometheus
metadata:
  name: prometheus
  labels:
    prometheus: k8s
spec:
  replicas: 2
  serviceAccountName: prometheus-k8s
  securityContext: {}
  serviceMonitorSelector: {}
  ruleSelector: {}
  alerting:
    alertmanagers:
    - namespace: IWANTTOREPLACETHIS
      name: alertmanager
      port: alertmanager
EOF
```


### Preparation Step Resource5

<!-- @createResource5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/service-monitor.yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: prometheus-service-monitor
spec:
  selector:
    namespaceSelector:
      matchNames:
      - IWANTTOREPLACETHIS
  endpoints:
  - port: metrics
    honorlablels: true
    interval: 10s
    scrapeTimeout: 10s
    etc: etc
EOF
```


### Preparation Step Resource6

<!-- @createResource6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/set-catalog-source-image.yaml
apiVersion: operators.coreos.com/v1alpha1
kind: CatalogSource
metadata:
  name: foocorp-operators-registry
spec:
  image: docker-registry.default.svc:/foocorp/foocorp-operators-registry@sha256:NOTIMPORTANTFORTHETEST
EOF
```


### Preparation Step Resource7

<!-- @createResource7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/set-svc-mon-ns.yaml
- op: replace
  path: /spec/selector/namespaceSelector/matchNames/0
  value: $(NS)
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build $DEMO_HOME -o ${DEMO_HOME}/actual
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/monitoring.coreos.com_v1_prometheus_prometheus.yaml
apiVersion: monitoring.coreos.com/v1
kind: Prometheus
metadata:
  labels:
    prometheus: k8s
  name: prometheus
  namespace: kustomizedns
spec:
  alerting:
    alertmanagers:
    - name: alertmanager
      namespace: kustomizedns
      port: alertmanager
  replicas: 2
  ruleSelector: {}
  securityContext: {}
  serviceAccountName: prometheus-k8s
  serviceMonitorSelector: {}
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/monitoring.coreos.com_v1_servicemonitor_prometheus-service-monitor.yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: prometheus-service-monitor
  namespace: kustomizedns
spec:
  endpoints:
  - etc: etc
    honorlablels: true
    interval: 10s
    port: metrics
    scrapeTimeout: 10s
  selector:
    namespaceSelector:
      matchNames:
      - kustomizedns
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/operators.coreos.com_v1alpha1_catalogsource_foocorp-operators-registry.yaml
apiVersion: operators.coreos.com/v1alpha1
kind: CatalogSource
metadata:
  name: foocorp-operators-registry
  namespace: kustomizedns
spec:
  displayName: Foo
  image: docker-registry.default.svc:/foocorp/foocorp-operators-registry@sha256:NOTIMPORTANTFORTHETEST
  publisher: Foo Corp.
  sourceType: grpc
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

