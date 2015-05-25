# env CGO_ENABLED=0 go build -v -a -tags leveldb,netgo -installsuffix netgo,cgo
FROM neowaylabs/neosearch-dev-env:latest

ADD ./neosearch-import /neosearch-import

VOLUME ["/data"]

CMD ["./neosearch-import"]