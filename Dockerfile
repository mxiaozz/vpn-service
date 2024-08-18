FROM debian:12.6-slim
WORKDIR /app
CMD ["vpn-web"]
ENV PATH=/app:$PATH
ADD deploy/ /app