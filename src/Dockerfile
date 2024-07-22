FROM alpine
ENV TZ=Asia/Shanghai

ENV API_HOST_LS :80
RUN apk add --no-cache tzdata && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
#COPY push_service-linux-amd64 /push_service

RUN chmod +x push_service-linux-amd64
EXPOSE 80

CMD ["push_service-linux-amd64 pro"]
