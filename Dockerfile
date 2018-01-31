ARG arch
FROM multiarch/alpine:${arch}-edge

COPY ./expino-ticker /expino-ticker

ENV INFLUXURL=""

CMD /expino-ticker