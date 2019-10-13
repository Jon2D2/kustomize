# Feature Test for Issue 1251_f


This folder contains files describing how to address [Issue 1251_f](https://github.com/kubernetes-sigs/kustomize/issues/1251_f)

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

resources:
- values.yaml
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomization.yaml
resources:
- ./environment
- ./projects/foo
- ./projects/bar
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/projects/bar/kustomization.yaml
namespace: bar
resources:
  - manifests/ingress.yaml
EOF
```


### Preparation Step KustomizationFile3

<!-- @createKustomizationFile3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/projects/foo/kustomization.yaml
namespace: foo
resources:
  - manifests/ingress.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/environment/values.yaml
apiVersion: kustomize.config.k8s.io/v1
kind: Values
metadata:
  name: shared
spec:
  env: dev
  domain: domain.com
  branch: -branch
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/projects/bar/manifests/ingress.yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: bar
spec:
  rules:
    - host: bar$(Values.shared.spec.branch).$(Values.shared.spec.env).$(Values.shared.spec.domain)
      http:
        paths:
        - backend:
            serviceName: bar
            servicePort: http
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/projects/foo/manifests/ingress.yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: foo
spec:
  rules:
    - host: foo$(Values.shared.spec.branch).$(Values.shared.spec.env).$(Values.shared.spec.domain)
      http:
        paths:
        - backend:
            serviceName: foo
            servicePort: http
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME} -o ${DEMO_HOME}/actual
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
cat <<'EOF' >${DEMO_HOME}/expected/default_kustomize.config.k8s.io_v1_values_shared.yaml
apiVersion: kustomize.config.k8s.io/v1
kind: Values
metadata:
  name: shared
spec:
  branch: -branch
  domain: domain.com
  env: dev
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

