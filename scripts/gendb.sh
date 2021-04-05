#!/bin/bash

echo "valid-data" > recipes/.db
for i in recipes/*.yml ; do
    app=$(basename $i | sed 's|.yml||g')
    desc=$(cat $i | grep desc | awk '{print $2}')
    if [[ -e tmp/pkg/$app ]] ; then
        echo "$app $desc" >> recipes/.db
    fi
done