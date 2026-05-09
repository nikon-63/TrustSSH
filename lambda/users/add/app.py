import os
import json
import logging
import boto3
from datetime import datetime, timezone

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
            
        email = body.get('email')
        ssh_principals = body.get('ssh_principals', [])
        max_duration_seconds = body.get('max_duration_seconds', 1800)
        
        if not email:
            return {"statusCode": 400, "body": json.dumps({"error": "email is required"})}
            
        if not ssh_principals:
            return {"statusCode": 400, "body": json.dumps({"error": "ssh_principals is required"})}
            
        response = cognito_client.admin_create_user(
            UserPoolId=USER_POOL_ID,
            Username=email,
            UserAttributes=[
                {'Name': 'email', 'Value': email},
                {'Name': 'email_verified', 'Value': 'true'}
            ],
            MessageAction='SUPPRESS'
        )
        
        sub = None
        for attr in response['User']['Attributes']:
            if attr['Name'] == 'sub':
                sub = attr['Value']
                break
                
        if not sub:
            return {"statusCode": 500, "body": json.dumps({"error": "Failed to get sub from newly created Cognito user"})}
            
        now = datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
        dynamodb_client.put_item(
            TableName=USER_MAPPING_TABLE,
            Item={
                'cognito_sub': {'S': sub},
                'email': {'S': email},
                'enabled': {'BOOL': True},
                'ssh_principals': {'L': [{'S': str(p)} for p in ssh_principals]},
                'max_duration_seconds': {'N': str(max_duration_seconds)},
                'created_at': {'S': now},
                'updated_at': {'S': now}
            }
        )
        
        return {
            "statusCode": 200, 
            "body": json.dumps({
                "message": "User created successfully", 
                "cognito_sub": sub,
                "email": email
            })
        }
    except cognito_client.exceptions.UsernameExistsException:
        logger.error(f"User {email} already exists")
        return {"statusCode": 400, "body": json.dumps({"error": "User already exists"})}
    except Exception as e:
        logger.error(f"Error creating user: {e}", exc_info=True)
        return {"statusCode": 500, "body": json.dumps({"error": "Internal server error"})}
