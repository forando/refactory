#!/bin/bash

# Copy this script into terragrunt folder where imports.csv is generated.
# Run it to import all necessary resources.

rollback () {
  echo ""
  echo "Rolling back imports..."
  echo "---"
  for address in "${$1[@]}"
  do
    terragrunt state rm "${address}"
  done
}

addresses=()
while IFS=";" read -r module address id
do
  echo "================"
  echo "Module: $module"
  echo "Address: $address"
  echo "ID: $id"
  echo "---"
  terragrunt import "$(echo "${address}" | sed "s/'/\"/g")" "${id}"
  error=$?
  addresses+=("${address}")
  if [ "${error}" -ne 0 ]; then
    rollback addresses
    break
  fi
  echo "================"
  echo ""
done < imports.csv