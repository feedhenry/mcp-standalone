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

To perform a release of mcp-standalone you'll need:

* A GitHub access token with repo access permissions
* `goreleaser` - `go get github.com/goreleaser/goreleaser`

Set the GitHub access token to the GITHUB_TOKEN environment variable.

```
export GITHUB_TOKEN=mysecrettoken
```

Run the following command with `TAG` set to a new tag.

```bash
make release TAG=0.0.1
```

This will perform a number of steps:
* Build docker images
* Tag the repo with `TAG`
* Push the new tag
* Create a draft release on GitHub with binaries attached
* Push the new image to DockerHub with the tag `TAG`

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
