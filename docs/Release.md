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

```

For MCP tag as above and run the following commands

```bash
make image TAG=$(TAG)

```

Next update the main template in ```artifacts/openshift/template.json``` change the IMAGE_TAG parameter
to match the the TAG you have just created.
```bash
docker push docker.io/feedhenry/mcp-standalone:$(TAG)
# builds the 3 different apbs for mcp (android, cordova, iOS) copying over the main template
make apbs TAG=$(TAG)
```



