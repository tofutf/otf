#!/usr/bin/env bash

set -e

PORT_FORWARD_TIMEOUT_SECONDS=10

function start_caddy(){
  coproc caddy reverse-proxy --from https://tofutf.pvotal.local:8081 --to http://tofutf.pvotal.local:8080 --internal-certs --disable-redirects 
  caddy_pid=$COPROC_PID
}

# Start a kubectl port-forward and wait until the port is active
#
# The local port is automatically selected and is returned as global variable
#
#    port_forward_local_port
#
function start_port_forward() {
    coproc kubectl port-forward service/postgres-postgresql :5432 </dev/null 2>&1
    port_forward_pid=$COPROC_PID
    while IFS='' read -r -t $PORT_FORWARD_TIMEOUT_SECONDS -u "${COPROC[0]}" LINE
    do
        if [[ "$LINE" == *"Forwarding from"* ]]; then
            port_forward_local_port="${LINE#Forwarding from 127.0.0.1:}"
            port_forward_local_port="${port_forward_local_port%% -> *}"
            if [ -z "${port_forward_local_port}" ]; then
              echo "ERROR: Failed to get local address for port-forward"
              echo "kubectl output line: $LINE"
              exit 1
            fi
            # Remaining output is on stderr, which we don't capture, so we
            # should be fine to ignore the stdout file descriptor now and
            # port_forward_pid remains set and will be used on cleanup
            #
            return
        else
            echo "kubectl: ${LINE}"
        fi
    done
    # if we reached here, read failed, likely due to the coproc exiting
    if [ -n "${port_forward_pid:-}" ]; then
      port_forward_ecode=
      wait $port_forward_pid || port_forward_ecode=$?
      echo "port-forward request failed? Exit code $port_forward_ecode"
    else
      echo "port forward request failed? Could not get kubectl port-forward's pid"
    fi
    exit 1
}

# Assumes there's only one coproc
function kill_port_forward() {
  if [ -n "${port_forward_pid:-}" ]; then
    kill ${port_forward_pid} || true
    wait -f ${port_forward_pid} || true
  fi
  port_foward_pid=
}

# Assumes there's only one coproc
function kill_caddy() {
  if [ -n "${caddy_pid:-}" ]; then
    kill ${caddy_pid} || true
    wait -f ${caddy_pid} || true
  fi
  caddy_pid=
}


trap kill_caddy EXIT

start_caddy

# here we grab all of the configuration from the running tofutf instance.
#$(kubectl get deployments -n tofutf tofutf -oyaml | yq '.spec.template.spec.containers[0].env | select(.value != null) | select(.name == "OTF*") | map({"key": (.name), "value": (.value)}) | from_entries' | awk '{ print "export " substr($1, 1, length($1)-1) "=" $2}')

export OTF_SKIP_TLS_VERIFICATION="true"
export OTF_HOSTNAME=tofutf.pvotal.local
export OTF_LOG_FORMAT=default
export OTF_LOG_HTTP_REQUESTS="false"
export OTF_DATABASE=postgres://tofutf:AtxcIOp3Ons83SKd@localhost:5432/tofutf?sslmode=disable
export OTF_SECRET=2876cb147697052eec5b3cdb56211681
export OTF_PROVIDER_PROXY_URL=https://registry.terraform.io/v1/providers/
export OTF_GITHUB_CLIENT_ID
export OTF_GITHUB_CLIENT_SECRET
export OTF_GITHUB_HOSTNAME=github.com
export OTF_GITLAB_CLIENT_ID
export OTF_GITLAB_CLIENT_SECRET
export OTF_GITLAB_HOSTNAME=gitlab.com
export OTF_OIDC_NAME=internal-client
export OTF_OIDC_ISSUER_URL=https://dex.pvotal.local
export OTF_OIDC_CLIENT_ID=internal-client
export OTF_SITE_TOKEN=site-token
export OTF_SANDBOX="false"


export OTF_SANDBOX=false
export OTF_V="100"
export OTF_OIDC_CLIENT_SECRET=5GCVLQ7qAQcqfnmRlaFqWHas4RijCKH7
export OTF_PROVIDER_PROXY_URL=https://registry.opentofu.org/v1/providers/
export OTF_PROVIDER_PROXY_IS_ARTIFACTORY=false
export OTF_SKIP_TLS_VERIFICATION=true

air