FROM google/cloud-sdk:221.0.0-alpine

LABEL maintainer="estafette.io" \
      description="The estafette-extension-gke component is an Estafette extension to run commands against Kubernetes Engine"

RUN gcloud components install kubectl \
    && rm -rf /var/cache/apk/*

COPY estafette-extension-gke /
COPY templates /templates

RUN du -h / -d 1

ENTRYPOINT ["/estafette-extension-gke"]