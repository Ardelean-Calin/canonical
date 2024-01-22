FROM ubuntu:22.04

RUN apt-get update && apt-get install -y sudo
RUN mkdir -p /tmp/calinos

# COPY run.sh /tmp/calinos/
# COPY config /tmp/calinos/config
# COPY scripts /tmp/calinos/scripts
# RUN chmod +x /tmp/calinos/run.sh

WORKDIR /tmp/calinos
CMD ["/tmp/calinos/run.sh"]
