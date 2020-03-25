// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package krusty_test

import (
	"testing"

	kusttest_test "sigs.k8s.io/kustomize/api/testutils/kusttest"
)

type inlineTest struct{}

func (ut *inlineTest) writeKustomization(th kusttest_test.Harness) {
	th.WriteK("/inlinesimple/", `
resources:
- cronjob.yaml
- deployment.yaml
- values.yaml

configurations:
- kustomizeconfig.yaml

vars:
- name : Values.shared.spec.env
  objref:
    apiVersion: v1
    kind: Values
    name: shared
  fieldref:
    fieldpath: spec.env
`)
}

func (ut *inlineTest) writeKustomConfig(th kusttest_test.Harness) {
	th.WriteF("/inlinesimple/kustomizeconfig.yaml", `
varReference:
- kind: Deployment
  path: spec/template/spec/containers[]/env

- kind: CronJob
  path: spec/jobTemplate/spec/template/spec/containers[]/env
`)
}

func (ut *inlineTest) writeCronJob(th kusttest_test.Harness) {
	th.WriteF("/inlinesimple/cronjob.yaml", `
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: wordpress-cron
  labels:
    app: wordpress
spec:
  schedule: "*/10 * * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - image: wordpress:4.8-apache
            name: wordpress
            command:
            - php
            args:
            - /path/to/wp-cron.php
            env: $(Values.shared.spec.env)
          restartPolicy: OnFailure
`)
}

func (ut *inlineTest) writeDnsPatch(th kusttest_test.Harness) {
	th.WriteF("/inlinesimple/deployment.yaml", `
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: wordpress
  labels:
    app: wordpress
spec:
  selector:
    matchLabels:
      app: wordpress
  template:
    metadata:
      labels:
        app: wordpress
    spec:
      containers:
      - image: wordpress:4.8-apache
        name: wordpress
        ports:
        - containerPort: 80
          name: wordpress
        env: $(Values.shared.spec.env)
`)
}

func (ut *inlineTest) writeValues(th kusttest_test.Harness) {
	th.WriteF("/inlinesimple/values.yaml", `
apiVersion: v1
kind: Values
metadata:
  name: shared
spec:
  env:
  - name: WORDPRESS_DB_USER
    valueFrom:
      secretKeyRef:
        name: wordpress-db-auth
        key: user
  - name: WORDPRESS_DB_PASSWORD
    valueFrom:
      secretKeyRef:
        name: wordpress-db-auth
        key: password
`)
}

func TestSimpleInline(t *testing.T) {
	ut := &inlineTest{}
	th := kusttest_test.MakeHarness(t)
	ut.writeKustomization(th)
	ut.writeKustomConfig(th)
	ut.writeCronJob(th)
	ut.writeDnsPatch(th)
	ut.writeValues(th)
	m := th.Run("/inlinesimple", th.MakeDefaultOptions())
	th.AssertActualEqualsExpected(m, `
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  labels:
    app: wordpress
  name: wordpress-cron
spec:
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - args:
            - /path/to/wp-cron.php
            command:
            - php
            env:
            - name: WORDPRESS_DB_USER
              valueFrom:
                secretKeyRef:
                  key: user
                  name: wordpress-db-auth
            - name: WORDPRESS_DB_PASSWORD
              valueFrom:
                secretKeyRef:
                  key: password
                  name: wordpress-db-auth
            image: wordpress:4.8-apache
            name: wordpress
          restartPolicy: OnFailure
  schedule: '*/10 * * * *'
---
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  labels:
    app: wordpress
  name: wordpress
spec:
  selector:
    matchLabels:
      app: wordpress
  template:
    metadata:
      labels:
        app: wordpress
    spec:
      containers:
      - env:
        - name: WORDPRESS_DB_USER
          valueFrom:
            secretKeyRef:
              key: user
              name: wordpress-db-auth
        - name: WORDPRESS_DB_PASSWORD
          valueFrom:
            secretKeyRef:
              key: password
              name: wordpress-db-auth
        image: wordpress:4.8-apache
        name: wordpress
        ports:
        - containerPort: 80
          name: wordpress
---
apiVersion: v1
kind: Values
metadata:
  name: shared
spec:
  env:
  - name: WORDPRESS_DB_USER
    valueFrom:
      secretKeyRef:
        key: user
        name: wordpress-db-auth
  - name: WORDPRESS_DB_PASSWORD
    valueFrom:
      secretKeyRef:
        key: password
        name: wordpress-db-auth
`)
}

type inlineCompositionTest struct{}

