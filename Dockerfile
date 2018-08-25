FROM node:8-alpine
MAINTAINER andygertjejansen@gmail.com
EXPOSE 8080
CMD ./listener
ADD listener .