## Required docker images
In order to reduce bandwidth consumption on the day, please pre-pull the required docker images before attending the Face 2 Face, using the following command:
```
for i in quay.io/3scale/apicast:master feedhenry/cordova-app-apb:0.0.6 feedhenry/ios-app-apb:0.0.6 feedhenry/android-app-apb:0.0.6 feedhenry/mcp-standalone:0.0.6 feedhenry/keycloak-apb:0.0.6 feedhenry/3scale-apb:0.0.6 feedhenry/aerogear-digger-apb:0.0.6 ansibleplaybookbundle/origin-ansible-service-broker:sprint139.1 openshift/origin-haproxy-router:v3.7.0-rc.0 openshift/origin-deployer:v3.7.0-rc.0 openshift/origin:v3.7.0-rc.0 openshift/origin-docker-registry:v3.7.0-rc.0 openshift/origin-pod:v3.7.0-rc.0 openshift/origin-service-catalog:v3.7.0-rc.0 quay.io/coreos/etcd:latest jimmidyson/keycloak-openshift:2.5.4.Final rhmap/redis:2.18.22 aerogear/digger-android-slave-image:AGDIGGER-177 aerogear/digger-android-sdk-image:FH-v3.19 feedhenry/fh-sync-server-apb:0.0.6 feedhenry/fh-sync-server:0.0.6 centos/mongodb-32-centos7 openshift/jenkins-2-centos7; do docker pull $i; done
```

## Local setup
Export your Docker credentials:
```
export DOCKER_USER=<docker user>
export DOCKER_PASS=<docker pass>
```

For the face to face, we are working from the 0.0.6 git tag.
```
cd /path/to/this/repo
git checkout 0.0.6
ansible-galaxy install -r ./installer/requirements.yml
ansible-playbook installer/playbook.yml -e "dockerhub_username=$DOCKER_USER" -e "dockerhub_password=$DOCKER_PASS" -e "dockerhub_tag=0.0.6" --ask-become-pass
```
For more detailed instructions, look [here](https://github.com/feedhenry/mcp-standalone/blob/master/docs/walkthroughs/local-setup.adoc#local-setup).

## 3 Scale
If you intend to experiment with 3 Scale, you will need to set up a trial account with them, this can take a day or two to be provisioned, start the process [here](https://www.3scale.net/signup/)

## Docker on Mac
Edit the settings of docker on mac and allow it 6Gb of RAM as follows: r-click systray icon > preferences... > Advanced

