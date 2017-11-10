## Required docker images
In order to reduce bandwidth consumption on the day, please pre-pull the required docker images before attending the Face 2 Face, using the following command:
```
for i in quay.io/3scale/apicast:master docker.io/feedhenry/cordova-app-apb:0.0.6 docker.io/feedhenry/ios-app-apb:0.0.6 docker.io/feedhenry/android-app-apb:0.0.6 docker.io/feedhenry/mcp-standalone:0.0.6 docker.io/feedhenry/keycloak-apb:0.0.6 docker.io/feedhenry/3scale-apb:0.0.6 docker.io/feedhenry/aerogear-digger-apb:0.0.6 docker.io/ansibleplaybookbundle/origin-ansible-service-broker:sprint139.1 docker.io/openshift/origin-haproxy-router:v3.7.0-rc.0 docker.io/openshift/origin-deployer:v3.7.0-rc.0 docker.io/openshift/origin:v3.7.0-rc.0 docker.io/openshift/origin-docker-registry:v3.7.0-rc.0 docker.io/openshift/origin-pod:v3.7.0-rc.0 docker.io/openshift/origin-service-catalog:v3.7.0-rc.0 quay.io/coreos/etcd:latest docker.io/jimmidyson/keycloak-openshift:2.5.4.Final docker.io/rhmap/redis:2.18.22 docker.io/aerogear/digger-android-slave-image:AGDIGGER-177 docker.io/aerogear/digger-android-sdk-image:FH-v3.19 docker.io/feedhenry/fh-sync-server-apb:0.0.6 docker.io/feedhenry/fh-sync-server:0.0.6 docker.io/centos/mongodb-32-centos7 docker.io/openshift/jenkins-2-centos7; do docker pull $i; done
```

## Local setup

Follow the local setup guide from [here](https://github.com/feedhenry/mcp-standalone/blob/master/docs/walkthroughs/local-setup.adoc#requirements), taking care to setup prerequisites and any firewalld rules (if on Linux)

When prompted with `DockerHub Tag (defaults to latest)` in the installer, use `0.0.6`.

## 3 Scale
If you intend to experiment with 3 Scale, you will need to set up a trial account with them, this can take a day or two to be provisioned, start the process [here](https://www.3scale.net/signup/)

## Docker on Mac
Edit the settings of docker on mac and allow it 6Gb of RAM as follows: r-click systray icon > preferences... > Advanced

