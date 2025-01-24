4. Deploying and Configuring k8sgpt-operator
k8sgpt-operator can automate k8sgpt in the cluster. You can install it using Helm.

helm repo add k8sgpt https://charts.k8sgpt.ai/
helm repo update
helm install release k8sgpt/k8sgpt-operator -n k8sgpt --create-namespace
k8sgpt-operator provides two CRDs: K8sGPT to configure k8sgpt and Result to output analysis results.

kubectl api-resources  | grep -i gpt
k8sgpts                                        core.k8sgpt.ai/v1alpha1                true         K8sGPT
results                                        core.k8sgpt.ai/v1alpha1                true         Result
Configure K8sGPT, using Ollama's IP address for baseUrl.

kubectl apply -n k8sgpt -f - << EOF
apiVersion: core.k8sgpt.ai/v1alpha1
kind: K8sGPT
metadata:
  name: k8sgpt-ollama
spec:
  ai:
    enabled: true
    model: llama3
    backend: localai
    baseUrl: http://198.19.249.3:11434/v1
  noCache: false
  filters: ["Pod"]
  repository: ghcr.io/k8sgpt-ai/k8sgpt
  version: v0.3.8
EOF
After creating the K8sGPT CR, the operator will automatically create a pod for it. Checking the Result CR will show the same results.

kubectl get result -n k8sgpt -o jsonpath='{.items[].spec}' | jq .
{
  "backend": "localai",
  "details": "Error: Kubernetes is unable to pull the image \"image-not-exist\" due to it not existing.\n\nSolution: \n1. Check if the image actually exists.\n2. If not, create the image or use an alternative one.\n3. If the image does exist, ensure that the Docker daemon and registry are properly configured.",
  "error": [
    {
      "text": "Back-off pulling image \"image-not-exist\""
    }
  ],
  "kind": "Pod",
  "name": "default/k8sgpt-test",
  "parentObject": ""
}