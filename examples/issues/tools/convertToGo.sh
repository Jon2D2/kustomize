#!/bin/bash
ISSUENUMBER=1295
ISSUENAME=issue_${ISSUENUMBER}
NEWFILE=${ISSUENAME}_test.go
APPDIR=recreate${ISSUENUMBER}
ALLINONEDIR=$HOME/src/sigs.k8s.io/kustomize/examples/issues/${ISSUENAME}
ISSUEDIR=$ALLINONEDIR
KUSTFILES=`find $ISSUEDIR -name kustomization.yaml -print`
CONFIGFILES=`find $ISSUEDIR -name kustomizeconfig.yaml -print`
RESOURCES=`find $ISSUEDIR -name "*.yaml" -print | grep -v expected | grep -v kustomization.yaml | grep -v kustomizeconfig.yaml`

cat <<EOF >$NEWFILE
// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package target_test

import (
	kusttest_test "sigs.k8s.io/kustomize/v3/pkg/kusttest"
	"testing"
)

type ${APPDIR}Test struct {}
EOF

KUSTCOUNTER=0
for i in $KUSTFILES
do
echo 'func (ut *'$APPDIR'Test) writeKustFile'$KUSTCOUNTER'(th *kusttest_test.KustTestHarness) {' >> $NEWFILE
KUSTDIR=`echo $i | sed -e "s;$ISSUEDIR;;g" -e "s/kustomization.yaml//g"`
echo 'th.WriteK("/'${APPDIR}${KUSTDIR}'", `' >> $NEWFILE
cat $i >> $NEWFILE
echo '`)' >> $NEWFILE
echo '}' >> $NEWFILE
KUSTCOUNTER=$((KUSTCOUNTER +1))
done


CONFIGCOUNTER=0
for i in $CONFIGFILES
do
FILENAME=`echo $i | sed -e "s;$ISSUEDIR;;g"`
echo 'func (ut *'$APPDIR'Test) writeConfigFile'$CONFIGCOUNTER'(th *kusttest_test.KustTestHarness) {' >> $NEWFILE
echo 'th.WriteF("/'${APPDIR}${FILENAME}'", `' >> $NEWFILE
cat $i >> $NEWFILE
echo '`)' >> $NEWFILE
echo '}' >> $NEWFILE
CONFIGCOUNTER=$((CONFIGCOUNTER +1))
done

RESOURCECOUNTER=0
for i in $RESOURCES
do
FILENAME=`echo $i | sed -e "s;$ISSUEDIR;;g"`
echo 'func (ut *'$APPDIR'Test) writeResourceFile'$RESOURCECOUNTER'(th *kusttest_test.KustTestHarness) {' >> $NEWFILE
echo 'th.WriteF("/'${APPDIR}${FILENAME}'", `' >> $NEWFILE
cat $i >> $NEWFILE
echo '`)' >> $NEWFILE
echo '}' >> $NEWFILE
RESOURCECOUNTER=$((RESOURCECOUNTER +1))
done

cat <<EOF >>$NEWFILE
func TestIssue${ISSUENUMBER}(t *testing.T) {
        ut := &${APPDIR}Test{}
	th := kusttest_test.NewKustTestHarness(t, "/${APPDIR}/myapp")
EOF

x=$KUSTCOUNTER
while ((x--)); do
echo 'ut.writeKustFile'$x'(th)' >> $NEWFILE
done

x=$CONFIGCOUNTER
while ((x--)); do
echo 'ut.writeConfigFile'$x'(th)' >> $NEWFILE
done

x=$RESOURCECOUNTER
while ((x--)); do
echo 'ut.writeResourceFile'$x'(th)' >> $NEWFILE
done

cat <<EOF >>$NEWFILE
	m, err := th.MakeKustTarget().MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	th.AssertActualEqualsExpected(m, \`
EOF

kustomize build $ISSUEDIR/myapp >> $NEWFILE

cat <<EOF >>$NEWFILE
\`)
}
EOF

go fmt $NEWFILE
go test -v $NEWFILE
