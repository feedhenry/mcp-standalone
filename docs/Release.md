# Upstream Release Process

## Release of independent apbs

- keycloak https://github.com/feedhenry/keycloak-apb
- 3scale https://github.com/feedhenry/3scale-apb
- sync https://github.com/feedhenry/fh-sync-server-apb
- digger (mobile ci cd) https://github.com/feedhenry/aerogear-digger-apb  

Run the following command in each of the above repos:

```bash
make apb_release VERSION=<version> ORIGIN=<remote-origin>
```
This creates a `VERSION` tag in github, and afterwards a build on Dockerhub is triggered to release the image from this `VERSION` tag.

## MCP Release

For MCP tag as above and run the following commands

```bash
make image TAG=$(TAG)
docker push docker.io/feedhenry/mcp-standalone:$(TAG)
```

### MCP included APBs

Next update the main template in ```artifacts/openshift/template.json``` change the IMAGE_TAG parameter
to match the the TAG you have just created.
```bash
# builds the 3 different apbs for mcp (android, cordova, iOS) copying over the main template
make apbs TAG=$(TAG)
```

### Tag the sources in Github

The above adds automated commits to your local branch, afterwards go and look for the latest commit, and create a tag and push it to github:

```bash
git rev-parse HEAD

## use that commit hash when creating the tag 
git tag -a $(TAG) HASH -m "signing tag"
git push origin $(TAG)
```
