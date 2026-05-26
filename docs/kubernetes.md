# Kubernetes Deployment

This app can run on Docker Desktop Kubernetes as a single web deployment.

Build the local image:

```bash
docker build -t pqc-readiness:dev .
```

Deploy:

```bash
kubectl apply -k k8s
kubectl -n pqc get pods,svc
```

Open the UI:

```bash
kubectl -n pqc port-forward svc/pqc-readiness 8080:8080
```

Then visit:

```text
http://localhost:8080
```

The app supports:

- `GET /healthz`
- `POST /api/scan/sample`
- `POST /api/scan/upload` with multipart field `repo` containing a `.zip`
- `POST /api/scan/git` with multipart field `repo_url` containing a public `https://github.com/owner/repo` URL and optional `branch`
- `POST /api/scan/domains` with multipart field `domains` containing one or more TLS domains separated by newlines, commas, or semicolons
