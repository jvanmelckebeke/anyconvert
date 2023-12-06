#!/bin/bash

path_to_image="/home/jarivm/Downloads/sample.webp"
retries=1

ARGS=$(getopt --options "r:" --longoptions "retries:" -- "${@}")

# parse args

eval "set -- ${ARGS}"

while true; do
  case "${1}" in
    (-r | --retries)
      retries=$2
      shift 2
    ;;
    (--)
      shift
      break
  esac
done

echo "retries=$retries"



# Step 1: Send a POST request to upload the image

upload_response=$(curl -X POST -F "file=@$path_to_image" http://localhost:8888/api/image)

echo $upload_response

# Extract the ID from the response
id=$(echo $upload_response | jq -r '.status' | cut -d'/' -f3)

for ((i=1;i<=retries;i++));
do
  sleep 1

  # Step 2: Send a GET request to check the status
  status_response=$(curl http://localhost:8888/api/status/$id)

  # Print the final status response
  echo "Status Response: "
  echo "$status_response" | jq .

  status=$(echo "$status_response" | jq  -r .status)


  if [ "$status" = "done" ]; then
    echo "request completed successfully"
    exit 0

  elif [ "$status" = "error" ]; then
    echo "request contained an error"
    exit 1
  fi
done


echo "request could not be completed"
