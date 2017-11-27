#!/bin/sh

set -x

# Stop OpenShift
oc cluster down

# Remove openshift config dir things
SCRIPT_PATH=$(dirname $0)
OPENSHIFT_DATA_DIR=$(cd $SCRIPT_PATH/../ui && pwd)

sudo rm -rf ${OPENSHIFT_DATA_DIR}/master
sudo rm -rf ${OPENSHIFT_DATA_DIR}/node-localhost
sudo rm -rf ${OPENSHIFT_DATA_DIR}/openshift-data
sudo rm -rf ${OPENSHIFT_DATA_DIR}/openshift-pvs

if [ "$(uname -s)" == "Linux"  ]
then
  findmnt -lo TARGET | grep ${OPENSHIFT_DATA_DIR}/openshift-volumes | xargs -r sudo umount
fi

sudo rm -rf ${OPENSHIFT_DATA_DIR}/openshift-volumes
