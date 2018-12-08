FROM git.harknet.ml:5005/becrespi/hypertube/apache2
ARG FRONT_DOMAIN_NAME
# RUN apt-get update && apt-get install apache2 -y
RUN a2enmod proxy_wstunnel proxy proxy_http headers
EXPOSE 80
COPY ./App/virtualhost.conf /etc/apache2/sites-available/000-default.conf
COPY ./App/front /var/www/html/
RUN sed -i "1iDefine FRONT_DOMAIN_NAME $FRONT_DOMAIN_NAME" /etc/apache2/sites-available/000-default.conf
RUN echo "text/vtt                                        vtt" >> /etc/mime.types
ENTRYPOINT service apache2 start && bash