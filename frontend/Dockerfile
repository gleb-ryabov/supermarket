FROM httpd:2.4-alpine

COPY . /usr/local/apache2/htdocs/
COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

EXPOSE 80

CMD ["/entrypoint.sh"]