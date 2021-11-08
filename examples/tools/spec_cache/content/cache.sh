#!/bin/bash

rm -f Gemfile.lock
echo -e "\xF0\x9F\x9A\x80 Caching gems..."
bundle config --global set cache_path /cache
bundle cache --no-install
echo -e "\xE2\x9C\x85 Done"
rm -rf .bundle
rm -f Gemfile.lock
