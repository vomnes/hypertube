FROM golang
EXPOSE 4003
COPY ./App/ /app/
RUN cd /app/api/ && make vendor_get
RUN cd /app/api/ && make build
RUN apt-get update --fix-missing && apt-get install cron -y --fix-missing
ADD ./crontab /etc/cron.d/update-db-cron
RUN crontab /etc/cron.d/update-db-cron
ENTRYPOINT cron && cd /app/api/ && make launch	