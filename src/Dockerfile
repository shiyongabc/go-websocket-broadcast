FROM alpine
ENV TZ=Asia/Shanghai
ENV PRO_CONFIG=pro
ENV API_HOST_LS :80
RUN apk add --no-cache tzdata && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
COPY push_service-linux-amd64 /push_service
COPY config.pro.json /src/config.pro.json
#RUN ls /
RUN chmod +x /push_service
EXPOSE 80

#CMD ["/push_service"]
ENTRYPOINT ["/push_service","$PRO_CONFIG"]
