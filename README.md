# goddns

A DDNS tool based for home lab.

```
docker run -d --name=goddns --restart=always -p 8001:8001 -v /path/to/config.json:/goddns/config.json baqihuanxiong/goddns:latest
```

## Example Config

```
{
  "ddns": {
    "checkInterval": "10s"
  },
  "api": {
    "listen": ":8001",
    "allowInsecureHttp": false
  },
  "orm": {
    "driver": "sqlite",
    "dsn": "goddns.db",
    "adminPassword": "admin",
    "ipScanners": [
      {
        "provider": "jsonip",
        "jsonipScanner": {
          "url": "https://jsonip.com"
        }
      }
    ],
    "dnsProviders": [
      {
        "provider": "alidns",
        "alidns": {
          "region": "cn-shanghai",
          "accessKey": "your accessKey",
          "accessKeySecret": "your accessKeySecret"
        }
      }
    ]
  },
  "metrics": {
    "enabled": true,
    "listen": ":9002",
    "metricsPath": "/metrics",
    "healthCheckPath": "/ping"
  }
}
```