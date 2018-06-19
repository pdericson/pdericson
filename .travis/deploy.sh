#!/bin/sh

set -e
set -x

ssh -o BatchMode=yes -o StrictHostKeyChecking=no -o User=deploy pdericson.com mkdir -p bin
scp -o BatchMode=yes -o StrictHostKeyChecking=no -o User=deploy pdericson pdericson.com:bin/

ssh -o BatchMode=yes -o StrictHostKeyChecking=no -o User=deploy pdericson.com sh -c "pdericson -littleboss=reload || nohup pdericson -littleboss=start &"
