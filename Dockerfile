##############################################################
## Stage 1 - Go Build
##############################################################

FROM golang:1.23 AS builder
WORKDIR /opt
COPY . .
RUN go build -o app .

#############################################################
## Stage 2 - Application Setup
##############################################################
FROM ubuntu:22.04
ARG UID=0
ARG GID=0
ARG USER=root
ARG GROUP=root
#-------------------------------------------------------------
RUN apt-get update
RUN apt-get install ca-certificates -y
#-------------------------------------------------------------
# User Set
RUN if [ "${USER}" != "root" ]; then \
        groupadd -g ${GID} ${GROUP} && \
        useradd -m -u ${UID} -g ${GID} -s /bin/bash ${USER}; \
    fi
#-------------------------------------------------------------
# Copy App and Web
RUN mkdir -p /app/log
RUN chown -R ${USER}:${GROUP} /app
USER ${USER}
COPY --from=builder --chown=${USER}:${GROUP} /opt/app /app/app
COPY --from=builder --chown=${USER}:${GROUP} /opt/web /app/web
#-------------------------------------------------------------
# Add entrypoint script
COPY --chown=${USER}:${GROUP} entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh
#-------------------------------------------------------------
WORKDIR /app
EXPOSE 3300
ENTRYPOINT ["/app/entrypoint.sh"]