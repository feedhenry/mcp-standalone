#!/bin/bash

readonly DOCKERHUB_USER="${1}"
readonly DOCKERHUB_PASS="${2}"
readonly DOCKERHUB_ORG="${3}"
readonly LAUNCH_APB_ON_BIND="${4}"
readonly ROUTING_SUFFIX="${5}"
curl -s https://raw.githubusercontent.com/openshift/ansible-service-broker/d93206936751437cc7edaa15a2c88ef317d4c698/templates/deploy-ansible-service-broker.template.yaml > /tmp/deploy-ansible-service-broker.template.yaml

oc login -u system:admin
oc new-project ansible-service-broker
oc process -f /tmp/deploy-ansible-service-broker.template.yaml \
    -n ansible-service-broker \
    -p DOCKERHUB_USER="${DOCKERHUB_USER}" \
    -p DOCKERHUB_PASS="${DOCKERHUB_PASS}" \
    -p DOCKERHUB_ORG="${DOCKERHUB_ORG}" \
    -p BROKER_IMAGE="ansibleplaybookbundle/origin-ansible-service-broker:sprint137" \
    -p ENABLE_BASIC_AUTH="false" \
    -p SANDBOX_ROLE="admin" \
    -p ROUTING_SUFFIX="${ROUTING_SUFFIX}" \
    -p LAUNCH_APB_ON_BIND="${LAUNCH_APB_ON_BIND}" | oc create -f -

if [ "${?}" -ne 0 ]; then
	echo "Error processing template and creating deployment"
	exit
fi