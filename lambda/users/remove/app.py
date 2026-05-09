import os
import json
import logging
import boto3

logger = logging.getLogger()
logger.setLevel(logging.INFO)

cognito_client = boto3.client('cognito-idp')
dynamodb_client = boto3.client('dynamodb')

USER_POOL_ID = os.environ.get('USER_POOL_ID')
USER_MAPPING_TABLE = os.environ.get('USER_MAPPING_TABLE')

def handler(event, context):
    try:
        # For if later this is behind API Gateway
        if 'body' in event and isinstance(event['body'], str):
            body = json.loads(event['body'])
        else:
            body = event
            
        sub = body.get('sub')
        
        if not sub:
            return {"statusCode": 400, "body": json.dumps({"error": "sub is required"})}
            
        try:
            cognito_client.admin_delete_user(
                UserPoolId=USER_POOL_ID,
                Username=sub
            )
            logger.info(f"Deleted user {sub} from Cognito")
        except cognito_client.exceptions.UserNotFoundException:
            logger.warning(f"User {sub} not found in Cognito, proceeding to remove from DynamoDB")
            
        dynamodb_client.delete_item(
            TableName=USER_MAPPING_TABLE,
            Key={
                'cognito_sub': {'S': sub}
            }
        )
        logger.info(f"Deleted user {sub} from DynamoDB")
        
        return {
            "statusCode": 200, 
            "body": json.dumps({
                "message": "User removed successfully", 
                "cognito_sub": sub
            })
        }
    except Exception as e:
        logger.error(f"Error removing user: {e}", exc_info=True)
        return {"statusCode": 500, "body": json.dumps({"error": "Internal server error"})}
