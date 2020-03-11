FROM scratch
COPY alertmanager-webhook-example.linux /alertmanager-webhook-example
ENTRYPOINT [ "/alertmanager-webhook-example" ]
