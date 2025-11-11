# Single-Repo Sample App with Argo CD + GHCR + MCP

This repo demonstrates:
- Single repo for **app code + k8s manifests (kustomize)**
- GitHub Actions builds changed services and pushes images to **GHCR**
- Workflow updates image tags in **this repo**; Argo CD syncs
- Argo CD `Application` points to `deploy/overlays/prod`
- Optional MCP server for Argo CD (use in your IDE)

## Quick start

1. **Replace placeholders**
    - In `argocd/application.yaml`, change `https://github.com/eknathdj/argocd-mcdp-gitops.git` to your repo URL.

2. **Push to GitHub** (branch `main`).

3. **Ensure GHCR works**  
   By default the workflow logs into GHCR with `${{ secrets.GITHUB_TOKEN }}` which is usually sufficient. If your org requires a PAT, create `GHCR_TOKEN` with `packages:write` scope and change the login step.

4. **Install Argo CD** on your cluster (if not already).

5. **Apply the Argo CD Application**
   ```bash
   kubectl apply -n argocd -f argocd/application.yaml
   ```

6. **Trigger the workflow** by committing a change under `services/*`.
   Argo CD will pick up `deploy/overlays/prod` and deploy to namespace `app`.

7. **Access the app** at `http://app.127.0.0.1.nip.io` (Ingress class: `nginx`).

## Dev and Prod overlays

- `deploy/overlays/dev` sets host `app.dev.127.0.0.1.nip.io`
- `deploy/overlays/prod` sets host `app.127.0.0.1.nip.io`

## MCP for Argo CD (optional)

Create `.vscode/mcp.json` (or Cursor equivalent):
```json
{
  "servers": {
    "argocd-mcp-stdio": {
      "type": "stdio",
      "command": "npx",
      "args": ["argocd-mcp@latest", "stdio"],
      "env": {
        "ARGOCD_BASE_URL": "https://your-argocd.example.com",
        "ARGOCD_API_TOKEN": "REDACTED"
      }
    }
  }
}
```

## Add a new service

1. Create `services/<name>/Dockerfile` and code.
2. Add Deployment/Service in `deploy/base` and include the image in `deploy/base/kustomization.yaml`.
3. Commit — the CI will build it when changes occur under `services/<name>/`.
