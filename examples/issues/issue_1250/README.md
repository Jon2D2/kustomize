# Feature Test for Issue 1250


This folder contains files describing how to address [Issue 1250](https://github.com/kubernetes-sigs/kustomize/issues/1250)

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
resources:
  - configmap.yaml
  - ingress.yaml
  - service.yaml

configurations:
- kustomizeconfig.yaml

vars:
  - name: CNAME
    objref:
      apiVersion: networking.k8s.io/v1beta1
      kind: Ingress
      name: my-ingress
    fieldref:
      fieldpath: spec.rules[0].host
EOF
```


### Preparation Step KustomizeConfig0

<!-- @createKustomizeConfig0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/kustomizeconfig.yaml
varReference:
- path: data/CNAME
  kind: ConfigMap
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: configmap
data:
  ZZZ: $(Ingress.my-ingress.spec.rules[0].http.paths[0].backend.servicePort)
  CNAME: $(CNAME)
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/ingress.yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: my-ingress
spec:
  rules:
  - host: CNAME.DOMAIN.COM
    http:
      paths:
      - path: /
        backend:
          serviceName: service
          servicePort: 80
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: service
spec:
  ports:
  - name: http
    protocol: TCP
    port: 80
    targetPort: 80
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
cat <<'EOF' >${DEMO_HOME}/expected/networking.k8s.io_v1beta1_ingress_my-ingress.yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: my-ingress
spec:
  rules:
  - host: CNAME.DOMAIN.COM
    http:
      paths:
      - backend:
          serviceName: service
          servicePort: 80
        path: /
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_configmap_configmap.yaml
apiVersion: v1
data:
  CNAME: CNAME.DOMAIN.COM
  ZZZ: 80
kind: ConfigMap
metadata:
  name: configmap
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_service_service.yaml
apiVersion: v1
kind: Service
metadata:
  name: service
spec:
  ports:
  - name: http
    port: 80
    protocol: TCP
    targetPort: 80
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

