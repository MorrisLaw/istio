#!/bin/bash

# NON-PRODUCTION script that runs on your dev VM for the purpose of working on Ztunnel.
# This script assumes you already have the necessary tools installed to run the referenced 
# scripts.
#
# How to use the script?
# - save the script somewhere and run chmod +x on the file to make it executable.
# - you can now use the path to run it (or, you can set an alias to the script 
#   path in your bash_profile like I did).
# - make a change to a local copy of the ztunnel repo.
# - navigate back to your istio repo (this whole thing assumes that istio is in 
#   the same root directory as ztunnel).
# - when you call the script (wherever you stored it doesn't matter) make sure you're 
#   within the istio root directory.
# - you now have a locally built version of the ztunnel pod running in your cluster!

# check if we're in the istio directory
directory="$(echo $PWD | awk -F "/" '{print $NF}')"
if [[ "${directory}" != "istio" ]]
then
	echo "You need to run this from within the istio root directory... not $(echo $PWD)" >&2
	exit 1
fi

# ensure that BUILD_ZTUNNEL_REPO is unset before running script 
# to build kind cluster or else the kind script will not work
export BUILD_ZTUNNEL_REPO=""

# run the command that actually brings up the kind cluster
./prow/integ-suite-kind.sh --skip-cleanup

# install ambient mesh components such as ztunnel
istioctl install -y --set profile=ambient 

# set the appropriate env vars to make building 
# ztunnel from local repo possible
export BUILD_ZTUNNEL_REPO="$PWD/../ztunnel"
export HUB=localhost:5000
export TAG=istio-testing

# build ztunnel from local ztunnel repo and replace the currently
# running ztunnel pod
BUILD_WITH_CONTAINER=0 BUILD_ALL=false DOCKER_TARGETS=docker.ztunnel make dockerx.pushx && kubectl -n istio-system set image daemonset/ztunnel istio-proxy=localhost:5000/ztunnel:${TAG}
