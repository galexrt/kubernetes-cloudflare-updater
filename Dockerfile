FROM quay.io/prometheus/busybox:latest

ENV ARCH="linux_amd64"

ADD output/kubernetes-cloudflare-updater_$ARCH /bin/kubernetes-cloudflare-updater

ENTRYPOINT ["/bin/kubernetes-cloudflare-updater"]
