#!/bin/bash
set -e

readonly service="$1"
readonly output_dir="$2"
readonly package="$3"

# Function to convert kebab-case to snake_case
kebab_to_snake() {
    input_string="$1"
    echo "${input_string//-/_}"
}

IFS='/' read -r -a folders <<< "$output_dir"
for i in "${!folders[@]}"
do
    folder_path="${folders[*]:0:i+1}"
    folder_path=$(echo "$folder_path" | tr ' ' '/')
    mkdir -p "$folder_path"
done

for i in "${!folders[@]}"
do
    folder_path="${folders[*]:0:i+1}"
    folder_path=$(echo "$folder_path" | tr ' ' '/')
    mkdir -p "$folder_path"
done

oapi-codegen -generate skip-prune,types -o "$output_dir/openapi_shared_types.gen.go" -package "$package" "api/openapi/shared_components.yml"
oapi-codegen -generate types -o "$output_dir/openapi_types.gen.go" -import-mapping "./shared_components.yml:github.com/layarda-durianpay/go-skeleton/$output_dir" -package "$package" "api/openapi/$service.yml"
oapi-codegen -generate gorilla-server -o "$output_dir/openapi_api.gen.go" -import-mapping "./shared_components.yml:github.com/layarda-durianpay/go-skeleton/$output_dir" -package "$package" "api/openapi/$service.yml"

echo "done generating $output_dir"