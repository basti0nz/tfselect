#!/bin/bash
rm  ../dist/tfselect-alpine.tar.gz || true
env ../version
cd .. && docker build  -f alpine/Dockerfile.builder  -t alpine-build .
cd dist && docker run --rm -iv${PWD}:/host-volume alpine-build  sh -s <<EOF
chown -v $(id -u):$(id -g) /app/tfselect
cp -va /app/tfselect  /host-volume/tfselect-alpine
EOF
cd ../dist/ && tar -czvf tfselect-alpine.tar.gz tfselect-alpine


