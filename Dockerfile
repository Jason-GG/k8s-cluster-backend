FROM xxxxxxxxxxxxxxxxx.dkr.ecr.us-east-1.amazonaws.com/base-image:0.0.6
LABEL org.opencontainers.image.authors="Jason-GG"
VOLUME /tmp
ENV GOPROXY https://goproxy.io,direct
ENV RUN_ENV="prod"
ENV COMMIT=""
ADD config /root/config
ADD gaget /root/gaget
ADD cicd_linux /root/cicd_linux
RUN chmod +x /root/cicd_linux
RUN source /etc/bashrc
WORKDIR /root
ENTRYPOINT ./cicd_linux
EXPOSE 3000