#!/bin/sh

set -x

# Stop OpenShift
oc cluster down

# Remove openshift config dir things
SCRIPT_PATH=$(dirname $0)
SCRIPT_ABSOLUTE_PATH=$(cd $SCRIPT_PATH && pwd)

sudo rm -rf ${SCRIPT_ABSOLUTE_PATH}/master
sudo rm -rf ${SCRIPT_ABSOLUTE_PATH}/node-localhost
sudo rm -rf ${SCRIPT_ABSOLUTE_PATH}/openshift-data
sudo rm -rf ${SCRIPT_ABSOLUTE_PATH}/openshift-pvs

if [ "$(uname -s)" == "Linux"  ]
then
  findmnt -lo TARGET | grep ${SCRIPT_ABSOLUTE_PATH}/openshift-volumes | xargs -r sudo umount
fi

sudo rm -rf ${SCRIPT_ABSOLUTE_PATH}/openshift-volumes
