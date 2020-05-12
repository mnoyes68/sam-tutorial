set -e

# Burner account mnoyes-sam-test

# Make the bucket
# aws s3 mb s3://mnoyes-sam-test-app

# Building the stack
sam build
sam package --s3-bucket mnoyes-sam-test-app --output-template-file output.yaml
sam deploy --template-file output.yaml --stack-name mnoyes-test-stack-1 --capabilities CAPABILITY_IAM

# Inspecting the stack
# aws cloudformation describe-stacks --stack-name mnoyes-test-stack-1
# aws cloudformation describe-stacks --stack-name mnoyes-test-stack-1 --query Stacks[].Outputs

# aws cloudformation describe-stack-resources --stack-name mnoyes-test-stack-1


# To remove the bucket when finished
# aws s3 rm s3://mnoyes-test-app --recursive # old bucket
# aws s3 rm s3://mnoyes-sam-test-app --recursive

# To delete the stack
# aws cloudformation delete-stack --stack-name mnoyes-test-stack-1

# To get logs
# sam logs -n HelloWorldFunction --stack-name mnoyes-test-stack-1




