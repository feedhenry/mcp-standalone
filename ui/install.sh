#!/bin/sh

set -x
set -e

SCRIPT_PATH=$(dirname $0)
SCRIPT_ABSOLUTE_PATH=$(cd $SCRIPT_PATH && pwd)

# master-config.yaml location
OPENSHIFT_CONFIG_DIR=$SCRIPT_ABSOLUTE_PATH
OPENSHIFT_MASTER_CONFIG=$OPENSHIFT_CONFIG_DIR/master/master-config.yaml

# Enable Extension Development
sudo chmod 666 $OPENSHIFT_MASTER_CONFIG
cd $SCRIPT_ABSOLUTE_PATH
npm i && ./node_modules/.bin/bower install
node update_master_config.js $OPENSHIFT_MASTER_CONFIG
sudo chmod 644 $OPENSHIFT_MASTER_CONFIG

