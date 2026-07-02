#!/bin/sh

echo "window.APP_CONFIG = { API_BASE_URL: \"$API_BASE_URL\" };" > /usr/local/apache2/htdocs/js/config.js

httpd-foreground