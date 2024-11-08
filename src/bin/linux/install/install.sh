#! /bin/bash

root=$(cd "$(dirname "$0")"; pwd)

if [[ ! -d $root/config ]]; then
    mkdir $root/config
fi

touch $root/config/env.sh

if [[ -z $(grep "VM_HOME" $root/config/env.sh) ]] ; then
    echo "export PATH="'$PATH:$VM_HOME/bin' >> $root/config/env.sh
fi

touch $HOME/.bashrc

if [[ -z $(grep "VM_HOME" $HOME/.bashrc) ]] ; then
    echo "export VM_HOME=$root" >> $HOME/.bashrc
    echo "[[ -s "'$VM_HOME'"/config/env.sh ]] && source "'$VM_HOME'"/config/env.sh" >> $HOME/.bashrc
fi

source $HOME/.bashrc
