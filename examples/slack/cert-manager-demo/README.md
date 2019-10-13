# Kustomize Regression Test based on cert-manager-demo

This folder is only used for kustomize regression testing.
The original files are located [here](https://github.com/jetstack/kustomize-cert-manager-demo)

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
mkdir -p ${DEMO_HOME}/overlays/cert-manager
mkdir -p ${DEMO_HOME}/overlays/development
mkdir -p ${DEMO_HOME}/overlays/multi-environment
mkdir -p ${DEMO_HOME}/overlays/production
mkdir -p ${DEMO_HOME}/overlays/staging
mkdir -p ${DEMO_HOME}/overlays/staging/secret
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

commonLabels:
  app: helloweb

namespace: helloweb

resources:
- namespace.yaml
- deployment.yaml
- service.yaml
- ingress.yaml

EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/cert-manager/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

bases:
- ../../base

commonLabels:
  app: helloweb

namespace: helloweb

resources:
- issuer.yaml
- certificate.yaml

patchesStrategicMerge:
- ingress.yaml

configurations:
- cert-manager-configuration.yaml

EOF
```


### Preparation Step KustomizationFile2

<!-- @createKustomizationFile2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/development/kustomization.yaml
bases:
- ../cert-manager

commonLabels:
  app: helloweb-development

namespace: helloweb-development

nameSuffix: -development

patchesStrategicMerge:
- ingress.yaml
- issuer.yaml
- certificate.yaml

resources:
- selfsigned-issuer.yaml
- selfsigned-certificate.yaml

EOF
```


### Preparation Step KustomizationFile3

<!-- @createKustomizationFile3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/multi-environment/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

bases:
- ../development
- ../staging
- ../production

EOF
```


### Preparation Step KustomizationFile4

<!-- @createKustomizationFile4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/production/kustomization.yaml
bases:
  - ../cert-manager

commonLabels:
  app: helloweb-production

namespace: helloweb-production

nameSuffix: -production

patchesStrategicMerge:
- ingress.yaml
- issuer.yaml
- certificate.yaml

EOF
```


### Preparation Step KustomizationFile5

<!-- @createKustomizationFile5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/staging/kustomization.yaml
bases:
- ../cert-manager

commonLabels:
  app: helloweb-staging

namespace: helloweb-staging

nameSuffix: -staging

patchesStrategicMerge:
- ingress.yaml
- issuer.yaml
- certificate.yaml

secretGenerator:
- name: ca-secret
  files:
    - secret/tls.crt
    - secret/tls.key
  type: "kubernetes.io/tls"

EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: helloweb-deployment
spec:
  selector:
    matchLabels:
      app: helloweb
  template:
    metadata:
      labels:
        app: helloweb
    spec:
      containers:
      - name: hello-app
        image: gcr.io/google-samples/hello-app:1.0
        ports:
        - containerPort: 8080

EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/ingress.yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: ingress
  annotations:
    kubernetes.io/ingress.global-static-ip-name: web-static-ip
spec:
  rules:
  - http:
      paths:
      - path: /
        backend:
          serviceName: helloweb-service
          servicePort: 8080

EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: helloweb

EOF
```


### Preparation Step Resource3

<!-- @createResource3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: helloweb-service
spec:
  type: NodePort
  selector:
    app: helloweb
  ports:
  - protocol: TCP
    port: 8080
    targetPort: 8080

EOF
```


### Preparation Step Resource4

<!-- @createResource4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/cert-manager/certificate.yaml
apiVersion: certmanager.k8s.io/v1alpha1
kind: Certificate
metadata:
  name: certificate
spec:
  secretName: secret-demo
  renewBefore: 360h # 15d
  # Change this to your own domain
  commonName: demo.example.net
  dnsNames:
  # Change this to your own domain
  - demo.example.net
  issuerRef:
    name: issuer
    kind: Issuer
  acme:
    config:
    - http01:
        ingress: ingress
      domains:
      # Change this to your own domain
      - demo.example.net

EOF
```


### Preparation Step Resource5

<!-- @createResource5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/cert-manager/cert-manager-configuration.yaml
nameReference:
- kind: Issuer
  fieldSpecs:
  - path: spec/issuerRef/name
    kind: Certificate
- kind: Secret
  fieldSpecs:
  - path: spec/ca/secretName
    kind: Issuer
- kind: Ingress
  fieldSpecs:
  - path: spec/acme/config/http01/ingress
    kind: Certificate

EOF
```


### Preparation Step Resource6

<!-- @createResource6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/cert-manager/ingress.yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: ingress
spec:
  tls:
  - hosts:
    # Change this to your own domain
    - demo.example.net
    secretName: secret-demo

EOF
```


### Preparation Step Resource7

<!-- @createResource7 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/cert-manager/issuer.yaml
apiVersion: certmanager.k8s.io/v1alpha1
kind: Issuer
metadata:
  name: issuer
spec:
  acme:
    # Change this to your own email
    email: demo@example.net
    server: https://acme-staging-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: issuer-letsencrypt-account-key-secret
    http01: {}

EOF
```


### Preparation Step Resource8

<!-- @createResource8 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/development/certificate.yaml
apiVersion: certmanager.k8s.io/v1alpha1
kind: Certificate
metadata:
  name: certificate
spec:
  secretName: secret-development-demo
  # Change this to your own domain
  commonName: development.demo.example.net
  dnsNames:
  # Change this to your own domain
  - development.demo.example.net
  issuerRef:
    name: issuer
    kind: Issuer
  acme: null

EOF
```


### Preparation Step Resource9

<!-- @createResource9 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/development/ingress.yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: ingress
spec:
  tls:
  - hosts:
    # Change this to your own domain
    - development.demo.example.net
    secretName: secret-development-demo

EOF
```


### Preparation Step Resource10

<!-- @createResource10 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/development/issuer.yaml
apiVersion: certmanager.k8s.io/v1alpha1
kind: Issuer
metadata:
  name: issuer
spec:
  ca:
    secretName: ca-secret-development
  acme: null

EOF
```


### Preparation Step Resource11

<!-- @createResource11 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/development/selfsigned-certificate.yaml
apiVersion: certmanager.k8s.io/v1alpha1
kind: Certificate
metadata:
  name: selfsigned-certificate
spec:
  secretName: ca-secret-development
  commonName: "development"
  isCA: true
  issuerRef:
    name: selfsigned-issuer
    kind: Issuer

EOF
```


### Preparation Step Resource12

<!-- @createResource12 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/development/selfsigned-issuer.yaml
apiVersion: certmanager.k8s.io/v1alpha1
kind: Issuer
metadata:
  name: selfsigned-issuer
spec:
  selfSigned: {}

EOF
```


### Preparation Step Resource13

<!-- @createResource13 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/production/certificate.yaml
apiVersion: certmanager.k8s.io/v1alpha1
kind: Certificate
metadata:
  name: certificate
spec:
  secretName: secret-www-demo
  # Change this to your own domain
  commonName: www.demo.example.net
  dnsNames:
  # Change this to your own domain
  - www.demo.example.net
  issuerRef:
    name: issuer
    kind: Issuer
  acme:
    config:
    - http01:
        ingress: ingress
      domains:
      # Change this to your own domain
      - www.demo.example.net

EOF
```


### Preparation Step Resource14

<!-- @createResource14 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/production/ingress.yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: ingress
spec:
  tls:
  - hosts:
    # Change this to your own domain
    - www.demo.example.net
    secretName: secret-www-demo

EOF
```


### Preparation Step Resource15

<!-- @createResource15 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/production/issuer.yaml
apiVersion: certmanager.k8s.io/v1alpha1
kind: Issuer
metadata:
  name: issuer
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: issuer-letsencrypt-account-key-secret-production

EOF
```


### Preparation Step Resource16

<!-- @createResource16 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/staging/certificate.yaml
apiVersion: certmanager.k8s.io/v1alpha1
kind: Certificate
metadata:
  name: certificate
spec:
  secretName: secret-staging-demo
  # Change this to your own domain
  commonName: staging.demo.example.net
  dnsNames:
  # Change this to your own domain
  - staging.demo.example.net
  issuerRef:
    name: issuer
    kind: Issuer
  acme: null

EOF
```


### Preparation Step Resource17

<!-- @createResource17 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/staging/ingress.yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: ingress
spec:
  tls:
  - hosts:
    # Change this to your own domain
    - staging.demo.example.net
    secretName: secret-staging-demo

EOF
```


### Preparation Step Resource18

<!-- @createResource18 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/staging/issuer.yaml
apiVersion: certmanager.k8s.io/v1alpha1
kind: Issuer
metadata:
  name: issuer
spec:
  ca:
    secretName: ca-secret
  acme: null

EOF
```


### Preparation Step Other0

<!-- @createOther0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/staging/secret/tls.crt
-----BEGIN CERTIFICATE-----
MIIC/zCCAeegAwIBAgIJAJxgKkyFg6K6MA0GCSqGSIb3DQEBCwUAMBYxFDASBgNV
BAMMC2NvbW1vbl9uYW1lMB4XDTE5MDMxMTEyMzk1NloXDTI5MDMwODEyMzk1Nlow
FjEUMBIGA1UEAwwLY29tbW9uX25hbWUwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAw
ggEKAoIBAQDIilLJRRoRORC0/9pjOnHpOSUroFAYfK4TyKngSIjgYNaCXAVqz5tz
w7nrEIUfieib1k0v5xs72YDI+NS3JO95fYsLHQ6Um1YR7uI5jlA0yNtlzDThdtiS
yVMiUMgYDx8kc1B42w0QTApYADA9dTTl9M/dwhDlQeFZDEXB/7PIIZ6N6vDT4P43
147rAc9SAKE3FXRfnUvfyhsKW8kNzdZnj+FoPMlk6r87SOzQ/CIR14A9lOj0e3KM
r5M7rXqf0Nr7hLcOXV7XcJsR1ter4UCjuR/E0I9cYfSZD0UWFc7M2k2J437U4vE3
+wUTjzqV6tLpC4GtM4IsgSbM7N6/UlCVAgMBAAGjUDBOMB0GA1UdDgQWBBTtfotG
xSup51IMTkFodaIZ8KZGVTAfBgNVHSMEGDAWgBTtfotGxSup51IMTkFodaIZ8KZG
VTAMBgNVHRMEBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQBlWJO0SQz2gpH4nJgP
gwU+46RZwXqJDu4fqLyXjYM9SXIwOadxD4jzNIGE10ez352Ms3Cg8C570TP1DMxW
2SC+ZxGRk3/DYxhcINeZLPlM3tNmFbWwCBhejWdru+mMqZCm8CqPkOnIDdkOBlE4
xpsDcC4oR+vutpuo44wijPD3WDcnCSOzeRVTp/WGUwYfSeST9ZZisxyKcB1QcrC8
LPGPqUI/AJWoQylIN/QAYvXbs2TGaODah+Fzv48B6A5Kt1wZrB0jH2fzrgkx4gwh
DU/KszT5eIM5BthBxd0I2YUJYJpD5nsGQ2BdepanOIh7qGwUaK+HpNGkrCHawnK/
hx/X
-----END CERTIFICATE-----
EOF
```


### Preparation Step Other1

<!-- @createOther1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlays/staging/secret/tls.key
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAyIpSyUUaETkQtP/aYzpx6TklK6BQGHyuE8ip4EiI4GDWglwF
as+bc8O56xCFH4nom9ZNL+cbO9mAyPjUtyTveX2LCx0OlJtWEe7iOY5QNMjbZcw0
4XbYkslTIlDIGA8fJHNQeNsNEEwKWAAwPXU05fTP3cIQ5UHhWQxFwf+zyCGejerw
0+D+N9eO6wHPUgChNxV0X51L38obClvJDc3WZ4/haDzJZOq/O0js0PwiEdeAPZTo
9HtyjK+TO616n9Da+4S3Dl1e13CbEdbXq+FAo7kfxNCPXGH0mQ9FFhXOzNpNieN+
1OLxN/sFE486lerS6QuBrTOCLIEmzOzev1JQlQIDAQABAoIBAQC6RloNoIVNGC+f
oMRvRVuH4k/XjSq5BB8CO3Mn5NhXazv1jJpvk3X2+whYA1lUaVpKXq4F3+qZFjic
9R1JHSLgO7AK26uud/dj4vv5sGpqDWRV03APOcCD4EO7bUPfrTQlPIO0LuychbVZ
9prYi2VecJ8ggmIFQcObXl3xjJ2nN4VvYF29nmmtm5vCMxU/jc56JTuAoU17Zmdu
EbZNvEoPXBrazjNNec4lqqXQS36OOyLxJhTovaOnOaCEi/HF8fQefJ51GYQqGsDJ
JheGEKZn1R8taGi3ZpZ60ngeUgz1YKrbGdWgJzI/67DTni86MBE81snLXgkQO6CA
j1H+pbLBAoGBAPGRXFovCZhkxBld4suLq1Lfy8vZpKmJIQ6SQAnaSLDectOGRXHo
WnJ6yHo/Jc82puUrv1ZksSGrf0lakbetV9aaJX6ug8uBNO2rw0xGCYRdoif6r2RR
yZ5TnnvC6sN6nwFYYSAjETUkkJJFxq0ArzKMVANbGeBj4rBhqMXNBWsRAoGBANSF
eJXrsCb76sRMCj0NkCJmSyk8mOhUR9WotGYsHtPPfkwuUoEvEknV7to8WHeHJcBD
wowchy84fpbAj7bHAerdXhKIIPXF1sVMT41H+cV/jsLb4/Rm99/XzauEsIMYi1ue
mv9dJC27qLfMrnsJLEENglTCOTzLqBO8WZJgMyVFAoGAWs+3dRurssNmyNZ3lOdL
n5sMJPULpsQrTiwCsPGDVCI77nLSlnCv18t6pCIrF4vHD+3zPwoZYLv03OGUWAVt
OPq3z7jRSOaovBRPFdRabY05kWf3GXJ5pfBvar0qvhPRxJKx6H/mTyEQzDw45P6V
3h3M03oi7yz9oisEZF+fgtECgYBdJTK44tgN/hPjfUBvieZGbXc716ddDLN/XbXT
ojrQsvyz/wmCPVNSsUVCuXg8yyssnYZDSq2lcKlrAXL7tTWN7wAwNyHbFp8PUmb7
kTRT75hup1m94PN7VGZ8amfCzZsmyNk+W2Bj+v/zru46VsbirD0XURktIEXGgKLx
mOBR0QKBgDFZUO2uvb9ZFa1u0iCtt3xSHRtC66QbkQaneOUc0QR0Zm6qZj1tSMbm
fRwlU+DlumNs4SBXr9id0dU+5Mmr38lkb78rj4DMGNGbJVegJ7669qXuvJIHIUEv
69RjFZ8OlB9rX7l/4Uy5Y54xDvcwQsnoZGBnS7ie8o49nYRcRKf7
-----END RSA PRIVATE KEY-----
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/overlays/cert-manager -o ${DEMO_HOME}/actual/cert-manager.yaml
kustomize build ${DEMO_HOME}/overlays/development -o ${DEMO_HOME}/actual/development.yaml
kustomize build ${DEMO_HOME}/overlays/multi-environment -o ${DEMO_HOME}/actual/multi-environment.yaml
kustomize build ${DEMO_HOME}/overlays/production -o ${DEMO_HOME}/actual/production.yaml
kustomize build ${DEMO_HOME}/overlays/staging -o ${DEMO_HOME}/actual/staging.yaml
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/cert-manager.yaml
apiVersion: v1
kind: Namespace
metadata:
  labels:
    app: helloweb
  name: helloweb
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: helloweb
  name: helloweb-service
  namespace: helloweb
spec:
  ports:
  - port: 8080
    protocol: TCP
    targetPort: 8080
  selector:
    app: helloweb
  type: NodePort
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: helloweb
  name: helloweb-deployment
  namespace: helloweb
spec:
  selector:
    matchLabels:
      app: helloweb
  template:
    metadata:
      labels:
        app: helloweb
    spec:
      containers:
      - image: gcr.io/google-samples/hello-app:1.0
        name: hello-app
        ports:
        - containerPort: 8080
---
apiVersion: certmanager.k8s.io/v1alpha1
kind: Certificate
metadata:
  labels:
    app: helloweb
  name: certificate
  namespace: helloweb
spec:
  acme:
    config:
    - domains:
      - demo.example.net
      http01:
        ingress: ingress
  commonName: demo.example.net
  dnsNames:
  - demo.example.net
  issuerRef:
    kind: Issuer
    name: issuer
  renewBefore: 360h
  secretName: secret-demo
---
apiVersion: certmanager.k8s.io/v1alpha1
kind: Issuer
metadata:
  labels:
    app: helloweb
  name: issuer
  namespace: helloweb
spec:
  acme:
    email: demo@example.net
    http01: {}
    privateKeySecretRef:
      name: issuer-letsencrypt-account-key-secret
    server: https://acme-staging-v02.api.letsencrypt.org/directory
---
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.global-static-ip-name: web-static-ip
  labels:
    app: helloweb
  name: ingress
  namespace: helloweb
spec:
  rules:
  - http:
      paths:
      - backend:
          serviceName: helloweb-service
          servicePort: 8080
        path: /
  tls:
  - hosts:
    - demo.example.net
    secretName: secret-demo
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/development.yaml
apiVersion: v1
kind: Namespace
metadata:
  labels:
    app: helloweb-development
  name: helloweb-development
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: helloweb-development
  name: helloweb-service-development
  namespace: helloweb-development
spec:
  ports:
  - port: 8080
    protocol: TCP
    targetPort: 8080
  selector:
    app: helloweb-development
  type: NodePort
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: helloweb-development
  name: helloweb-deployment-development
  namespace: helloweb-development
spec:
  selector:
    matchLabels:
      app: helloweb-development
  template:
    metadata:
      labels:
        app: helloweb-development
    spec:
      containers:
      - image: gcr.io/google-samples/hello-app:1.0
        name: hello-app
        ports:
        - containerPort: 8080
---
apiVersion: certmanager.k8s.io/v1alpha1
kind: Certificate
metadata:
  labels:
    app: helloweb-development
  name: certificate-development
  namespace: helloweb-development
spec:
  commonName: development.demo.example.net
  dnsNames:
  - development.demo.example.net
  issuerRef:
    kind: Issuer
    name: issuer-development
  renewBefore: 360h
  secretName: secret-development-demo
---
apiVersion: certmanager.k8s.io/v1alpha1
kind: Certificate
metadata:
  labels:
    app: helloweb-development
  name: selfsigned-certificate-development
  namespace: helloweb-development
spec:
  commonName: development
  isCA: true
  issuerRef:
    kind: Issuer
    name: selfsigned-issuer-development
  secretName: ca-secret-development
---
apiVersion: certmanager.k8s.io/v1alpha1
kind: Issuer
metadata:
  labels:
    app: helloweb-development
  name: issuer-development
  namespace: helloweb-development
spec:
  ca:
    secretName: ca-secret-development
---
apiVersion: certmanager.k8s.io/v1alpha1
kind: Issuer
metadata:
  labels:
    app: helloweb-development
  name: selfsigned-issuer-development
  namespace: helloweb-development
spec:
  selfSigned: {}
---
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.global-static-ip-name: web-static-ip
  labels:
    app: helloweb-development
  name: ingress-development
  namespace: helloweb-development
spec:
  rules:
  - http:
      paths:
      - backend:
          serviceName: helloweb-service-development
          servicePort: 8080
        path: /
  tls:
  - hosts:
    - development.demo.example.net
    secretName: secret-development-demo
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/multi-environment.yaml
apiVersion: v1
kind: Namespace
metadata:
  labels:
    app: helloweb-development
  name: helloweb-development
---
apiVersion: v1
kind: Namespace
metadata:
  labels:
    app: helloweb-production
  name: helloweb-production
---
apiVersion: v1
kind: Namespace
metadata:
  labels:
    app: helloweb-staging
  name: helloweb-staging
---
apiVersion: v1
data:
  tls.crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUMvekNDQWVlZ0F3SUJBZ0lKQUp4Z0treUZnNks2TUEwR0NTcUdTSWIzRFFFQkN3VUFNQll4RkRBU0JnTlYKQkFNTUMyTnZiVzF2Ymw5dVlXMWxNQjRYRFRFNU1ETXhNVEV5TXprMU5sb1hEVEk1TURNd09ERXlNemsxTmxvdwpGakVVTUJJR0ExVUVBd3dMWTI5dGJXOXVYMjVoYldVd2dnRWlNQTBHQ1NxR1NJYjNEUUVCQVFVQUE0SUJEd0F3CmdnRUtBb0lCQVFESWlsTEpSUm9ST1JDMC85cGpPbkhwT1NVcm9GQVlmSzRUeUtuZ1NJamdZTmFDWEFWcXo1dHoKdzduckVJVWZpZWliMWswdjV4czcyWURJK05TM0pPOTVmWXNMSFE2VW0xWVI3dUk1amxBMHlOdGx6RFRoZHRpUwp5Vk1pVU1nWUR4OGtjMUI0MncwUVRBcFlBREE5ZFRUbDlNL2R3aERsUWVGWkRFWEIvN1BJSVo2TjZ2RFQ0UDQzCjE0N3JBYzlTQUtFM0ZYUmZuVXZmeWhzS1c4a056ZFpuaitGb1BNbGs2cjg3U096US9DSVIxNEE5bE9qMGUzS00KcjVNN3JYcWYwTnI3aExjT1hWN1hjSnNSMXRlcjRVQ2p1Ui9FMEk5Y1lmU1pEMFVXRmM3TTJrMko0MzdVNHZFMword1VUanpxVjZ0THBDNEd0TTRJc2dTYk03TjYvVWxDVkFnTUJBQUdqVURCT01CMEdBMVVkRGdRV0JCVHRmb3RHCnhTdXA1MUlNVGtGb2RhSVo4S1pHVlRBZkJnTlZIU01FR0RBV2dCVHRmb3RHeFN1cDUxSU1Ua0ZvZGFJWjhLWkcKVlRBTUJnTlZIUk1FQlRBREFRSC9NQTBHQ1NxR1NJYjNEUUVCQ3dVQUE0SUJBUUJsV0pPMFNRejJncEg0bkpnUApnd1UrNDZSWndYcUpEdTRmcUx5WGpZTTlTWEl3T2FkeEQ0anpOSUdFMTBlejM1Mk1zM0NnOEM1NzBUUDFETXhXCjJTQytaeEdSazMvRFl4aGNJTmVaTFBsTTN0Tm1GYld3Q0JoZWpXZHJ1K21NcVpDbThDcVBrT25JRGRrT0JsRTQKeHBzRGNDNG9SK3Z1dHB1bzQ0d2lqUEQzV0RjbkNTT3plUlZUcC9XR1V3WWZTZVNUOVpaaXN4eUtjQjFRY3JDOApMUEdQcVVJL0FKV29ReWxJTi9RQVl2WGJzMlRHYU9EYWgrRnp2NDhCNkE1S3Qxd1pyQjBqSDJmenJna3g0Z3doCkRVL0tzelQ1ZUlNNUJ0aEJ4ZDBJMllVSllKcEQ1bnNHUTJCZGVwYW5PSWg3cUd3VWFLK0hwTkdrckNIYXduSy8KaHgvWAotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
  tls.key: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb3dJQkFBS0NBUUVBeUlwU3lVVWFFVGtRdFAvYVl6cHg2VGtsSzZCUUdIeXVFOGlwNEVpSTRHRFdnbHdGCmFzK2JjOE81NnhDRkg0bm9tOVpOTCtjYk85bUF5UGpVdHlUdmVYMkxDeDBPbEp0V0VlN2lPWTVRTk1qYlpjdzAKNFhiWWtzbFRJbERJR0E4ZkpITlFlTnNORUV3S1dBQXdQWFUwNWZUUDNjSVE1VUhoV1F4RndmK3p5Q0dlamVydwowK0QrTjllTzZ3SFBVZ0NoTnhWMFg1MUwzOG9iQ2x2SkRjM1daNC9oYUR6SlpPcS9PMGpzMFB3aUVkZUFQWlRvCjlIdHlqSytUTzYxNm45RGErNFMzRGwxZTEzQ2JFZGJYcStGQW83a2Z4TkNQWEdIMG1ROUZGaFhPek5wTmllTisKMU9MeE4vc0ZFNDg2bGVyUzZRdUJyVE9DTElFbXpPemV2MUpRbFFJREFRQUJBb0lCQVFDNlJsb05vSVZOR0MrZgpvTVJ2UlZ1SDRrL1hqU3E1QkI4Q08zTW41TmhYYXp2MWpKcHZrM1gyK3doWUExbFVhVnBLWHE0RjMrcVpGamljCjlSMUpIU0xnTzdBSzI2dXVkL2RqNHZ2NXNHcHFEV1JWMDNBUE9jQ0Q0RU83YlVQZnJUUWxQSU8wTHV5Y2hiVloKOXByWWkyVmVjSjhnZ21JRlFjT2JYbDN4akoybk40VnZZRjI5bm1tdG01dkNNeFUvamM1NkpUdUFvVTE3Wm1kdQpFYlpOdkVvUFhCcmF6ak5OZWM0bHFxWFFTMzZPT3lMeEpoVG92YU9uT2FDRWkvSEY4ZlFlZko1MUdZUXFHc0RKCkpoZUdFS1puMVI4dGFHaTNacFo2MG5nZVVnejFZS3JiR2RXZ0p6SS82N0RUbmk4Nk1CRTgxc25MWGdrUU82Q0EKajFIK3BiTEJBb0dCQVBHUlhGb3ZDWmhreEJsZDRzdUxxMUxmeTh2WnBLbUpJUTZTUUFuYVNMRGVjdE9HUlhIbwpXbko2eUhvL0pjODJwdVVydjFaa3NTR3JmMGxha2JldFY5YWFKWDZ1Zzh1Qk5PMnJ3MHhHQ1lSZG9pZjZyMlJSCnlaNVRubnZDNnNONm53RllZU0FqRVRVa2tKSkZ4cTBBcnpLTVZBTmJHZUJqNHJCaHFNWE5CV3NSQW9HQkFOU0YKZUpYcnNDYjc2c1JNQ2owTmtDSm1TeWs4bU9oVVI5V290R1lzSHRQUGZrd3VVb0V2RWtuVjd0bzhXSGVISmNCRAp3b3djaHk4NGZwYkFqN2JIQWVyZFhoS0lJUFhGMXNWTVQ0MUgrY1YvanNMYjQvUm05OS9YemF1RXNJTVlpMXVlCm12OWRKQzI3cUxmTXJuc0pMRUVOZ2xUQ09UekxxQk84V1pKZ015VkZBb0dBV3MrM2RSdXJzc05teU5aM2xPZEwKbjVzTUpQVUxwc1FyVGl3Q3NQR0RWQ0k3N25MU2xuQ3YxOHQ2cENJckY0dkhEKzN6UHdvWllMdjAzT0dVV0FWdApPUHEzejdqUlNPYW92QlJQRmRSYWJZMDVrV2YzR1hKNXBmQnZhcjBxdmhQUnhKS3g2SC9tVHlFUXpEdzQ1UDZWCjNoM00wM29pN3l6OW9pc0VaRitmZ3RFQ2dZQmRKVEs0NHRnTi9oUGpmVUJ2aWVaR2JYYzcxNmRkRExOL1hiWFQKb2pyUXN2eXovd21DUFZOU3NVVkN1WGc4eXlzc25ZWkRTcTJsY0tsckFYTDd0VFdON3dBd055SGJGcDhQVW1iNwprVFJUNzVodXAxbTk0UE43VkdaOGFtZkN6WnNteU5rK1cyQmordi96cnU0NlZzYmlyRDBYVVJrdElFWEdnS0x4Cm1PQlIwUUtCZ0RGWlVPMnV2YjlaRmExdTBpQ3R0M3hTSFJ0QzY2UWJrUWFuZU9VYzBRUjBabTZxWmoxdFNNYm0KZlJ3bFUrRGx1bU5zNFNCWHI5aWQwZFUrNU1tcjM4bGtiNzhyajRETUdOR2JKVmVnSjc2NjlxWHV2SklISVVFdgo2OVJqRlo4T2xCOXJYN2wvNFV5NVk1NHhEdmN3UXNub1pHQm5TN2llOG80OW5ZUmNSS2Y3Ci0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg==
kind: Secret
metadata:
  labels:
    app: helloweb-staging
  name: ca-secret-staging-tbd8c845kt
  namespace: helloweb-staging
type: kubernetes.io/tls
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: helloweb-development
  name: helloweb-service-development
  namespace: helloweb-development
spec:
  ports:
  - port: 8080
    protocol: TCP
    targetPort: 8080
  selector:
    app: helloweb-development
  type: NodePort
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: helloweb-production
  name: helloweb-service-production
  namespace: helloweb-production
spec:
  ports:
  - port: 8080
    protocol: TCP
    targetPort: 8080
  selector:
    app: helloweb-production
  type: NodePort
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: helloweb-staging
  name: helloweb-service-staging
  namespace: helloweb-staging
spec:
  ports:
  - port: 8080
    protocol: TCP
    targetPort: 8080
  selector:
    app: helloweb-staging
  type: NodePort
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: helloweb-development
  name: helloweb-deployment-development
  namespace: helloweb-development
spec:
  selector:
    matchLabels:
      app: helloweb-development
  template:
    metadata:
      labels:
        app: helloweb-development
    spec:
      containers:
      - image: gcr.io/google-samples/hello-app:1.0
        name: hello-app
        ports:
        - containerPort: 8080
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: helloweb-production
  name: helloweb-deployment-production
  namespace: helloweb-production
spec:
  selector:
    matchLabels:
      app: helloweb-production
  template:
    metadata:
      labels:
        app: helloweb-production
    spec:
      containers:
      - image: gcr.io/google-samples/hello-app:1.0
        name: hello-app
        ports:
        - containerPort: 8080
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: helloweb-staging
  name: helloweb-deployment-staging
  namespace: helloweb-staging
spec:
  selector:
    matchLabels:
      app: helloweb-staging
  template:
    metadata:
      labels:
        app: helloweb-staging
    spec:
      containers:
      - image: gcr.io/google-samples/hello-app:1.0
        name: hello-app
        ports:
        - containerPort: 8080
---
apiVersion: certmanager.k8s.io/v1alpha1
kind: Certificate
metadata:
  labels:
    app: helloweb-development
  name: certificate-development
  namespace: helloweb-development
spec:
  commonName: development.demo.example.net
  dnsNames:
  - development.demo.example.net
  issuerRef:
    kind: Issuer
    name: issuer-development
  renewBefore: 360h
  secretName: secret-development-demo
---
apiVersion: certmanager.k8s.io/v1alpha1
kind: Certificate
metadata:
  labels:
    app: helloweb-development
  name: selfsigned-certificate-development
  namespace: helloweb-development
spec:
  commonName: development
  isCA: true
  issuerRef:
    kind: Issuer
    name: selfsigned-issuer-development
  secretName: ca-secret-development
---
apiVersion: certmanager.k8s.io/v1alpha1
kind: Certificate
metadata:
  labels:
    app: helloweb-production
  name: certificate-production
  namespace: helloweb-production
spec:
  acme:
    config:
    - domains:
      - www.demo.example.net
      http01:
        ingress: ingress-production
  commonName: www.demo.example.net
  dnsNames:
  - www.demo.example.net
  issuerRef:
    kind: Issuer
    name: issuer-production
  renewBefore: 360h
  secretName: secret-www-demo
---
apiVersion: certmanager.k8s.io/v1alpha1
kind: Certificate
metadata:
  labels:
    app: helloweb-staging
  name: certificate-staging
  namespace: helloweb-staging
spec:
  commonName: staging.demo.example.net
  dnsNames:
  - staging.demo.example.net
  issuerRef:
    kind: Issuer
    name: issuer-staging
  renewBefore: 360h
  secretName: secret-staging-demo
---
apiVersion: certmanager.k8s.io/v1alpha1
kind: Issuer
metadata:
  labels:
    app: helloweb-development
  name: issuer-development
  namespace: helloweb-development
spec:
  ca:
    secretName: ca-secret-development
---
apiVersion: certmanager.k8s.io/v1alpha1
kind: Issuer
metadata:
  labels:
    app: helloweb-development
  name: selfsigned-issuer-development
  namespace: helloweb-development
spec:
  selfSigned: {}
---
apiVersion: certmanager.k8s.io/v1alpha1
kind: Issuer
metadata:
  labels:
    app: helloweb-production
  name: issuer-production
  namespace: helloweb-production
spec:
  acme:
    email: demo@example.net
    http01: {}
    privateKeySecretRef:
      name: issuer-letsencrypt-account-key-secret-production
    server: https://acme-v02.api.letsencrypt.org/directory
---
apiVersion: certmanager.k8s.io/v1alpha1
kind: Issuer
metadata:
  labels:
    app: helloweb-staging
  name: issuer-staging
  namespace: helloweb-staging
spec:
  ca:
    secretName: ca-secret-staging-tbd8c845kt
---
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.global-static-ip-name: web-static-ip
  labels:
    app: helloweb-development
  name: ingress-development
  namespace: helloweb-development
spec:
  rules:
  - http:
      paths:
      - backend:
          serviceName: helloweb-service-development
          servicePort: 8080
        path: /
  tls:
  - hosts:
    - development.demo.example.net
    secretName: secret-development-demo
---
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.global-static-ip-name: web-static-ip
  labels:
    app: helloweb-production
  name: ingress-production
  namespace: helloweb-production
spec:
  rules:
  - http:
      paths:
      - backend:
          serviceName: helloweb-service-production
          servicePort: 8080
        path: /
  tls:
  - hosts:
    - www.demo.example.net
    secretName: secret-www-demo
---
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.global-static-ip-name: web-static-ip
  labels:
    app: helloweb-staging
  name: ingress-staging
  namespace: helloweb-staging
spec:
  rules:
  - http:
      paths:
      - backend:
          serviceName: helloweb-service-staging
          servicePort: 8080
        path: /
  tls:
  - hosts:
    - staging.demo.example.net
    secretName: secret-staging-demo
EOF
```


### Verification Step Expected3

<!-- @createExpected3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/production.yaml
apiVersion: v1
kind: Namespace
metadata:
  labels:
    app: helloweb-production
  name: helloweb-production
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: helloweb-production
  name: helloweb-service-production
  namespace: helloweb-production
spec:
  ports:
  - port: 8080
    protocol: TCP
    targetPort: 8080
  selector:
    app: helloweb-production
  type: NodePort
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: helloweb-production
  name: helloweb-deployment-production
  namespace: helloweb-production
spec:
  selector:
    matchLabels:
      app: helloweb-production
  template:
    metadata:
      labels:
        app: helloweb-production
    spec:
      containers:
      - image: gcr.io/google-samples/hello-app:1.0
        name: hello-app
        ports:
        - containerPort: 8080
---
apiVersion: certmanager.k8s.io/v1alpha1
kind: Certificate
metadata:
  labels:
    app: helloweb-production
  name: certificate-production
  namespace: helloweb-production
spec:
  acme:
    config:
    - domains:
      - www.demo.example.net
      http01:
        ingress: ingress-production
  commonName: www.demo.example.net
  dnsNames:
  - www.demo.example.net
  issuerRef:
    kind: Issuer
    name: issuer-production
  renewBefore: 360h
  secretName: secret-www-demo
---
apiVersion: certmanager.k8s.io/v1alpha1
kind: Issuer
metadata:
  labels:
    app: helloweb-production
  name: issuer-production
  namespace: helloweb-production
spec:
  acme:
    email: demo@example.net
    http01: {}
    privateKeySecretRef:
      name: issuer-letsencrypt-account-key-secret-production
    server: https://acme-v02.api.letsencrypt.org/directory
---
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.global-static-ip-name: web-static-ip
  labels:
    app: helloweb-production
  name: ingress-production
  namespace: helloweb-production
spec:
  rules:
  - http:
      paths:
      - backend:
          serviceName: helloweb-service-production
          servicePort: 8080
        path: /
  tls:
  - hosts:
    - www.demo.example.net
    secretName: secret-www-demo
EOF
```


### Verification Step Expected4

<!-- @createExpected4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/staging.yaml
apiVersion: v1
kind: Namespace
metadata:
  labels:
    app: helloweb-staging
  name: helloweb-staging
---
apiVersion: v1
data:
  tls.crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUMvekNDQWVlZ0F3SUJBZ0lKQUp4Z0treUZnNks2TUEwR0NTcUdTSWIzRFFFQkN3VUFNQll4RkRBU0JnTlYKQkFNTUMyTnZiVzF2Ymw5dVlXMWxNQjRYRFRFNU1ETXhNVEV5TXprMU5sb1hEVEk1TURNd09ERXlNemsxTmxvdwpGakVVTUJJR0ExVUVBd3dMWTI5dGJXOXVYMjVoYldVd2dnRWlNQTBHQ1NxR1NJYjNEUUVCQVFVQUE0SUJEd0F3CmdnRUtBb0lCQVFESWlsTEpSUm9ST1JDMC85cGpPbkhwT1NVcm9GQVlmSzRUeUtuZ1NJamdZTmFDWEFWcXo1dHoKdzduckVJVWZpZWliMWswdjV4czcyWURJK05TM0pPOTVmWXNMSFE2VW0xWVI3dUk1amxBMHlOdGx6RFRoZHRpUwp5Vk1pVU1nWUR4OGtjMUI0MncwUVRBcFlBREE5ZFRUbDlNL2R3aERsUWVGWkRFWEIvN1BJSVo2TjZ2RFQ0UDQzCjE0N3JBYzlTQUtFM0ZYUmZuVXZmeWhzS1c4a056ZFpuaitGb1BNbGs2cjg3U096US9DSVIxNEE5bE9qMGUzS00KcjVNN3JYcWYwTnI3aExjT1hWN1hjSnNSMXRlcjRVQ2p1Ui9FMEk5Y1lmU1pEMFVXRmM3TTJrMko0MzdVNHZFMword1VUanpxVjZ0THBDNEd0TTRJc2dTYk03TjYvVWxDVkFnTUJBQUdqVURCT01CMEdBMVVkRGdRV0JCVHRmb3RHCnhTdXA1MUlNVGtGb2RhSVo4S1pHVlRBZkJnTlZIU01FR0RBV2dCVHRmb3RHeFN1cDUxSU1Ua0ZvZGFJWjhLWkcKVlRBTUJnTlZIUk1FQlRBREFRSC9NQTBHQ1NxR1NJYjNEUUVCQ3dVQUE0SUJBUUJsV0pPMFNRejJncEg0bkpnUApnd1UrNDZSWndYcUpEdTRmcUx5WGpZTTlTWEl3T2FkeEQ0anpOSUdFMTBlejM1Mk1zM0NnOEM1NzBUUDFETXhXCjJTQytaeEdSazMvRFl4aGNJTmVaTFBsTTN0Tm1GYld3Q0JoZWpXZHJ1K21NcVpDbThDcVBrT25JRGRrT0JsRTQKeHBzRGNDNG9SK3Z1dHB1bzQ0d2lqUEQzV0RjbkNTT3plUlZUcC9XR1V3WWZTZVNUOVpaaXN4eUtjQjFRY3JDOApMUEdQcVVJL0FKV29ReWxJTi9RQVl2WGJzMlRHYU9EYWgrRnp2NDhCNkE1S3Qxd1pyQjBqSDJmenJna3g0Z3doCkRVL0tzelQ1ZUlNNUJ0aEJ4ZDBJMllVSllKcEQ1bnNHUTJCZGVwYW5PSWg3cUd3VWFLK0hwTkdrckNIYXduSy8KaHgvWAotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
  tls.key: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb3dJQkFBS0NBUUVBeUlwU3lVVWFFVGtRdFAvYVl6cHg2VGtsSzZCUUdIeXVFOGlwNEVpSTRHRFdnbHdGCmFzK2JjOE81NnhDRkg0bm9tOVpOTCtjYk85bUF5UGpVdHlUdmVYMkxDeDBPbEp0V0VlN2lPWTVRTk1qYlpjdzAKNFhiWWtzbFRJbERJR0E4ZkpITlFlTnNORUV3S1dBQXdQWFUwNWZUUDNjSVE1VUhoV1F4RndmK3p5Q0dlamVydwowK0QrTjllTzZ3SFBVZ0NoTnhWMFg1MUwzOG9iQ2x2SkRjM1daNC9oYUR6SlpPcS9PMGpzMFB3aUVkZUFQWlRvCjlIdHlqSytUTzYxNm45RGErNFMzRGwxZTEzQ2JFZGJYcStGQW83a2Z4TkNQWEdIMG1ROUZGaFhPek5wTmllTisKMU9MeE4vc0ZFNDg2bGVyUzZRdUJyVE9DTElFbXpPemV2MUpRbFFJREFRQUJBb0lCQVFDNlJsb05vSVZOR0MrZgpvTVJ2UlZ1SDRrL1hqU3E1QkI4Q08zTW41TmhYYXp2MWpKcHZrM1gyK3doWUExbFVhVnBLWHE0RjMrcVpGamljCjlSMUpIU0xnTzdBSzI2dXVkL2RqNHZ2NXNHcHFEV1JWMDNBUE9jQ0Q0RU83YlVQZnJUUWxQSU8wTHV5Y2hiVloKOXByWWkyVmVjSjhnZ21JRlFjT2JYbDN4akoybk40VnZZRjI5bm1tdG01dkNNeFUvamM1NkpUdUFvVTE3Wm1kdQpFYlpOdkVvUFhCcmF6ak5OZWM0bHFxWFFTMzZPT3lMeEpoVG92YU9uT2FDRWkvSEY4ZlFlZko1MUdZUXFHc0RKCkpoZUdFS1puMVI4dGFHaTNacFo2MG5nZVVnejFZS3JiR2RXZ0p6SS82N0RUbmk4Nk1CRTgxc25MWGdrUU82Q0EKajFIK3BiTEJBb0dCQVBHUlhGb3ZDWmhreEJsZDRzdUxxMUxmeTh2WnBLbUpJUTZTUUFuYVNMRGVjdE9HUlhIbwpXbko2eUhvL0pjODJwdVVydjFaa3NTR3JmMGxha2JldFY5YWFKWDZ1Zzh1Qk5PMnJ3MHhHQ1lSZG9pZjZyMlJSCnlaNVRubnZDNnNONm53RllZU0FqRVRVa2tKSkZ4cTBBcnpLTVZBTmJHZUJqNHJCaHFNWE5CV3NSQW9HQkFOU0YKZUpYcnNDYjc2c1JNQ2owTmtDSm1TeWs4bU9oVVI5V290R1lzSHRQUGZrd3VVb0V2RWtuVjd0bzhXSGVISmNCRAp3b3djaHk4NGZwYkFqN2JIQWVyZFhoS0lJUFhGMXNWTVQ0MUgrY1YvanNMYjQvUm05OS9YemF1RXNJTVlpMXVlCm12OWRKQzI3cUxmTXJuc0pMRUVOZ2xUQ09UekxxQk84V1pKZ015VkZBb0dBV3MrM2RSdXJzc05teU5aM2xPZEwKbjVzTUpQVUxwc1FyVGl3Q3NQR0RWQ0k3N25MU2xuQ3YxOHQ2cENJckY0dkhEKzN6UHdvWllMdjAzT0dVV0FWdApPUHEzejdqUlNPYW92QlJQRmRSYWJZMDVrV2YzR1hKNXBmQnZhcjBxdmhQUnhKS3g2SC9tVHlFUXpEdzQ1UDZWCjNoM00wM29pN3l6OW9pc0VaRitmZ3RFQ2dZQmRKVEs0NHRnTi9oUGpmVUJ2aWVaR2JYYzcxNmRkRExOL1hiWFQKb2pyUXN2eXovd21DUFZOU3NVVkN1WGc4eXlzc25ZWkRTcTJsY0tsckFYTDd0VFdON3dBd055SGJGcDhQVW1iNwprVFJUNzVodXAxbTk0UE43VkdaOGFtZkN6WnNteU5rK1cyQmordi96cnU0NlZzYmlyRDBYVVJrdElFWEdnS0x4Cm1PQlIwUUtCZ0RGWlVPMnV2YjlaRmExdTBpQ3R0M3hTSFJ0QzY2UWJrUWFuZU9VYzBRUjBabTZxWmoxdFNNYm0KZlJ3bFUrRGx1bU5zNFNCWHI5aWQwZFUrNU1tcjM4bGtiNzhyajRETUdOR2JKVmVnSjc2NjlxWHV2SklISVVFdgo2OVJqRlo4T2xCOXJYN2wvNFV5NVk1NHhEdmN3UXNub1pHQm5TN2llOG80OW5ZUmNSS2Y3Ci0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg==
kind: Secret
metadata:
  labels:
    app: helloweb-staging
  name: ca-secret-staging-tbd8c845kt
  namespace: helloweb-staging
type: kubernetes.io/tls
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: helloweb-staging
  name: helloweb-service-staging
  namespace: helloweb-staging
spec:
  ports:
  - port: 8080
    protocol: TCP
    targetPort: 8080
  selector:
    app: helloweb-staging
  type: NodePort
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: helloweb-staging
  name: helloweb-deployment-staging
  namespace: helloweb-staging
spec:
  selector:
    matchLabels:
      app: helloweb-staging
  template:
    metadata:
      labels:
        app: helloweb-staging
    spec:
      containers:
      - image: gcr.io/google-samples/hello-app:1.0
        name: hello-app
        ports:
        - containerPort: 8080
---
apiVersion: certmanager.k8s.io/v1alpha1
kind: Certificate
metadata:
  labels:
    app: helloweb-staging
  name: certificate-staging
  namespace: helloweb-staging
spec:
  commonName: staging.demo.example.net
  dnsNames:
  - staging.demo.example.net
  issuerRef:
    kind: Issuer
    name: issuer-staging
  renewBefore: 360h
  secretName: secret-staging-demo
---
apiVersion: certmanager.k8s.io/v1alpha1
kind: Issuer
metadata:
  labels:
    app: helloweb-staging
  name: issuer-staging
  namespace: helloweb-staging
spec:
  ca:
    secretName: ca-secret-staging-tbd8c845kt
---
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.global-static-ip-name: web-static-ip
  labels:
    app: helloweb-staging
  name: ingress-staging
  namespace: helloweb-staging
spec:
  rules:
  - http:
      paths:
      - backend:
          serviceName: helloweb-service-staging
          servicePort: 8080
        path: /
  tls:
  - hosts:
    - staging.demo.example.net
    secretName: secret-staging-demo
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

