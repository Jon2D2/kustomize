# Feature Test for Issue 0964


This folder contains files describing how to address [Issue 0964](https://github.com/kubernetes-sigs/kustomize/issues/0964)

## Setup the workspace

First, define a place to work:

<!-- @makeWorkplace @test -->
```bash
DEMO_HOME=$(mktemp -d)
```

## Preparation

<!-- @makeDirectories @test -->
```bash
mkdir -p ${DEMO_HOME}
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomization.yaml
---
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1


resources:
- resource.yaml
- values.yaml

vars:
- name: Values.file1.spec.Release.Namespace
  objref:
    apiVersion: kustomize.config.k8s.io/v1
    kind: Values
    name: file1
  fieldref:
    fieldpath: spec.Release.Namespace
- name: Values.file1.spec.url
  objref:
    apiVersion: kustomize.config.k8s.io/v1
    kind: Values
    name: file1
  fieldref:
    fieldpath: spec.url

configurations:
- kustomizeconfig.yaml
EOF
```


### Preparation Step KustomizeConfig0

<!-- @createKustomizeConfig0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomizeconfig.yaml
varReference:
- kind: Service
  path: metadata/annotations
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/resource.yaml
apiVersion: v1
kind: Service
metadata:
  name: name_of_the_service
  annotations:
    getambassador.io/config: |
      ---
      apiVersion: ambassador/v1
      kind:  Mapping
      name:  $(Service.name_of_the_service.metadata.name)
      prefix: /
      service: $(Service.name_of_the_service.metadata.name).$(Values.file1.spec.Release.Namespace)
      host: $(Values.file1.spec.url)
spec:
  type: ClusterIP
  ports:
    - port: 80
      targetPort: http
      protocol: TCP
      name: http
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/values.yaml
apiVersion: kustomize.config.k8s.io/v1
kind: Values
metadata:
  name: file1
spec:
  Release:
    Namespace: name_of_the_namespace_we_deploy_into
  url: a_configured_host
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
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_service_name_of_the_service.yaml
apiVersion: v1
kind: Service
metadata:
  annotations:
    getambassador.io/config: |
      ---
      apiVersion: ambassador/v1
      kind:  Mapping
      name:  name_of_the_service
      prefix: /
      service: name_of_the_service.name_of_the_namespace_we_deploy_into
      host: a_configured_host
  name: name_of_the_service
spec:
  ports:
  - name: http
    port: 80
    protocol: TCP
    targetPort: http
  type: ClusterIP
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/kustomize.config.k8s.io_v1_values_file1.yaml
apiVersion: kustomize.config.k8s.io/v1
kind: Values
metadata:
  name: file1
spec:
  Release:
    Namespace: name_of_the_namespace_we_deploy_into
  url: a_configured_host
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

