# Feature Test for Issue 1685


This folder contains files describing how to address [Issue 1685](https://github.com/kubernetes-sigs/kustomize/issues/1685)

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
# kustomization.yaml
resources:
- issuer.yaml
- ingress.yaml
namePrefix: staging-
configurations:
- issuer-kustomize-config.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/ingress.yaml
# ingress.yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: wordpress
  annotations:
    kubernetes.io/ingress.class: nginx
    certmanager.k8s.io/issuer: letsencrypt-staging # <--- Reference to Issuer
spec:
  tls:
  - hosts:
    - example.com
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
cat <<'EOF' >${DEMO_HOME}/issuer-kustomize-config.yaml
# issuer-kustomize-config.yaml
nameReference:
- kind: Issuer
  group: certmanager.k8s.io
  fieldSpecs:
    - kind: Ingress
      path: metadata/annotations/certmanager.k8s.io\/issuer

EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/issuer.yaml
# issuer.yaml
apiVersion: certmanager.k8s.io/v1alpha1
kind: Issuer
metadata:
  name: letsencrypt-staging
spec:
  acme:
    server: https://acme-staging-v02.api.letsencrypt.org/directory
    email: admin@example.com
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
cat <<'EOF' >${DEMO_HOME}/expected/certmanager.k8s.io_v1alpha1_issuer_staging-letsencrypt-staging.yaml
apiVersion: certmanager.k8s.io/v1alpha1
kind: Issuer
metadata:
  name: staging-letsencrypt-staging
spec:
  acme:
    email: admin@example.com
    server: https://acme-staging-v02.api.letsencrypt.org/directory
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/networking.k8s.io_v1beta1_ingress_staging-wordpress.yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  annotations:
    certmanager.k8s.io/issuer: staging-letsencrypt-staging
    kubernetes.io/ingress.class: nginx
  name: staging-wordpress
spec:
  rules:
  - host: example.com
    http:
      paths:
      - backend:
          serviceName: wordpress
          servicePort: 80
        path: /
  tls:
  - hosts:
    - example.com
    secretName: wordpress-cert
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

