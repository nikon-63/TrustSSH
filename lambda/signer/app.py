import json
import os


def handler(event, context):
    claims = (
        event.get("requestContext", {})
        .get("authorizer", {})
        .get("jwt", {})
        .get("claims", {})
    )

    cognito_sub = claims.get("sub", "")
    email = claims.get("email", "")

    body = {
        "error": "not_implemented",
        "message": "Certificate signing is not implemented yet",
        "authenticated": bool(cognito_sub),
        "cognito_sub": cognito_sub,
        "email": email,
        "tables": {
            "user_mappings": os.environ.get("USER_MAPPING_TABLE", ""),
            "audit_events": os.environ.get("AUDIT_EVENTS_TABLE", ""),
        },
    }

    return {
        "statusCode": 501,
        "headers": {
            "content-type": "application/json",
        },
        "body": json.dumps(body),
    }
