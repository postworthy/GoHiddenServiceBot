FROM kvrg/torssh:1.0 AS base
WORKDIR /app

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        ca-certificates

COPY build ./
RUN chmod +x /app/go-hidden-service-bot

ENTRYPOINT ["/app/go-hidden-service-bot"]

