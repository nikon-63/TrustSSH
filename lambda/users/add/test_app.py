import os
import json
import logging
import unittest
from unittest.mock import patch, MagicMock

os.environ['USER_POOL_ID'] = 'us-east-1_testpool'
os.environ['USER_MAPPING_TABLE'] = 'test-mappings-table'

import app

class TestAddUserLambda(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        # Suppress logging output during tests
        logging.getLogger().setLevel(logging.CRITICAL)

    def setUp(self):
        patcher = patch('app.cognito_client')
        self.mock_cognito = patcher.start()
        self.addCleanup(patcher.stop)
        
        self.mock_cognito.exceptions.UsernameExistsException = type('UsernameExistsException', (Exception,), {})

    @patch('app.dynamodb_client')
    def test_successful_creation(self, mock_dynamodb):
        # Mock Cognito response
        self.mock_cognito.admin_create_user.return_value = {
            'User': {
                'Attributes': [{'Name': 'sub', 'Value': 'test-sub-1234'}]
            }
        }
        
        event = {
            "email": "test@example.com",
            "ssh_principals": ["ubuntu", "admin"]
        }
        
        response = app.handler(event, None)
        
        self.assertEqual(response['statusCode'], 200)
        body = json.loads(response['body'])
        self.assertEqual(body['cognito_sub'], 'test-sub-1234')
        self.assertEqual(body['email'], 'test@example.com')
        
        # Ensure boto3 was called correctly
        self.mock_cognito.admin_create_user.assert_called_once()
        mock_dynamodb.put_item.assert_called_once()

    def test_missing_email(self):
        event = {"ssh_principals": ["ubuntu"]}
        response = app.handler(event, None)
        self.assertEqual(response['statusCode'], 400)
        self.assertIn('email is required', response['body'])

    def test_missing_principals(self):
        event = {"email": "test@example.com"}
        response = app.handler(event, None)
        self.assertEqual(response['statusCode'], 400)
        self.assertIn('ssh_principals is required', response['body'])

    def test_user_already_exists(self):
        self.mock_cognito.admin_create_user.side_effect = self.mock_cognito.exceptions.UsernameExistsException()
        
        event = {
            "email": "test@example.com",
            "ssh_principals": ["ubuntu"]
        }
        
        response = app.handler(event, None)
        self.assertEqual(response['statusCode'], 400)
        self.assertIn('already exists', response['body'])

    def test_missing_sub_in_response(self):
        self.mock_cognito.admin_create_user.return_value = {
            'User': {
                'Attributes': [{'Name': 'email', 'Value': 'test@example.com'}] # Missing sub
            }
        }
        
        event = {
            "email": "test@example.com",
            "ssh_principals": ["ubuntu"]
        }
        
        response = app.handler(event, None)
        self.assertEqual(response['statusCode'], 500)
        self.assertIn('Failed to get sub', response['body'])

    def test_general_exception(self):
        self.mock_cognito.admin_create_user.side_effect = Exception("General error")
        
        event = {
            "email": "test@example.com",
            "ssh_principals": ["ubuntu"]
        }
        
        response = app.handler(event, None)
        self.assertEqual(response['statusCode'], 500)
        self.assertIn('Internal server error', response['body'])

if __name__ == '__main__':
    unittest.main()
