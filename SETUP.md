# Setup Guide for ArgoCD GitOps Repository

This guide provides step-by-step instructions to set up and run the ArgoCD GitOps repository with sample API and web services.

## Prerequisites

- A GitHub account with access to GitHub Container Registry (GHCR)
- A Kubernetes cluster (local or cloud-based)
- ArgoCD installed on your cluster
- kubectl configured to access your cluster
- Git installed locally

## Step 1: Fork and Clone the Repository

1. Fork this repository to your GitHub account
2. Clone your forked repository:
   ```bash
   git clone https://github.com/eknathdj/argocd-mcp-gitops.git
   cd argocd-mcp-gitops
   ```

## Step 2: Configure Repository-Specific Settings

1. **Update ArgoCD Application URL**
   Edit `argocd/application.yaml` and replace the repoURL with your fork's URL:
   ```yaml
   source:
     repoURL: https://github.com/YOUR_USERNAME/argocd-mcp-gitops.git
   ```

2. **Verify Owner Configuration**
   The repository is already configured for owner `eknathdj`. If you're using a different GitHub username, update:
   - `deploy/base/api-deployment.yaml`: Change `ghcr.io/eknathdj/api:latest`
   - `deploy/base/kustomization.yaml`: Change image names to use your username
   - CI/CD workflow will automatically use your repository owner

## Step 3: Set Up GitHub Container Registry (GHCR)

1. Ensure GHCR is enabled for your GitHub account
2. The CI/CD workflow uses `GITHUB_TOKEN` for authentication, which should work by default
3. If your organization requires a Personal Access Token (PAT):
   - Create a PAT with `packages:write` scope
   - Add it as `GHCR_TOKEN` secret in your repository settings
   - Update the workflow's login step to use the PAT instead of `GITHUB_TOKEN`

## Step 4: Install ArgoCD on Your Cluster

```bash
# Install ArgoCD CLI (optional but recommended)
curl -sSL -o argocd-linux-amd64 https://github.com/argoproj/argo-cd/releases/latest/download/argocd-linux-amd64
chmod +x argocd-linux-amd64
sudo mv argocd-linux-amd64 /usr/local/bin/argocd

# Install ArgoCD on your cluster
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# Wait for ArgoCD to be ready
kubectl wait --for=condition=available --timeout=300s deployment/argocd-server -n argocd
```

## Step 5: Access ArgoCD Web UI

```bash
# Get initial admin password
argocd admin initial-password -n argocd

# Port forward ArgoCD server (in a separate terminal)
kubectl port-forward svc/argocd-server -n argocd 8080:80

# Access ArgoCD at https://localhost:8080
# Login with username: admin and the password from above
```

## Step 6: Deploy the Application

1. **Apply the ArgoCD Application**
   ```bash
   kubectl apply -n argocd -f argocd/application.yaml
   ```

2. **Trigger the CI/CD Pipeline**
   - Make a change to any file under `services/` directory
   - Commit and push the changes
   - The GitHub Actions workflow will:
     - Build Docker images for changed services
     - Push images to GHCR
     - Update the image tags in `deploy/base/kustomization.yaml`
     - ArgoCD will automatically sync the changes

## Step 7: Set Up Ingress Controller (NGINX)

Ensure NGINX Ingress Controller is installed on your cluster:

```bash
# For local clusters (like Minikube, Kind)
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.10.1/deploy/static/provider/cloud/deploy.yaml

# For cloud providers, follow their specific ingress setup
```

## Step 8: Access the Application

After deployment:

- **Production**: Access at `http://app.127.0.0.1.nip.io`
- **Development**: Access at `http://app.dev.127.0.0.1.nip.io`

The web application will display:
- Health status from the API
- Current server time from the API
- Auto-refresh every 5 seconds

## Step 9: Optional - Set Up MCP for ArgoCD

Create `.vscode/mcp.json` (or equivalent for Cursor) for IDE integration:

```json
{
  "servers": {
    "argocd-mcp-stdio": {
      "type": "stdio",
      "command": "npx",
      "args": ["argocd-mcp@latest", "stdio"],
      "env": {
        "ARGOCD_BASE_URL": "https://your-argocd-server-url",
        "ARGOCD_API_TOKEN": "your-api-token"
      }
    }
  }
}
```

## Troubleshooting

### Common Issues

1. **Images not pulling**: Ensure GHCR access and correct image names
2. **Ingress not working**: Check ingress controller installation and cluster networking
3. **ArgoCD sync issues**: Verify repository URL and access permissions
4. **Build failures**: Check GitHub Actions logs and ensure Dockerfile syntax

### Useful Commands

```bash
# Check ArgoCD application status
argocd app get app -n argocd

# Force sync
argocd app sync app -n argocd

# View application logs
kubectl logs -n app deployment/api
kubectl logs -n app deployment/web

# Check ingress
kubectl get ingress -n app
kubectl describe ingress app -n app
```

## Adding New Services

1. Create `services/<new-service>/Dockerfile` and source code
2. Add Deployment/Service manifests in `deploy/base/`
3. Update `deploy/base/kustomization.yaml` to include the new image
4. Commit changes to trigger CI/CD

## Environment Variables

The setup supports different overlays:
- `deploy/overlays/prod`: Production environment
- `deploy/overlays/dev`: Development environment

To deploy to dev environment, update the ArgoCD application path to `deploy/overlays/dev`.