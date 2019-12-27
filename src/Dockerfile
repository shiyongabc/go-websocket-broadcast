FROM alpine
ENV TZ=Asia/Shanghai

ENV API_HOST_LS :80
RUN apk add --no-cache tzdata && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
COPY build/service_push /service_push

RUN chmod +x /service_push
EXPOSE 80

CMD ["/service_push"]
