#!/bin/bash
# personal test script to test the api

path_to_videos=(
  "/home/jarivm/Downloads/example.gif"
  "/home/jarivm/Downloads/animated.webp"
  "/home/jarivm/Downloads/Wow-gif.webp"
  "/home/jarivm/Downloads/large.webm"
)

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
for path_to_video in "${path_to_videos[@]}" ; do
  tput bold
  echo "$path_to_video"
  echo ""
  tput sgr0


  upload_response=$(curl -s -X POST -F "file=@$path_to_video" http://127.0.0.1:8888/api/video)

#  echo $upload_response

  # Extract the ID from the response
  id=$(echo $upload_response | jq -r '.status' | cut -d'/' -f3)

  success=0

  for ((i=1;i<=retries;i++));
  do
    sleep 1

    # Step 2: Send a GET request to check the status
    status_response=$(curl -s http://localhost:8888/api/status/$id)

    # Print the final status response
    echo "Status Response: "
    echo "$status_response" | jq .

    status=$(echo "$status_response" | jq  -r .status)


    if [ "$status" = "done" ]; then
      echo "request completed successfully"
      # go to next video
      success=1
      break


    elif [ "$status" = "error" ]; then
      echo "request contained an error"
      exit 1
    fi
  done


  if [ $success -eq 0 ]; then
    echo "request could not be completed"
    exit 1
  fi

done
