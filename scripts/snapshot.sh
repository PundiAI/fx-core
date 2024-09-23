#!/usr/bin/env bash

# crontab -e
# 0 0 * * * NETWORK=testnet CHAIN_NAME=fxcore /usr/local/bin/snapshot.sh >> /home/ubuntu/fxcore-snapshot.log 2>&1

set -euo pipefail

export PATH=$PATH:/usr/local/bin # for aws
export NETWORK=${NETWORK:-"mainnet"}
export CHAIN_NAME=${CHAIN_NAME:-"fxcore"}
export SNAPSHOT_TYPE=${SNAPSHOT_TYPE:-"pruned"}
export AWS_S3_BUCKET=${AWS_S3_BUCKET:-"fx-mainnet"}
export DEBUG=${DEBUG:-"false"}
[[ "$DEBUG" == "true" ]] && set -x

s3_ID=$(aws s3api get-bucket-acl --bucket "$AWS_S3_BUCKET" | jq '.Owner.ID')
docker_container="${CHAIN_NAME}"
snapshot="$CHAIN_NAME-$SNAPSHOT_TYPE-snapshot-$NETWORK-$(date +%F).tar.gz"
s3_dir=/home/ubuntu/s3
snapshot_path="$s3_dir/$snapshot"
s3_uri="s3://$AWS_S3_BUCKET/"

# only support Linux
[[ "$(uname -s)" != "Linux" ]] && echo "Only support Linux" && exit 1

# install jq, curl, unzip
for cmd in "jq" "curl" "unzip"; do
  if ! command -v "$cmd" &>/dev/null; then
    sudo apt-get install "$cmd" -y
  fi
done

# install aws
if ! command -v aws &>/dev/null; then
  curl -s "https://awscli.amazonaws.com/awscli-exe-linux-$(uname -m).zip" -o "awscli-linux.zip"
  unzip -n awscli-linux.zip
  sudo ./aws/install
  rm -rf ./awscli-linux.zip ./aws
fi

# make sure there is enough disk space
disk_avail="$(df -h | grep "${CHAIN_NAME}" | awk '{print $4}' | sed 's/[^0-9]*//g')"
[[ "$disk_avail" -lt 5 ]] && echo "Disk space is less than 5GB" && exit 1

# check if docker container exists
[[ $(docker ps -f name="$docker_container" -q) == "" ]] && echo "Docker container $docker_container not found" && exit 1

# check if node is catching up
[[ $(docker exec "$docker_container" "/root/.${CHAIN_NAME}/cosmovisor/current/bin/${CHAIN_NAME}d" status -o json | jq -r '.SyncInfo.catching_up') != "false" ]] && echo "Node $CHAIN_NAME is catching up" && exit 1

# check if upgrade-info.json exists
[[ ! -f "/home/ubuntu/.$docker_container/data/upgrade-info.json" ]] && echo "upgrade-info.json not found" && exit 1

# recreate s3 directory
[[ -d "$s3_dir" ]] && rm -rf "$s3_dir"
mkdir -p "$s3_dir"
trap 'rm -rf "$s3_dir"' EXIT

echo "start create s3 snapshot for $CHAIN_NAME"

echo "stop $docker_container"
docker stop "$docker_container"

echo "wait $docker_container stop"
while true; do
  if [[ $(docker inspect "$docker_container" -f '{{.State.Status}}') != "running" ]]; then
    break
  fi
  sleep 1
done

sudo du -h -d 1 "/home/ubuntu/.$docker_container/data" # for debug

echo "prune data"
docker_image=$(docker inspect "$docker_container" -f '{{.Config.Image}}')
docker run --rm --entrypoint "${CHAIN_NAME}d" -v "/home/ubuntu/.$docker_container:/root/.$CHAIN_NAME" "$docker_image" data prune-compact all --enable_pruning --height 3600
# Repeated execution can further compress the size
docker run --rm --entrypoint "${CHAIN_NAME}d" -v "/home/ubuntu/.$docker_container:/root/.$CHAIN_NAME" "$docker_image" data prune-compact all --enable_pruning --height 3600

