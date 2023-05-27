#!/bin/bash
set -e

helm repo add jetstack https://charts.jetstack.io
helm upgrade --install cert-manager --create-namespace --namespace cert-manager jetstack/cert-manager --atomic --set installCRDs=true
