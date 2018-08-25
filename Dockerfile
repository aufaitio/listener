FROM ubuntu:18.04
MAINTAINER andygertjejansen@gmail.com
EXPOSE 8080
RUN mkdir -p /opt/app/config
RUN apt-get update && apt-get install -y nodejs npm
WORKDIR /opt/app
COPY builds/linux/listener /opt/app
COPY config/* /opt/app/config/
CMD ["/opt/app/listener"]
