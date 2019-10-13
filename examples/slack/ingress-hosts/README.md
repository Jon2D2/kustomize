# Feature Test for Issue ingress-hosts


This folder contains files describing how to address [Issue ISSUENUMBER](https://github.com/kubernetes-sigs/kustomize/issues/ISSUENUMBER)

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
mkdir -p ${DEMO_HOME}/base
mkdir -p ${DEMO_HOME}/overlays
mkdir -p ${DEMO_HOME}/overlays/production
mkdir -p ${DEMO_HOME}/overlays/staging
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomization.yaml
resources:
- ./ingress.yaml
- ./service.yaml
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/production/kustomization.yaml
resources:
- ../../base
namePrefix: production-
EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/staging/kustomization.yaml
resources:
- ../../base
namePrefix: staging-
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/ingress.yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: wordpress
  annotations:
    kubernetes.io/ingress.class: nginx
spec:
  tls:
  - hosts:
    - example.com
    - $(Ingress.wordpress.spec.rules[0].http.paths[0].backend.serviceName)
    secretName: wordpress-cert
  rules:
  - host: example.com
    http:
      paths:
      - backend:
          serviceName: wordpress
          servicePort: 80
        path: /
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: wordpress
  labels:
    app: wordpress
spec:
  type: NodePort
  selector:
    app: wordpress
  ports:
    - port: 80
      targetPort: 3000

EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/overlays/staging -o ${DEMO_HOME}/actual/staging.yaml
kustomize build ${DEMO_HOME}/overlays/production -o ${DEMO_HOME}/actual/production.yaml
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/production.yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: wordpress
  name: production-wordpress
spec:
  ports:
  - port: 80
    targetPort: 3000
  selector:
    app: wordpress
  type: NodePort
---
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.class: nginx
  name: production-wordpress
spec:
  rules:
  - host: example.com
    http:
      paths:
      - backend:
          serviceName: production-wordpress
          servicePort: 80
        path: /
  tls:
  - hosts:
    - example.com
    - production-wordpress
    secretName: wordpress-cert
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/staging.yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: wordpress
  name: staging-wordpress
spec:
  ports:
  - port: 80
    targetPort: 3000
  selector:
    app: wordpress
  type: NodePort
---
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.class: nginx
  name: staging-wordpress
spec:
  rules:
  - host: example.com
    http:
      paths:
      - backend:
          serviceName: staging-wordpress
          servicePort: 80
        path: /
  tls:
  - hosts:
    - example.com
    - staging-wordpress
    secretName: wordpress-cert
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

