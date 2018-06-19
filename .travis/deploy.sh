#!/bin/sh

set -e
set -x

ssh -o BatchMode=yes -o StrictHostKeyChecking=no -o User=deploy pdericson.com mkdir -p bin
scp -o BatchMode=yes -o StrictHostKeyChecking=no -o User=deploy pdericson pdericson.com:bin/pdericson-$TRAVIS_COMMIT

ssh -o BatchMode=yes -o StrictHostKeyChecking=no -o User=deploy pdericson.com bin/pdericson-$TRAVIS_COMMIT -littleboss=reload
