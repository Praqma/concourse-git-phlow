#!/bin/sh

set -e -u -x

echo '{"source":{"url":"https://github.com/praqma/phlow-test.git"}}' | ./check