echo "compress data"
sudo chown -R ubuntu:ubuntu "/home/ubuntu/.$docker_container/data"
rm -rf "/home/ubuntu/.$docker_container/data/tx_index.db"
rm -rf "/home/ubuntu/.$docker_container/data/snapshots"
du -h -d 1 "/home/ubuntu/.$docker_container/data"
tar -zcf "$snapshot_path" -C "/home/ubuntu/.$docker_container" "data"
du -h "$snapshot_path"

echo "calculate md5"
md5sum "$snapshot_path" >"$snapshot_path".md5

echo "put s3"
aws s3 cp "$snapshot_path" "$s3_uri"
aws s3 cp "$snapshot_path".md5 "$s3_uri"

echo "put s3 acl"
aws s3api put-object-acl --bucket "$AWS_S3_BUCKET" --key "$snapshot" --grant-full-control id="$s3_ID" --grant-read uri=http://acs.amazonaws.com/groups/global/AllUsers
aws s3api put-object-acl --bucket "$AWS_S3_BUCKET" --key "$snapshot".md5 --grant-full-control id="$s3_ID" --grant-read uri=http://acs.amazonaws.com/groups/global/AllUsers

echo "clear s3 old file"
num=$(aws s3 ls "$s3_uri" | grep "$CHAIN_NAME-$SNAPSHOT_TYPE-snapshot-$NETWORK" | grep -oP '[0-9]+-[0-9]+-[0-9]+' | sort -u | wc -l)
if [ "$num" -gt 3 ]; then
  dates=$(aws s3 ls "$s3_uri" | grep "$CHAIN_NAME-$SNAPSHOT_TYPE-snapshot-$NETWORK" | grep -oP '[0-9]+-[0-9]+-[0-9]+' | sort -u | head -n-3)
  for i in $dates; do
    aws s3 rm "${s3_uri}$CHAIN_NAME-$SNAPSHOT_TYPE-snapshot-$NETWORK-${i}.tar.gz"
    aws s3 rm "${s3_uri}$CHAIN_NAME-$SNAPSHOT_TYPE-snapshot-$NETWORK-${i}.tar.gz.md5"
  done
fi

echo "reset data"
mv "/home/ubuntu/.$docker_container/data/upgrade-info.json" "/home/ubuntu/.$docker_container/upgrade-info.json"
rm -rf "/home/ubuntu/.$docker_container/data" && mkdir -p "/home/ubuntu/.$docker_container/data"
echo -e '{\n  "height": "0",\n  "round": 0,\n  "step": 0\n}' >"/home/ubuntu/.$docker_container/data/priv_validator_state.json"
mv "/home/ubuntu/.$docker_container/upgrade-info.json" "/home/ubuntu/.$docker_container/data/upgrade-info.json"

statesync_node_rpc=$(grep 'rpc_servers' "/home/ubuntu/.$docker_container/config/config.toml" | grep -oE '([0-9]{1,3}\.){3}[0-9]{1,3}:[0-9]{1,5}' | head -n 1)
latest_block_height=$(curl -s "$statesync_node_rpc/block" | jq -r '.result.block.header.height')
block_height=$(("$latest_block_height" - 1000))
block_hash=$(curl -s "$statesync_node_rpc/block?height=$block_height" | jq -r '.result.block_id.hash')

sed -i "s/trust_height = [0-9]\{1,\}/trust_height = $block_height/g" "/home/ubuntu/.$docker_container/config/config.toml"
sed -i "s/trust_hash = \"[a-f0-9A-F]\{64\}\"/trust_hash = \"$block_hash\"/g" "/home/ubuntu/.$docker_container/config/config.toml"

echo "restart $docker_container"
docker start "$docker_container"

echo "All complete!"
