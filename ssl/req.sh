#!/bin/bash

cfssl genkey csr.json | cfssljson -bare server
KUBE_CSR=$(cat kube-csr.json)
CSR_B64=$(cat server.csr | base64 -w0)
echo "${KUBE_CSR}" | jq --arg csr "${CSR_B64}" '.spec.request |= $csr' | yq eval -P | kubectl apply -f -
sleep 10
kubectl get certificaterequests.cert-manager.io mutate-ingress-for-imported-service -o json | jq -r '.status.certificate|@base64d' > server.pem
kubectl create secret -n nginx-ingress
rm *.pem *.csr