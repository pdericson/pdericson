#!/bin/sh

set -e
set -x

ssh -o BatchMode=yes -o StrictHostKeyChecking=no -o User=deploy pdericson.com mkdir -p bin
scp -o BatchMode=yes -o StrictHostKeyChecking=no -o User=deploy pdericson pdericson.com:bin/

ssh -o BatchMode=yes -o StrictHostKeyChecking=no -o User=deploy pdericson.com sh -c '"bin/pdericson -littleboss=reload || nohup bin/pdericson -littleboss=start > pdericson.log 2>&1 &"'
