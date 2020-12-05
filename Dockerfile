#
# Контейнер сборки
#
FROM golang:1.14 as builder

ARG DRONE
ARG DRONE_TAG
ARG DRONE_COMMIT
ARG DRONE_BRANCH

ENV CGO_ENABLED=0

COPY . /go/src/github.com/KubaiDoLove/conductor
WORKDIR /go/src/github.com/KubaiDoLove/conductor
RUN \
    if [ -z "$DRONE" ] ; then echo "no drone" && version=`git describe --abbrev=6 --always --tag`; \
    else version=${DRONE_TAG}${DRONE_BRANCH}-`echo ${DRONE_COMMIT} | cut -c 1-7` ; fi && \
    echo "version=$version" && \
    cd cmd/conductor && \
    go build -a -tags conductor -installsuffix conductor -ldflags "-X apiserver.version=${version} -s -w" -o /go/bin/conductor

#
# Контейнер для получения актуальных SSL/TLS сертификатов
#
FROM alpine as alpine
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
RUN addgroup -S conductor && adduser -S conductor -G conductor

ENTRYPOINT [ "/bin/conductor" ]

#
# Контейнер рантайма
#
FROM scratch
COPY --from=builder /go/bin/conductor /bin/conductor

# копируем сертификаты из alpine
COPY --from=alpine /etc/ssl/certs /etc/ssl/certs

# копируем документацию
COPY --from=alpine /usr/share/conductor /usr/share/conductor

# копируем пользователя и группу из alpine
COPY --from=alpine /etc/passwd /etc/passwd
COPY --from=alpine /etc/group /etc/group

USER conductor

ENTRYPOINT ["/bin/conductor"]



