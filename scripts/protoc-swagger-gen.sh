#!/usr/bin/env bash

set -eo pipefail

proto_dirs=$(find ./proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)

for dir in $proto_dirs; do

  # generate swagger files (filter query files)
  query_file=$(find "${dir}" -maxdepth 2 -name 'query.proto')
  if [[ ! -z "$query_file" ]]; then
    buf protoc  \
      -I "proto" \
      -I "third_party/proto" \
      "$query_file" \
      --swagger_out=./client/docs \
      --swagger_opt=logtostderr=true --swagger_opt=fqn_for_swagger_name=true --swagger_opt=simple_operation_ids=true
  fi
done

rm -f ./client/docs/swagger.yaml

# download Cosmos SDK swagger doc
SDK_VERSION=$(go list -m -f '{{ .Version }}' github.com/cosmos/cosmos-sdk)
echo "SDK version ${SDK_VERSION}"
wget "https://raw.githubusercontent.com/cosmos/cosmos-sdk/${SDK_VERSION}/client/docs/swagger-ui/swagger.yaml" -P ./client/docs
mv ./client/docs/swagger.yaml ./client/docs/swagger-sdk.yaml

# # download IBC swagger doc
IBC_VERSION=$(go list -m -f '{{ .Version }}' github.com/cosmos/ibc-go/v2)
echo "IBC version ${IBC_VERSION}"
wget "https://raw.githubusercontent.com/cosmos/ibc-go/${IBC_VERSION}/docs/client/swagger-ui/swagger.yaml" -P ./client/docs
mv ./client/docs/swagger.yaml ./client/docs/swagger-ibc.yaml

# combine swagger files
# uses nodejs package `swagger-combine`.
# all the individual swagger files need to be configured in `config.json` for merging
swagger-combine ./client/docs/config.json -o ./client/docs/swagger.yaml -f yaml --continueOnConflictingPaths true --includeDefinitions true

if [ ! -d ./client/docs/swagger-ui ]; then
  wget https://github.com/swagger-api/swagger-ui/archive/refs/tags/v4.11.0.zip -O ./client/docs/swagger-ui-4.11.0.zip
  unzip ./client/docs/swagger-ui-4.11.0.zip -d ./client/docs
  mv ./client/docs/swagger-ui-4.11.0/dist ./client/docs/swagger-ui
  rm ./client/docs/swagger-ui-4.11.0.zip
  rm -rf ./client/docs/swagger-ui-4.11.0
fi

# move generated swagger file to swagger-ui directory
cp ./client/docs/swagger.yaml ./client/docs/swagger-ui/

# move generated swagger file to swagger-ui directory
statik -src=./client/docs/swagger-ui -dest=client/docs -f -m

if [ -n "$(git status --porcelain)" ]; then
  echo "\033[91mSwagger docs updated\033[0m"
else
  echo "\033[92mSwagger docs already in sync\033[0m"
fi
