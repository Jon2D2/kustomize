# Feature Test for Issue 1478


This folder contains files describing how to address [Issue 1478](https://github.com/kubernetes-sigs/kustomize/issues/1478)

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
mkdir -p ${DEMO_HOME}/app
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/app/kustomization.yaml
kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1

nameSuffix: -dev

resources:
- ingress.yaml
- service.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/app/ingress.yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: web-private
spec:
  rules:
  - http:
      paths:
      - backend:
          serviceName: my-service-name
          servicePort: 3333
        path: /*
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/app/service.yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: service
  name: my-service-name
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/app -o ${DEMO_HOME}/actual
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/networking.k8s.io_v1beta1_ingress_web-private-dev.yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: web-private-dev
spec:
  rules:
  - http:
      paths:
      - backend:
          serviceName: my-service-name-dev
          servicePort: 3333
        path: /*
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_service_my-service-name-dev.yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: service
  name: my-service-name-dev
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

