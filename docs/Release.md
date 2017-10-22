# Upstream Release Process

## Tag dependent apbs

- keycloak https://github.com/feedhenry/keycloak-apb
- 3scale https://github.com/feedhenry/3scale-apb
- sync https://github.com/feedhenry/fh-sync-server-apb
- digger (mobile ci cd) https://github.com/feedhenry/aerogear-digger-apb  


Run the following command in each of the above repos:

```bash
TAG=0.0.4 #example version
git checkout $(TAG)
make build_and_push TAG=$(TAG)

```

For MCP tag as above and run the following commands

```bash
make image TAG=$(TAG)

docker push feedhenry/mcp-standalone:$(TAG)
# builds the 3 different apbs for mcp (android, cordova, iOS)
make apbs TAG=$(TAG)
```


