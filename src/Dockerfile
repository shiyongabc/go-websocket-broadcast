FROM alpine
ENV TZ=Asia/Shanghai

ENV API_HOST_LS :80
RUN apk add --no-cache tzdata && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
COPY build/push_service /push_service

RUN chmod +x /push_service
EXPOSE 80

CMD ["/push_service"]
