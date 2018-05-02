#!/bin/bash -e
scriptdir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cd ~/go/src/k8s.io/code-generator && ./generate-groups.sh \
  all \
  github.com/srleyva/aws-operator/pkg/client \
  github.com/srleyva/aws-operator/pkg/apis \
  "sleyva:v1alpha1" \
