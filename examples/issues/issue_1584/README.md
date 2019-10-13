# Feature Test for Issue 1584


This folder contains files describing how to address [Issue 1584](https://github.com/kubernetes-sigs/kustomize/issues/1584)

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

namePrefix: demo-

commonLabels:
  app: api
  project: demo

configurations:
  - commonlabels.yaml

resources:
  - networkpolicy.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/commonlabels.yaml
commonLabels:
  - path: spec/egress/to/podSelector/matchLabels
    skip: true
    group: networking.k8s.io
    version: v1
    kind: NetworkPolicy
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/networkpolicy.yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: network-policy
  labels:
    service: network-policy
spec:
  policyTypes:
    - Ingress
    - Egress
  podSelector:
    matchLabels:
  ingress:
    - from:
        - podSelector:
            matchLabels:
    - ports:
        - port: 80
          protocol: TCP
  egress:
    - to:
        - podSelector:
            matchLabels:
              app: web
    - ports:
        - port: 53
          protocol: TCP
        - port: 53
          protocol: UDP
        - port: 80
          protocol: TCP
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
cat <<'EOF' >${DEMO_HOME}/expected/networking.k8s.io_v1_networkpolicy_demo-network-policy.yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  labels:
    app: api
    project: demo
    service: network-policy
  name: demo-network-policy
spec:
  egress:
  - to:
    - podSelector:
        matchLabels:
          app: web
  - ports:
    - port: 53
      protocol: TCP
    - port: 53
      protocol: UDP
    - port: 80
      protocol: TCP
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: api
          project: demo
  - ports:
    - port: 80
      protocol: TCP
  podSelector:
    matchLabels:
      app: api
      project: demo
  policyTypes:
  - Ingress
  - Egress
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

