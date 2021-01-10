FROM kvrg/torssh:1.0 AS base
WORKDIR /app

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        ca-certificates

COPY ./login-notify.sh /etc/profile.d/login-notify.sh
COPY build/linux/amd64/go-hidden-service-bot ./
RUN chmod +x /etc/profile.d/login-notify.sh
RUN chmod +x /app/go-hidden-service-bot

ENTRYPOINT ["/app/go-hidden-service-bot"]

