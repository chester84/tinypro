FROM alpine:latest

#RUN sed -i -e 's/http:/https:/' /etc/apk/repositories
#RUN apk add curl

LABEL maintainer="yctech tinypro backend service"

ENV TIME_ZONE=Asia/Shanghai

#RUN  echo http://mirrors.ustc.edu.cn/alpine/v3.14/main > /etc/apk/repositories && \
#  echo http://mirrors.ustc.edu.cn/alpine/v3.14/community >> /etc/apk/repositories


RUN apk add tzdata \
    && cp /usr/share/zoneinfo/$TIME_ZONE /etc/localtime \
    && echo "$TIME_ZONE" > /etc/timezone \
    && apk del tzdata

COPY build /app/

EXPOSE 8976

WORKDIR /app

CMD ["./api"]
