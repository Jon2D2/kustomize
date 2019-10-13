# Configuring Intel Device Plugins for Kubernetes

This folder is using kustomize to customize the default values for Intel demos.
See [here](https://github.com/intel/intel-device-plugins-for-kubernetes.git)

## Setup the workspace

First, define a place to work:

<!-- @makeWorkplace @test -->
```bash
DEMO_HOME=$(mktemp -d)
```

## Preparation

<!-- @makeDirectories @test -->
```bash
mkdir -p ${DEMO_HOME}/overlay
mkdir -p ${DEMO_HOME}/base
```

### Preparation Step KustomizationFile0

<!-- @createKustomizationFile0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay/kustomization.yaml
resources:
- ../base
- ./values.yaml

patchesStrategicMerge:
- intelfga-patch.yaml

patchesJson6902:
- target:
    version: v1
    kind: Pod
    name: dpdkqatuio
  path: dpdk-patch.yaml
EOF
```


### Preparation Step KustomizationFile1

<!-- @createKustomizationFile1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/kustomization.yaml
resources:
- crypto-perf-dpdk-pod-requesting-qat.yaml
- intelfpga-job.yaml
- intelgpu-job.yaml
- openssl-qat-engine-pod.yaml
EOF
```


### Preparation Step Resource0

<!-- @createResource0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay/dpdk-patch.yaml
- op: replace
  path: /spec/containers/0/resources/limits/memory
  value: 64Mi
EOF
```


### Preparation Step Resource1

<!-- @createResource1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay/intelfga-patch.yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: intelfpga-demo-job
spec:
  template:
    spec:
      containers:
      - command:
        name: intelfpga-demo-job-1
        resources:
          limits:
            hugepages-2Mi: 10Mi
EOF
```


### Preparation Step Resource2

<!-- @createResource2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/overlay/values.yaml
apiVersion: v1
kind: Values
metadata:
  name: global
spec:
  gpulimit: 3
EOF
```


### Preparation Step Resource3

<!-- @createResource3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/openssl-qat-engine-pod.yaml
apiVersion: v1
kind: Pod
metadata:
  name: openssl-qat-engine
spec:
  containers:
  - args:
    - while true; do sleep 300000; done;
    command:
    - /bin/bash
    - -c
    - --
    image: docker.io/library/openssl-qat-engine:devel
    imagePullPolicy: Never
    name: openssl-qat-engine
    resources:
      limits:
        qat.intel.com/generic: "1"
      requests:
        qat.intel.com/generic: "1"
    volumeMounts:
    - mountPath: /dev
      name: dev-mount
    - mountPath: /etc/c6xxvf_dev0.conf
      name: dev0
  runtimeClassName: kata-containers
  volumes:
  - hostPath:
      path: /dev
    name: dev-mount
  - hostPath:
      path: /etc/c6xxvf_dev0.conf
    name: dev0
EOF
```


### Preparation Step Resource4

<!-- @createResource4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/crypto-perf-dpdk-pod-requesting-qat.yaml
apiVersion: v1
kind: Pod
metadata:
  name: dpdkqatuio
spec:
  containers:
  - args:
    - while true; do sleep 300000; done;
    command:
    - /bin/bash
    - -c
    - --
    image: crypto-perf:devel
    imagePullPolicy: IfNotPresent
    name: dpdkcontainer
    resources:
      limits:
        cpu: "3"
        hugepages-2Mi: 1Gi
        memory: 128Mi
        qat.intel.com/generic: "4"
      requests:
        cpu: "3"
        hugepages-2Mi: 1Gi
        memory: 128Mi
        qat.intel.com/generic: "4"
    securityContext:
      capabilities:
        add:
        - IPC_LOCK
    volumeMounts:
    - mountPath: /dev/hugepages
      name: hugepage
  volumes:
  - emptyDir:
      medium: HugePages
    name: hugepage
EOF
```


### Preparation Step Resource5

<!-- @createResource5 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/intelgpu-job.yaml
apiVersion: batch/v1
kind: Job
metadata:
  labels:
    jobgroup: intelgpu-demo
  name: intelgpu-demo-job
spec:
  template:
    metadata:
      labels:
        jobgroup: intelgpu-demo
    spec:
      containers:
      - command:
        - /run-opencl-example.sh
        - /root/6-1/fft
        image: ubuntu-demo-opencl:devel
        imagePullPolicy: IfNotPresent
        name: intelgpu-demo-job-1
        resources:
          limits:
            gpu.intel.com/i915: $(Values.global.spec.gpulimit)
      restartPolicy: Never
EOF
```


### Preparation Step Resource6

<!-- @createResource6 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/base/intelfpga-job.yaml
apiVersion: batch/v1
kind: Job
metadata:
  labels:
    jobgroup: intelfpga-demo
  name: intelfpga-demo-job
spec:
  template:
    metadata:
      labels:
        jobgroup: intelfpga-demo
    spec:
      containers:
      - command:
        - /usr/bin/test_fpga.sh
        image: ubuntu-demo-opae:devel
        imagePullPolicy: IfNotPresent
        name: intelfpga-demo-job-1
        resources:
          limits:
            cpu: 1
            fpga.intel.com/af-d8424dc4a4a3c413f89e433683f9040b: 1
            hugepages-2Mi: 20Mi
        securityContext:
          capabilities:
            add:
            - IPC_LOCK
      restartPolicy: Never
EOF
```

## Execution

<!-- @build @test -->
```bash
mkdir ${DEMO_HOME}/actual
kustomize build ${DEMO_HOME}/overlay -o ${DEMO_HOME}/actual
```

## Verification

<!-- @createExpectedDir @test -->
```bash
mkdir ${DEMO_HOME}/expected
```


### Verification Step Expected0

<!-- @createExpected0 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_values_global.yaml
apiVersion: v1
kind: Values
metadata:
  name: global
spec:
  gpulimit: 3
EOF
```


### Verification Step Expected1

<!-- @createExpected1 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/batch_v1_job_intelfpga-demo-job.yaml
apiVersion: batch/v1
kind: Job
metadata:
  labels:
    jobgroup: intelfpga-demo
  name: intelfpga-demo-job
spec:
  template:
    metadata:
      labels:
        jobgroup: intelfpga-demo
    spec:
      containers:
      - image: ubuntu-demo-opae:devel
        imagePullPolicy: IfNotPresent
        name: intelfpga-demo-job-1
        resources:
          limits:
            cpu: 1
            fpga.intel.com/af-d8424dc4a4a3c413f89e433683f9040b: 1
            hugepages-2Mi: 10Mi
        securityContext:
          capabilities:
            add:
            - IPC_LOCK
      restartPolicy: Never
EOF
```


### Verification Step Expected2

<!-- @createExpected2 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/batch_v1_job_intelgpu-demo-job.yaml
apiVersion: batch/v1
kind: Job
metadata:
  labels:
    jobgroup: intelgpu-demo
  name: intelgpu-demo-job
spec:
  template:
    metadata:
      labels:
        jobgroup: intelgpu-demo
    spec:
      containers:
      - command:
        - /run-opencl-example.sh
        - /root/6-1/fft
        image: ubuntu-demo-opencl:devel
        imagePullPolicy: IfNotPresent
        name: intelgpu-demo-job-1
        resources:
          limits:
            gpu.intel.com/i915: 3
      restartPolicy: Never
EOF
```


### Verification Step Expected3

<!-- @createExpected3 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_pod_dpdkqatuio.yaml
apiVersion: v1
kind: Pod
metadata:
  name: dpdkqatuio
spec:
  containers:
  - args:
    - while true; do sleep 300000; done;
    command:
    - /bin/bash
    - -c
    - --
    image: crypto-perf:devel
    imagePullPolicy: IfNotPresent
    name: dpdkcontainer
    resources:
      limits:
        cpu: "3"
        hugepages-2Mi: 1Gi
        memory: 64Mi
        qat.intel.com/generic: "4"
      requests:
        cpu: "3"
        hugepages-2Mi: 1Gi
        memory: 128Mi
        qat.intel.com/generic: "4"
    securityContext:
      capabilities:
        add:
        - IPC_LOCK
    volumeMounts:
    - mountPath: /dev/hugepages
      name: hugepage
  volumes:
  - emptyDir:
      medium: HugePages
    name: hugepage
EOF
```


### Verification Step Expected4

<!-- @createExpected4 @test -->
```bash
cat <<'EOF' >${DEMO_HOME}/expected/~g_v1_pod_openssl-qat-engine.yaml
apiVersion: v1
kind: Pod
metadata:
  name: openssl-qat-engine
spec:
  containers:
  - args:
    - while true; do sleep 300000; done;
    command:
    - /bin/bash
    - -c
    - --
    image: docker.io/library/openssl-qat-engine:devel
    imagePullPolicy: Never
    name: openssl-qat-engine
    resources:
      limits:
        qat.intel.com/generic: "1"
      requests:
        qat.intel.com/generic: "1"
    volumeMounts:
    - mountPath: /dev
      name: dev-mount
    - mountPath: /etc/c6xxvf_dev0.conf
      name: dev0
  runtimeClassName: kata-containers
  volumes:
  - hostPath:
      path: /dev
    name: dev-mount
  - hostPath:
      path: /etc/c6xxvf_dev0.conf
    name: dev0
EOF
```


<!-- @compareActualToExpected @test -->
```bash
test 0 == \
$(diff -r $DEMO_HOME/actual $DEMO_HOME/expected | wc -l); \
echo $?
```

