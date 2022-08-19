#!/bin/bash
cat sort | awk '{ printf "{\"name\":\"%s\",\"count\":\"%s\"}\n", $1, $2 }'|jq -s  > sort.json;