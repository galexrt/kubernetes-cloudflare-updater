FROM quay.io/prometheus/busybox:latest
LABEL maintainer="Alexander Trost <galexrt@googlemail.com>"

ADD .build/linux-amd64/kubernetes-cloudflare-updater /bin/kubernetes-cloudflare-updater

ENTRYPOINT ["/bin/kubernetes-cloudflare-updater"]
