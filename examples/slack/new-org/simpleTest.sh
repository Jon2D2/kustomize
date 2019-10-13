#!/bin/bash
for i in `cat serviceList`
do
	kustomize build $i | kubectl apply -f -
	sleep 10
	kustomize build $i | kubectl delete -f -
done
