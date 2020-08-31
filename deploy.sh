#!/bin/bash
set -ev

tar -zcvf vss-mac.tar.gz build/darwin/vss
tar -zcvf vss-linux.tar.gz build/linux/vss
tar -zcvf vss-windows.tar.gz build/windows/vss
