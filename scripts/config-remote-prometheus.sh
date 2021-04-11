#!/bin/bash

if [[ -z "${REMOTE_PROM_UNAME}" ]]; then
  echo "Skipping configuration of remote prometheus due to lack of REMOTE_PROM_UNAME environment variable."
  exit 1
fi

if [[ -z "${REMOTE_PROM_PWORD}" ]]; then
  echo "Skipping configuration of remote prometheus due to lack of REMOTE_PROM_PWORD environment variable."
  exit 1
fi

configFile="configs/prometheus/prometheus.yml"
remoteUname=${REMOTE_PROM_UNAME}
remotePword=${REMOTE_PROM_PWORD}

if grep -F -q "remote_write" "$configFile"; then
  echo "Skip adding 'remote_write' configuration."
else
  echo "Adding 'remote_write' configuration..."
  cp configs/prometheus/prometheus.yml configs/prometheus/prometheus.yml.bak
  echo "" >> configs/prometheus/prometheus.yml
  echo "remote_write:" >> configs/prometheus/prometheus.yml
  echo "- url: https://prometheus-us-central1.grafana.net/api/prom/push" >> configs/prometheus/prometheus.yml
  echo "  basic_auth:" >> configs/prometheus/prometheus.yml
  echo "    username: ${remoteUname}" >> configs/prometheus/prometheus.yml
  echo "    password: ${remotePword}" >> configs/prometheus/prometheus.yml
  echo "Done."
fi