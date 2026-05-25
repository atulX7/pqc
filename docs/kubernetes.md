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
