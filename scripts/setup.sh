#!/bin/sh
git config --local core.hooksPath scripts/githooks
chmod a+x scripts/githooks/pre-commit