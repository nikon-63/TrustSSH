import os
import json
import logging
import unittest
from unittest.mock import patch, MagicMock

# Setup environment variables before importing app
os.environ['USER_POOL_ID'] = 'us-east-1_testpool'
os.environ['USER_MAPPING_TABLE'] = 'test-mappings-table'

import app

class TestRemoveUserLambda(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        # Suppress logging output during tests
        logging.getLogger().setLevel(logging.CRITICAL)

    def setUp(self):
        # Ensure boto3 exceptions used in try/except are proper Exception classes
        patcher = patch('app.cognito_client')
        self.mock_cognito = patcher.start()
        self.addCleanup(patcher.stop)
        
        # Setup specific exceptions for boto3
        self.mock_cognito.exceptions.UserNotFoundException = type('UserNotFoundException', (Exception,), {})

    @patch('app.dynamodb_client')
    def test_successful_removal(self, mock_dynamodb):
        event = {
            "sub": "test-sub-1234"
        }
        
        response = app.handler(event, None)
        
        self.assertEqual(response['statusCode'], 200)
        body = json.loads(response['body'])
        self.assertEqual(body['cognito_sub'], 'test-sub-1234')
        self.assertEqual(body['message'], 'User removed successfully')
        
        # Ensure boto3 was called correctly
        self.mock_cognito.admin_delete_user.assert_called_once_with(
            UserPoolId='us-east-1_testpool',
            Username='test-sub-1234'
        )
        mock_dynamodb.delete_item.assert_called_once()

    def test_missing_sub(self):
        event = {}
        response = app.handler(event, None)
        self.assertEqual(response['statusCode'], 400)
        self.assertIn('sub is required', response['body'])

    @patch('app.dynamodb_client')
    def test_user_not_found_in_cognito(self, mock_dynamodb):
        self.mock_cognito.admin_delete_user.side_effect = self.mock_cognito.exceptions.UserNotFoundException()
        
        event = {
            "sub": "test-sub-1234"
        }
        
        response = app.handler(event, None)
        
        # Even if user not in Cognito, we proceed to remove from DynamoDB and return 200
        self.assertEqual(response['statusCode'], 200)
        mock_dynamodb.delete_item.assert_called_once()

    def test_general_exception(self):
        # Trigger an unhandled Exception
        self.mock_cognito.admin_delete_user.side_effect = Exception("General error")
        
        event = {
            "sub": "test-sub-1234"
        }
        
        response = app.handler(event, None)
        
        self.assertEqual(response['statusCode'], 500)
        self.assertIn('Internal server error', response['body'])

if __name__ == '__main__':
    unittest.main()
