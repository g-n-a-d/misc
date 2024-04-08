# Victoriametrics

## Components

### vminsert

**vminsert** serves as the entry point for data ingestion, accepting data from various sources and storing it efficiently for subsequent querying and analysis.

### vmselect

**vmselect** responsible for handling incoming queries and retrieving the relevant time-series data from the storage layer.

### vmstorage

**vmstorage** responsible for storing and managing time-series data efficiently.

## Architecture

![arch](https://docs.victoriametrics.com/Cluster-VictoriaMetrics_cluster-scheme.webp)

## Important Features

### Data Availability

The replication factor is a crucial configuration parameter for ensuring data durability and high availability in distributed storage systems. Which represents number of copies of data stored across the cluster for redundancy and fault tolerance. Specifically, it determines how many replicas of each data shard are maintained in the cluster.

Data availability can be ensured by passing **-replicationFactor=N** as an keyword argument to **vminsert**.

### Deduplication

Deduplication leaves a single raw sample with the biggest timestamp for each time series per each discrete interval. It is a common technique used to reduce storage space, improve data efficiency, and streamline data processing operations. This is similar to the staleness rules in Prometheus

Set **-dedup.minScrapeInterval** to define the discrete interval. The recommended value for **-dedup.minScrapeInterval** must equal to **scrape_interval** config from Prometheus configs. It is recommended to have a single scrape_interval across all the scrape targets.

**vmselect** must run with **-dedup.minScrapeInterval=1ms** for data de-duplication when -**replicationFactor** is greater than 1. Higher values for **-dedup.minScrapeInterval** at **vmselect** is OK.

### Downsampling

Downsampling in VictoriaMetrics refers to the process of aggregating or reducing the granularity of time-series data over a specific time window. Downsampling is commonly used to manage large volumes of time-series data efficiently, especially in scenarios where high-resolution data is not required for long-term analysis or visualization.

- VictoriaMetrics supports various aggregation functions for downsampling.

- Downsampling in VictoriaMetrics is performed over configurable time windows, known as "aggregation intervals" or "step intervals". Users can specify the duration of the aggregation window by setting **-downsampling.period**.

- Downsampling helps optimize storage utilization by reducing the volume of raw time-series data stored in the database.

The **-dedup.minScrapeInterval=D** is equivalent to **-downsampling.period=0s:D** if downsampling is enabled.

*Sample docker-compose.yml*

```
version: '3.5'
services:
  #  Metrics collector.
  #  It scrapes targets defined in --promscrape.config
  #  And forward them to --remoteWrite.url
  vmagent:
    container_name: vmagent
    image: victoriametrics/vmagent:v1.100.0
    depends_on:
      - "vminsert"
    ports:
      - 8429:8429
    volumes:
      - vmagentdata:/vmagentdata
      - ./prometheus-cluster.yml:/etc/prometheus/prometheus.yml
    command:
      - '--promscrape.config=/etc/prometheus/prometheus.yml'
      - '--remoteWrite.url=http://vminsert:8480/insert/0/prometheus/'
    restart: always

  # vmstorage shards. Each shard receives 1/N of all metrics sent to vminserts,
  # where N is number of vmstorages (2 in this case).
  vmstorage-1:
    container_name: vmstorage-1
    image: victoriametrics/vmstorage:v1.100.0-cluster
    ports:
      - 8482
      - 8400
      - 8401
    volumes:
      - strgdata-1:/storage
    command:
      - '--storageDataPath=/storage'
    restart: always
  vmstorage-2:
    container_name: vmstorage-2
    image: victoriametrics/vmstorage:v1.100.0-cluster
    ports:
      - 8482
      - 8400
      - 8401
    volumes:
      - strgdata-2:/storage
    command:
      - '--storageDataPath=/storage'
    restart: always

  # vminsert is ingestion frontend. It receives metrics pushed by vmagent,
  # pre-process them and distributes across configured vmstorage shards.
  vminsert:
    container_name: vminsert
    image: victoriametrics/vminsert:v1.100.0-cluster
    depends_on:
      - "vmstorage-1"
      - "vmstorage-2"
    command:
      - '--storageNode=vmstorage-1:8400'
      - '--storageNode=vmstorage-2:8400'
    ports:
      - 8480:8480
    restart: always

  # vmselect is a query fronted. It serves read queries in MetricsQL or PromQL.
  # vmselect collects results from configured `--storageNode` shards.
  vmselect-1:
    container_name: vmselect-1
    image: victoriametrics/vmselect:v1.100.0-cluster
    depends_on:
      - "vmstorage-1"
      - "vmstorage-2"
    command:
      - '--storageNode=vmstorage-1:8401'
      - '--storageNode=vmstorage-2:8401'
      - '--vmalert.proxyURL=http://vmalert:8880'
    ports:
      - 8481
    restart: always
  vmselect-2:
    container_name: vmselect-2
    image: victoriametrics/vmselect:v1.100.0-cluster
    depends_on:
      - "vmstorage-1"
      - "vmstorage-2"
    command:
      - '--storageNode=vmstorage-1:8401'
      - '--storageNode=vmstorage-2:8401'
      - '--vmalert.proxyURL=http://vmalert:8880'
    ports:
      - 8481
    restart: always

  # vmauth is a router and balancer for HTTP requests.
  # It is configured via --auth.config and balances
  # read requests from Grafana, vmui, vmalert among vmselects.
  # It can be used as an authentication proxy.
  vmauth:
    container_name: vmauth
    image: victoriametrics/vmauth:v1.100.0
    depends_on:
      - "vmselect-1"
      - "vmselect-2"
    volumes:
      - ./auth-cluster.yml:/etc/auth.yml
    command:
      - '--auth.config=/etc/auth.yml'
    ports:
      - 8427:8427
    restart: always

  # vmalert executes alerting and recording rules
  vmalert:
    container_name: vmalert
    image: victoriametrics/vmalert:v1.100.0
    depends_on:
      - "vmauth"
    ports:
      - 8880:8880
    volumes:
      - ./alerts-cluster.yml:/etc/alerts/alerts.yml
      - ./alerts-health.yml:/etc/alerts/alerts-health.yml
      - ./alerts-vmagent.yml:/etc/alerts/alerts-vmagent.yml
      - ./alerts-vmalert.yml:/etc/alerts/alerts-vmalert.yml
    command:
      - '--datasource.url=http://vmauth:8427/select/0/prometheus'
      - '--remoteRead.url=http://vmauth:8427/select/0/prometheus'
      - '--remoteWrite.url=http://vminsert:8480/insert/0/prometheus'
      - '--notifier.url=http://alertmanager:9093/'
      - '--rule=/etc/alerts/*.yml'
      # display source of alerts in grafana
      - '-external.url=http://127.0.0.1:3000' #grafana outside container
      # when copypaste the line below be aware of '$$' for escaping in '$expr'
      - '--external.alert.source=explore?orgId=1&left={"datasource":"VictoriaMetrics","queries":[{"expr":{{$$expr|jsonEscape|queryEscape}},"refId":"A"}],"range":{"from":"now-1h","to":"now"}}'
    restart: always

  # alertmanager receives alerting notifications from vmalert
  # and distributes them according to --config.file.
  alertmanager:
    container_name: alertmanager
    image: prom/alertmanager:v0.27.0
    volumes:
      - ./alertmanager.yml:/config/alertmanager.yml
    command:
      - '--config.file=/config/alertmanager.yml'
    ports:
      - 9093:9093
    restart: always

volumes:
  vmagentdata: {}
  strgdata-1: {}
  strgdata-2: {}
```
