FROM scratch

COPY sample-device-plugin /usr/bin/sample-device-plugin
ENTRYPOINT ["/usr/bin/sample-device-plugin"]
