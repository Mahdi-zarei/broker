FROM envoyproxy/envoy-dev:dc42d9817bdd2285b3f1c934a0e203af5db8bc66
COPY config.yaml /etc/envoy/envoy.yaml

EXPOSE 9595
EXPOSE 9901

RUN chmod go+r /etc/envoy/envoy.yaml