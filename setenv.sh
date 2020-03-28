#!/bin/bash

if [ "${BASH_SOURCE[0]}" == "${0}" ]; then
  echo "Script must be called with . or source:" >&2
  echo -e "\tsource setenv.sh" >&2
  exit 1
fi

THE_MODULE_ROOT="$(git rev-parse --show-toplevel 2>/dev/null)"
if [ -z "$THE_MODULE_ROOT" ]; then
  export THE_MODULE_ROOT="$(cd $(dirname ${BASH_SOURCE[0]})&&pwd)"
fi
cd "$THE_MODULE_ROOT"

for dir in _target _testing utl bin; do
  if [ -d "$dir" ]; then
    case ":${PATH}:" in
      *${THE_MODULE_ROOT}/${dir}*);;
      *) export PATH="$PATH:${THE_MODULE_ROOT}/${dir}";;
    esac
  fi
done
