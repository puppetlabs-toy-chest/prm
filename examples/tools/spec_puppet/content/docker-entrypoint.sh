#!/bin/bash

echo -e "\xF0\x9F\x9A\x80 Installing local gems..."
# install gems locally from the cache
bundle config  --global set cache_path /cache
bundle install --local
echo -e "\xE2\x9C\x85 Done"

echo -e "\xE2\x9C\xA8 Clean up "
# handle the possibility of multiple newer puppet gems being installed...
gem list puppet | grep puppet | head -n 1 | grep -oP '[0-9]+\.[0-9]+\.[0-9]+,' | cut -d , -f 1 | while read -r line ; do
    gem uninstall puppet --version "$line"
done
echo -e "\xE2\x9C\x85 Done"

echo -e "\xF0\x9F\x8E\xAF Running command: $CMD_ENTRY $*"
exec $CMD_ENTRY $*
