#!/bin/bash

echo "Enter your cluster name"
read cluster

echo "Enter pool name"
read pool

set -x

gcloud beta container node-pools create \
    --node-labels image=cos,node-type=transient,cloud.google.com/gke-preemptible=true \
    --tags=production,preemptible \
    --num-nodes=1 \
    --enable-autoscaling \
    --min-nodes=1 \
    --max-nodes=20 \
    --disk-size=80 \
    --image-type=COS \
    --machine-type=n1-standard-16 \
    --scopes=compute-rw,storage-rw \
    --cluster="$cluster" "$pool"

