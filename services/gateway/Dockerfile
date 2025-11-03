FROM caddy:2.10.2

# Copy Caddy configuration
COPY Caddyfile /etc/caddy/Caddyfile

# Create necessary directories
RUN mkdir -p /var/log/caddy

# Copy HTML files
COPY html /usr/share/caddy

EXPOSE 80 443

CMD ["caddy", "run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"]
