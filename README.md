# kubernetes-cloudflare-updater
A simple program to update Cloudflare DNS records on start and on exit of the program.
This solves the issue for users of Cloudflare that for example create a domain like: `k8s-masters.example.com` and add all nodes to it as a seperate `A`/`AAAA` record.

## Configuration
The configuration is done using environment variables:

### Environment Variables
#### Cloudflare
```
CF_API_KEY="YOUR_CLOUDFLARE_API_KEY"
CF_API_EMAIL="YOUR_CLOUDFLARE_EMAIL"
```

#### Domain Names and IP Adresses
```
DNS_NAMES="k8s-masters.example.com,k8s-workers.example.com"
IP_ADDRESS="1.1.1.1,1::"
```

### Kubernetes Manifest
See the file `bundle.yaml` in the repo root.
