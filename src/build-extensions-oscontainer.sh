#!/bin/bash
#Used by cmd/build-extensions-container.go
#Find the RHCOS ociarchive.
ostree_ociarchive=$(find -L ~+ -path '*/builds/latest/*/*-ostree*.ociarchive')
cd src/config || exit
#Start the build replacing the from line.
podman build --from oci-archive:"$ostree_ociarchive" --network=host --build-arg COSA=true -t localhost/extensions-container -f extensions/Dockerfile .
#Call skopeo to generate a extensions container ociarchive
extensions_ociarchive_dir=$(dirname "$ostree_ociarchive")
extensions_ociarchive="${extensions_ociarchive_dir}/extensions-container.ociarchive"
skopeo copy containers-storage:localhost/extensions-container oci-archive:"$extensions_ociarchive"

output=$1; echo "$extensions_ociarchive" > "$output"
