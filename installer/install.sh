#!/bin/bash

readonly SCRIPT_PATH=$(dirname $0)
readonly SCRIPT_ABSOLUTE_PATH=$(cd $SCRIPT_PATH && pwd)

readonly RED=$(tput setaf 1)

readonly VER_EQ=0
readonly VER_GT=1
readonly VER_LT=2

function banner() {
echo '  __  __  ____ ____'
echo ' |  \/  |/ ___|  _ \'
echo ' | |\/| | |   | |_) |'
echo ' | |  | | |___|  __/'
echo ' |_|  |_|\____|_|'
echo ''
}

oc_install_dir="/usr/local/bin"
oc_version_comparison=${VER_LT}

# Returns:
# 0 - =
# 1 - >
# 2 - <
function compare_version () {
  if [[ $1 == $2 ]]; then
    return 0
  fi
  local IFS=.
  local i ver1=(${1}) ver2=(${2})
  for ((i=${#ver1[@]}; i<${#ver2[@]}; i++)); do
    ver1[i]=0
  done
  for ((i=0; i<${#ver1[@]}; i++)); do
    if [[ -z ${ver2[i]} ]]; then
      ver2[i]=0
    fi
    if ((10#${ver1[i]} > 10#${ver2[i]})); then
      return 1
    fi
    if ((10#${ver1[i]} < 10#${ver2[i]})); then
      return 2
    fi
  done
  return 0
}

function does_not_exist_msg() {
  echo -e "${RED}${1} does not exist on host machine."
  echo -e "It can be installed using ${2}."
}

function check_exists_msg() {
  echo -e "\nChecking ${1} exists"
}

function check_version_msg() {
  echo "Checking ${1} version. Should be ${2}"
}

function check_passed_msg() {
  echo "✓ ${1} check passed."
}

function check_docker() {
  check_version_msg "Docker" "using Stable channel"
  docker_version=$(docker version --format '{{json .Client.Version}}')
  if [[ ${docker_version} == *"-rc"* ]]; then
    echo "${RED}Docker versions from the Edge channel are currently not supported. Switch to a release from the Stable channel"
    exit 1
  fi
  check_passed_msg "Docker"
}

function check_npm() {
  check_exists_msg "NPM"
  command -v npm &>/dev/null
  npm_exists=${?}; if [[ ${npm_exists} -ne 0 ]]; then
    does_not_exist_msg "NPM" "https://nodejs.org/en/download/"
    exit 1
  fi
  check_passed_msg "NPM"
}

function check_python() {
  check_exists_msg "Python"

  command -v python &>/dev/null
  python_exists=${?}; if [[ ${python_exists} -ne 0 ]]; then
    does_not_exist_msg "Python" "pip install ansible>=2.3"
  fi
  check_passed_msg "Python"

  readonly python_version=$(python -c 'import sys; print(".".join(map(str, sys.version_info[:3])))')

  check_version_msg "Python" ">= 2.7"
  compare_version ${python_version} 2.7
  python_version_comparison=${?}; if [[ ${python_version_comparison} -eq ${VER_LT} ]]; then
    echo -e "${RED}Python is < 2.7. Update to >= 2.7."
    exit 1
  fi
  check_passed_msg "Python"
}

function check_ansible() {
  check_exists_msg "Ansible"

  command -v ansible &>/dev/null
  ansible_exists=${?}; if [[ ${ansible_exists} -ne 0 ]]; then
    does_not_exist_msg "Ansible" "pip install ansible>=2.3"
    exit 1
  fi
  check_passed_msg "Ansible"

  readonly ansible_version=$(ansible --version | sed -n '1p' | cut -d " " -f2)

  check_version_msg "Ansible" ">= 2.3"
  compare_version ${ansible_version} 2.3
  ansible_version_comparison=${?}; if [[ ${ansible_version_comparison} -eq ${VER_LT} ]]; then
    echo -e "${RED}Ansible version is < 2.3. Install ansible>=2.3 using pip install ansible>=2.3"
    exit 1
  fi
  check_passed_msg "Ansible"
}

function check_oc() {
  # OPENSHIFT CLIENT TOOLS
  check_exists_msg "OpenShift client tools"

  command -v oc &>/dev/null
  oc_exists=${?}; if [[ ${oc_exists} -ne 0 ]]; then
    echo "? OpenShift Client tools do not exist on host. They will be installed by the MCP installer."
  else
    check_passed_msg "OpenShift Client Tools"
    check_version_msg "OpenShift client tools" ">= 3.7"

    readonly oc_version=$(oc version | sed -n "1p" | cut -d " " -f2 | cut -d "-" -f1 | cut -d "v" -f2)
    compare_version ${oc_version} 3.7

    oc_version_comparison=${?}; if [[ ${oc_version_comparison} -eq ${VER_LT} ]]; then
      echo -e "\n? OpenShift Client tools are less than 3.7"
      read -p "Allow the installer to delete and reinstall the OpenShift client tools? (y/n): " uninstall_client_tools
      if [[ ${uninstall_client_tools} == "y" ]]; then
        echo "Removing oc tool"
        oc_install_dir=$(dirname $(command -v oc))
        rm $(command -v oc)
      else
        echo -e "${RED}The Mobile Control Panel requires oc >= 3.7"
        exit 1
      fi
    fi
    check_passed_msg "OpenShift Client Tools"
  fi
}

function read_oc_install_dir() {
  read -p "Where do you want to install oc? (Defaults to ${oc_install_dir}): " user_oc_install_dir
  oc_install_dir=${user_oc_install_dir:-${oc_install_dir}}
  echo "Updating PATH to include specified directory"
  export PATH="${oc_install_dir}:${PATH}"
}

function run_installer() {
  echo -e "\nThe Mobile Control Panel installer requires valid DockerHub credentials
  to communicate with the DockerHub API. If you enter invalid credentials or then
  Mobile Services will not be available in the Service Catalog.\n"

  read -p "DockerHub Username: " dockerhub_username
  stty -echo
  read -p "DockerHub Password: " dockerhub_password
  stty echo

  echo -e "\nChecking DockerHub credentials are valid...\n"

  curl --fail -u ${dockerhub_username}:${dockerhub_password} https://cloud.docker.com/api/app/v1/service/ &> /dev/null

  if [[ ${?} -ne 0 ]]; then
    echo -e "${RED}Invalid Docker credentials. Run the script again."
    exit 1
  fi

  echo -e "Credentials are valid. Continuing...\n"

  read -p "DockerHub Tag (Defaults to latest): " dockerhub_tag
  dockerhub_tag=${dockerhub_tag:-"latest"}

  echo "Performing and clean and running the installer. You will be asked for your password."

  cd ${SCRIPT_ABSOLUTE_PATH}
  cd .. && make clean &>/dev/null

  set -e
  echo "Installing roles to ${SCRIPT_ABSOLUTE_PATH}/roles"
  ansible-galaxy install -r ./installer/requirements.yml --roles-path="${SCRIPT_ABSOLUTE_PATH}/roles"
  set +e

  if [[ ${oc_version_comparison} -ne ${VER_LT} ]]; then
    echo "Skipping OpenShift client tools installation..."
    ansible-playbook installer/playbook.yml -e "dockerhub_username=${dockerhub_username}" -e "dockerhub_password=${dockerhub_password}" -e "dockerhub_tag=${dockerhub_tag}" --skip-tags "install-oc" --ask-become-pass
  else
    ansible-playbook installer/playbook.yml -e "dockerhub_username=${dockerhub_username}" -e "dockerhub_password=${dockerhub_password}" -e "dockerhub_tag=${dockerhub_tag}" -e "oc_install_parent_dir=${oc_install_dir}" --ask-become-pass
  fi
}

banner
check_docker
check_npm
check_python
check_ansible
check_oc
if [[ ${oc_version_comparison} -eq ${VER_LT} ]]; then
  read_oc_install_dir
fi
run_installer
