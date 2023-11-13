
FROM envoyproxy/envoy-dev:latest

# Copy the Envoy configuration file into the container
COPY ./envoy.yaml /etc/envoy/envoy.yaml

# Ensure the Envoy process can read the configuration file
RUN chmod go+r /etc/envoy/envoy.yaml

# Run Envoy with the configuration file, logging set to trace level
# Logs are written to stdout/stderr for Cloud Run to capture
CMD /usr/local/bin/envoy -c /etc/envoy/envoy.yaml -l trace