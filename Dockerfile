FROM quay.io/prometheus/busybox:latest

LABEL maintainer="Alexander Trost <galexrt@googlemail.com>"

ENV ARCH="linux_amd64"

ADD output/kubernetes-cloudflare-updater_$ARCH /bin/kubernetes-cloudflare-updater

ENTRYPOINT ["/bin/kubernetes-cloudflare-updater"]
