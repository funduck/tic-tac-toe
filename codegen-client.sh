function gen() {
  echo "Generating $GEN_MODULE API client for url $GEN_API_URL"
  
  curl -sL $GEN_API_URL -o $GEN_OPENAPI_FILE
  
  rm -rf generated/$GEN_MODULE

  docker run --rm \
  -u $(id -u):$(id -g) \
  -v ${PWD}:/local openapitools/openapi-generator-cli generate \
  -i /local/$GEN_OPENAPI_FILE \
  -g go \
  -o /local/generated/$GEN_MODULE

  if [ $? -ne 0 ]; then
    echo "Error generating $GEN_MODULE API client"
    exit 1
  fi

  rm $GEN_OPENAPI_FILE
}

CODEGEN_SERVER_URL=http://localhost:8080

GEN_API_URL=$CODEGEN_SERVER_URL/swagger/doc.json
GEN_OPENAPI_FILE=server.json
GEN_MODULE=server
gen