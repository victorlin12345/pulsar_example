# pulsar_example

## quick start
### setup pulsar docker
```
docker run -it -d \
-p 6650:6650 \
-p 8080:8080 \
--name pulsar \
--mount source=pulsardata,target=/pulsar/data \
--mount source=pulsarconf,target=/pulsar/conf \
apachepulsar/pulsar:4.0.2 \
bin/pulsar standalone
```
### create tanent
```
bin/pulsar-admin tenants create --admin-roles admin --allowed-clusters standalone investments
```
### create namespace
```
bin/pulsar-admin namespaces create investments/stocks
```
### create topic
```
bin/pulsar-admin topics create persistent://investments/stocks/stock-ticker
```
### run up consumer
```
go run main.go consumer
```
### run up producer
```
go run main.go producer
```

## TBD
- [x] setup pulsar docker
- [x] create tanent, namespace, topic
- [x] implement producer
- [x] implement consumer
- [ ] rerange code layout

## Reference
- https://www.youtube.com/watch?v=Qzq52ADcBD8
- https://github.com/apache/pulsar-client-go
- https://pulsar.apache.org/docs/4.0.x/getting-started-docker/
- https://github.com/spf13/cobra/blob/main/site/content/user_guide.md
- https://github.com/spf13/cobra-cli/blob/main/README.md