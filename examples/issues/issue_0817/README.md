# Feature Test for Issue 0817


This folder contains files describing how to address [Issue 0817](https://github.com/kubernetes-sigs/kustomize/issues/0817)

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
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

configurations:
- commonlabels.yaml

commonLabels:
  app.kubernetes.io/part-of: argocd

resources:
- service.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/commonlabels.yaml
commonLabels:
- path: spec/selector
  create: false   # <<<<< I wish to disable this behavior
  version: v1
  kind: Service
  behavior: remove
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: redis
spec:
  selector:
    app: redis
  ports:
  - name: server
    port: 6379
    protocol: TCP
    targetPort: redis
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
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_service_redis.yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app.kubernetes.io/part-of: argocd
  name: redis
spec:
  ports:
  - name: server
    port: 6379
    protocol: TCP
    targetPort: redis
  selector:
    app: redis
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

