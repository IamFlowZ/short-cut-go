#!/bin/bash

BIN_INSTALL_LOC=/usr/bin
FOLDER_LOC=/var/shortcut

cp build/short-cut $BIN_INSTALL_LOC/short-cut

if [ ! -d $FOLDER_LOC ]; then
    mkdir $FOLDER_LOC
fi

touch $FOLDER_LOC/shortcuts

sudo chmod 666 $FOLDER_LOC/shortcuts
