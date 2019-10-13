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
mkdir -p ${DEMO_HOME}/environment
mkdir -p ${DEMO_HOME}/projects
mkdir -p ${DEMO_HOME}/projects/bar
mkdir -p ${DEMO_HOME}/projects/bar/manifests
mkdir -p ${DEMO_HOME}/projects/foo
mkdir -p ${DEMO_HOME}/projects/foo/manifests
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/environment/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
configMapGenerator:
  - name: environment
    files:
      - name
      - domain
      - branch
# vars:
#   - name: ENV
#     objref:
#       apiVersion: v1
#       kind: ConfigMap
#       name: environment
#     fieldref:
#       fieldpath: data.name
#   - name: DOMAIN
#     objref:
#       apiVersion: v1
#       kind: ConfigMap
#       name: environment
#     fieldref:
#       fieldpath: data.domain
#   - name: BRANCH
#     objref:
#       apiVersion: v1
#       kind: ConfigMap
#       name: environment
#     fieldref:
#       fieldpath: data.branch
generatorOptions:
 disableNameSuffixHash: true
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomization.yaml
resources:
  - environment
  - projects/foo
  - projects/bar
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/projects/bar/kustomization.yaml
namespace: bar
resources:
  # - ../../environment
  - manifests/ingress.yaml
EOF
```


### Preparation Step KustomizationFile3

<!-- @createKustomizationFile3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/projects/foo/kustomization.yaml
namespace: foo
resources:
  # - ../../environment
  - manifests/ingress.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/projects/bar/manifests/ingress.yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: bar
spec:
  rules:
    - host: bar$(ConfigMap.environment.data.branch).$(ConfigMap.environment.data.name).$(ConfigMap.environment.data.domain)
      http:
        paths:
        - backend:
            serviceName: bar
            servicePort: http
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/projects/foo/manifests/ingress.yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: foo
spec:
  rules:
    - host: foo$(ConfigMap.environment.data.branch).$(ConfigMap.environment.data.name).$(ConfigMap.environment.data.domain)
      http:
        paths:
        - backend:
            serviceName: foo
            servicePort: http
EOF
```


### Preparation Step Other0

<!-- @createOther0 @test -->
```bash
printf '%s' '-branch' >${DEMO_HOME}/environment/branch
```


### Preparation Step Other1

<!-- @createOther1 @test -->
```bash
printf '%s' 'domain.com' > ${DEMO_HOME}/environment/domain
```


### Preparation Step Other2

<!-- @createOther2 @test -->
```bash
printf '%s' 'dev' >${DEMO_HOME}/environment/name
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
cat <<'EOF' >${DEMO_HOME}/expected/bar_networking.k8s.io_v1beta1_ingress_bar.yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: bar
  namespace: bar
spec:
  rules:
  - host: bar-branch.dev.domain.com
    http:
      paths:
      - backend:
          serviceName: bar
          servicePort: http
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/default_~g_v1_configmap_environment.yaml
apiVersion: v1
data:
  branch: -branch
  domain: domain.com
  name: dev
kind: ConfigMap
metadata:
  name: environment
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/foo_networking.k8s.io_v1beta1_ingress_foo.yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: foo
  namespace: foo
spec:
  rules:
  - host: foo-branch.dev.domain.com
    http:
      paths:
      - backend:
          serviceName: foo
          servicePort: http
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

