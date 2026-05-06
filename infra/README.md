```bash
aws login --profile signin
aws sts get-caller-identity --profile signin
aws configure set region eu-west-2 --profile process
aws configure set credential_process \
  "aws configure export-credentials --profile signin --format process" \
  --profile process
export AWS_PROFILE=process
export AWS_SDK_LOAD_CONFIG=1
export AWS_EC2_METADATA_DISABLED=true
```