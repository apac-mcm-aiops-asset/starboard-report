# IBM_SOURCE_PROLOG_BEGIN_TAG
# *****************************************************************
#
# IBM Confidential
# OCO Source Materials
#
# (C) Copyright IBM Corp. 2020
#
# The source code for this program is not published or otherwise
# divested of its trade secrets, irrespective of what has been
# deposited with the U.S. Copyright Office.
#
# *****************************************************************
# IBM_SOURCE_PROLOG_END_TAG
# Multi Stage Build

#### Stage 1 - Node build
# On ppc64le, Build fails to download https://github.com/sass/node-sass/releases
# because it is not supported for the architecture, so ppc will just build node-saas.
FROM node:16.13-stretch as builder

RUN apt-get update && apt-get install -y \
    python \
    make \
    g++ \
    && apt-get clean

# The base node image sets a very verbose log level.
ENV NPM_CONFIG_LOGLEVEL warn

# Setup workdir
WORKDIR /project
ENV WORKDIR /project

# Install app dependencies
# A wildcard is used to ensure both package.json AND package-lock.json are copied
# where available
COPY package*.json ./
# Actions to eliminate possible security issues
# RUN npm i npm@latest -g && npm audit fix
RUN npm install

# Bundle app source after npm install to prevent layer changes
COPY . .

# Build node app. Output goes into dist/ directory
RUN npm run build

#### Stage 2 - Nginx build
FROM nginxinc/nginx-unprivileged:1.16.1-alpine

# Setup workdir
WORKDIR /project
ENV WORKDIR /project

# Place static UI code into nginx default location
ENV STATIC_DIR /usr/share/nginx/html
COPY --from=builder $WORKDIR/dist $STATIC_DIR

USER 65532:65532