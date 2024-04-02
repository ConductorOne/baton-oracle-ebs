FROM gcr.io/distroless/static-debian11:nonroot
ENTRYPOINT ["/baton-oracle-ebs"]
COPY baton-oracle-ebs /