#!/bin/sh
#
# This script is a wrapper to execute the build-truss.sh script potentially
# as another user. This is userful for not having files created inside the
# container, but written to a volume mounted from the host, from being owned
# by the default Docker user, root.
set -e

# Look for environment variables BUILD_UID and BUILD_GID to create a user and
# group, the hand off execution to gosu running as this user. If no BUILD_UID
# is provided then just execute the arguments passed to this script as the
# originally-invoking user.
# if [ -z "$BUILD_UID" ]; then
#     if [ $(id -u) -eq 1 ]; then
#         echo "$0: warning: files built will be owned by uid 0 (root) if host is Linux"
#     fi
#     exec "$@"
# fi
# if [ -z "$BUILD_GID" ]; then
#     echo "$0: error: BUILD_UID environment var set, but not BUILD_GID"
#     exit 1
# fi
set -x

WORK_DIR="$1"
BUILD_UID=$(stat -c %u $WORK_DIR)
BUILD_GID=$(stat -c %g $WORK_DIR)
groupadd -fg $BUILD_GID buildgrp
useradd --create-home --uid $BUILD_UID --gid $BUILD_GID buildusr

fix_owner () {
    chown -R $BUILD_UID:$BUILD_GID "$WORK_DIR"
}

trap fix_owner EXIT

shift
cd $WORK_DIR

/usr/local/bin/gosu $BUILD_UID:$BUILD_GID "$@"

