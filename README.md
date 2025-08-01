# Kupher Registry Policy Controller
A Kubernetes Validating Admission Controller to enforce image registry policies per namespace and across cluster.

## Overview
Kupher Registry Policy Controller is a Kubernetes-native admission controller that ensures workloads only use container images from trusted registries. You can define:

✅ Common registries allowed in all namespaces

🎯 Namespace-specific registries only allowed in selected namespaces

## Use Cases
- Enforce private registry usage per tenant (platform, frontend, etc.)

- Block public image pulls from docker.io in production

- Add FinOps and security enforcement into your platform layer

## Features
- Validates Pod, Deployment, DaemonSet, StatefulSet, Job, and CronJob

- Rejects workloads using unauthorized image registries

- Supports dynamic config via mounted ConfigMap (YAML)

- Written in Go with clean modular logic

- Easy to test locally (Minikube/k3d) or deploy in production

## How It Works
Kupher Registry Policy Controller intercepts workload creation using Kubernetes' Validating Admission Webhooks.

For each workload:

    - It extracts all container images

    - It checks whether the image is:

        - From a common trusted registry, or

        - From an approved registry for the namespace

    - If not allowed, it rejects the resource creation

## Maintainer
Rajesh Deshpande
Cloud-Native Architect | Golang + Kubernetes Practitioner
Follow on LinkedIn for more Kupher tools & insights.

## Tags
#Kubernetes #AdmissionController #DevSecOps #PlatformEngineering #Golang #ZeroTrust
