apiVersion: admissionregistration.k8s.io/v1
kind: ValidatingWebhookConfiguration
metadata:
  name: kupher-registry-policy-controller
webhooks:
  - name: registry-policy-controller.kupher.io
    rules:
      - apiGroups: ["*"]
        apiVersions: ["*"]
        operations: ["CREATE", "UPDATE"]
        resources: ["*"]
        scope: "*"
    namespaceSelector:
      matchExpressions:
      - key: name
        operator: NotIn
        values: ["kube-system", "kube-public", "kube-node-lease", "default"]
    failurePolicy: Ignore
    clientConfig:
      service:
        name: validate-registry
        namespace: default
        path: /validate-registry
        port: 443
      caBundle: ${CA_BUNDLE}
    admissionReviewVersions: ["v1"]
    sideEffects: None