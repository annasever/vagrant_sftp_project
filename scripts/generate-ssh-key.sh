#!/bin/bash

# Generate SSH keys for sftp user if not exist
if [ ! -f /home/sftp/.ssh/id_rsa ]; then
    ssh-keygen -t rsa -b 4096 -f /home/sftp/.ssh/id_rsa -N ""
fi
