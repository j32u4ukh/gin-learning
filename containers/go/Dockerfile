FROM golang:latest

COPY ./containers/go/entrypoint.sh /app/entrypoint.sh

WORKDIR /usr/src/app

RUN chmod +x /app/entrypoint.sh

CMD ["/app/entrypoint.sh"]