func (ut *inlineCompositionTest) writeKustFileProb(th kusttest_test.Harness) {
	th.WriteK("/inlinecomposition/probe/", `
resources:
#- ../base
- dep-patch.yaml

vars: 
- name: Deployment.probe.spec.template.spec.containers[0].livenessProbe
  objref:
    kind: Deployment
    name: probe
    apiVersion: apps/v1
  fieldref:
    fieldpath: spec.template.spec.containers[0].livenessProbe
`)
}
func (ut *inlineCompositionTest) writeKustFileComposite(th kusttest_test.Harness) {
	th.WriteK("/inlinecomposition/composite/", `
resources:
- ../base
- ../probe
- ../dns
- ../restart
`)
}
func (ut *inlineCompositionTest) writeKustFileDns(th kusttest_test.Harness) {
	th.WriteK("/inlinecomposition/dns/", `
resources:
#- ../base
- dep-patch.yaml

vars:
- name: Deployment.dns.spec.template.spec.dnsPolicy
  objref:
    kind: Deployment
    name: dns
    apiVersion: apps/v1
  fieldref:
    fieldpath: spec.template.spec.dnsPolicy
`)
}
func (ut *inlineCompositionTest) writeKustFileBase(th kusttest_test.Harness) {
	th.WriteK("/inlinecomposition/base/", `
resources:
- deployment.yaml

configurations:
- kustomizeconfig.yaml
`)
}
func (ut *inlineCompositionTest) writeKustFileRestart(th kusttest_test.Harness) {
	th.WriteK("/inlinecomposition/restart/", `
resources:
#- ../base
- dep-patch.yaml

vars:
- name: Deployment.restart.spec.template.spec.restartPolicy
  objref:
    kind: Deployment
    name: restart
    apiVersion: apps/v1
  fieldref:
    fieldpath: spec.template.spec.restartPolicy
`)
}
func (ut *inlineCompositionTest) writeKustomizeConfig(th kusttest_test.Harness) {
	th.WriteF("/inlinecomposition/base/kustomizeconfig.yaml", `
varReference:
- path: spec/template/spec/containers[]/livenessProbe
  kind: Deployment
- path: spec/template/spec/dnsPolicy
  kind: Deployment
- path: spec/template/spec/restartPolicy
  kind: Deployment
`)
}
func (ut *inlineCompositionTest) writeProbePatch(th kusttest_test.Harness) {
	th.WriteF("/inlinecomposition/probe/dep-patch.yaml", `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: probe
  namespace: patch
spec:
  template:
    spec:
      containers:
      - livenessProbe:
          httpGet:
            path: /healthz
            port: 8080
`)
}
func (ut *inlineCompositionTest) writeDnsPatch(th kusttest_test.Harness) {
	th.WriteF("/inlinecomposition/dns/dep-patch.yaml", `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: dns
  namespace: patch
spec:
  template:
    spec:
      dnsPolicy: ClusterFirst
`)
}
func (ut *inlineCompositionTest) writeBaseDeployment(th kusttest_test.Harness) {
	th.WriteF("/inlinecomposition/base/deployment.yaml", `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  template:
    spec:
      containers:
      - name: my-deployment
        image: my-image
        livenessProbe: $(Deployment.probe.spec.template.spec.containers[0].livenessProbe)
      dnsPolicy: $(Deployment.dns.spec.template.spec.dnsPolicy)
      restartPolicy: $(Deployment.restart.spec.template.spec.restartPolicy)
`)
}
func (ut *inlineCompositionTest) writeRestartPatch(th kusttest_test.Harness) {
	th.WriteF("/inlinecomposition/restart/dep-patch.yaml", `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: restart
  namespace: patch
spec:
  template:
    spec:
      restartPolicy: Always
`)
}
func TestInlineComposition(t *testing.T) {
	ut := &inlineCompositionTest{}
	th := kusttest_test.MakeHarness(t)
	ut.writeKustFileBase(th)
	ut.writeKustFileRestart(th)
	ut.writeKustFileDns(th)
	ut.writeKustFileProb(th)
	ut.writeKustFileComposite(th)
	ut.writeKustomizeConfig(th)
	ut.writeRestartPatch(th)
	ut.writeBaseDeployment(th)
	ut.writeDnsPatch(th)
	ut.writeProbePatch(th)
	m := th.Run("/inlinecomposition/composite", th.MakeDefaultOptions())
	th.AssertActualEqualsExpected(m, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  template:
    spec:
      containers:
      - image: my-image
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8080
        name: my-deployment
      dnsPolicy: ClusterFirst
      restartPolicy: Always
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: probe
  namespace: patch
spec:
  template:
    spec:
      containers:
      - livenessProbe:
          httpGet:
            path: /healthz
            port: 8080
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: dns
  namespace: patch
spec:
  template:
    spec:
      dnsPolicy: ClusterFirst
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: restart
  namespace: patch
spec:
  template:
    spec:
      restartPolicy: Always
`)
}
