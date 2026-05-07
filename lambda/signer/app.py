import base64
import json
import os
import secrets
import struct
import time


CERT_TYPE_USER = 1
DEFAULT_DURATION_SECONDS = int(os.environ.get("DEFAULT_DURATION_SECONDS", "1800"))
MAX_DURATION_SECONDS = int(os.environ.get("MAX_DURATION_SECONDS", "14400"))
CLOCK_SKEW_SECONDS = 60
CERT_KEY_TYPE = "ssh-ed25519-cert-v01@openssh.com"
ED25519_KEY_TYPE = "ssh-ed25519"

_dynamodb = None
_ssm = None
_ca_key = None


def handler(event, context):
    try:
        cognito_sub, email = authenticated_user(event)
        request = parse_request(event)
        principals = authorised_principals(cognito_sub, email, request["duration"])

        serial = secrets.randbits(63)
        now = int(time.time())
        valid_after = now - CLOCK_SKEW_SECONDS
        valid_before = now + request["duration"]
        key_id = f"trustssh:{cognito_sub}:{now}"

        certificate = sign_user_certificate(
            request["public_key"],
            request["comment"],
            load_ca_key(),
            serial,
            key_id,
            principals,
            valid_after,
            valid_before,
        )

        valid_until = iso8601(valid_before)
        audit(cognito_sub, email, principals, serial, valid_until, "issued")

        return json_response(
            200,
            {
                "certificate": certificate,
                "valid_until": valid_until,
                "principals": principals,
                "serial": serial,
            },
        )
    except AuthError as exc:
        return error_response(exc.status_code, exc.code, exc.message)
    except RequestError as exc:
        return error_response(400, "invalid_request", str(exc))
    except Exception:
        return error_response(500, "internal_error", "Certificate issuance failed")


class AuthError(Exception):
    def __init__(self, status_code, code, message):
        self.status_code = status_code
        self.code = code
        self.message = message


class RequestError(ValueError):
    pass


def authenticated_user(event):
    claims = event.get("requestContext", {}).get("authorizer", {}).get("jwt", {}).get("claims", {})
    cognito_sub = claims.get("sub")
    if not cognito_sub:
        raise AuthError(401, "unauthenticated", "Authenticated Cognito subject was not found")
    return cognito_sub, claims.get("email", "")


def parse_request(event):
    raw_body = event.get("body") or "{}"
    if event.get("isBase64Encoded"):
        raw_body = base64.b64decode(raw_body).decode("utf-8")

    try:
        body = json.loads(raw_body)
    except json.JSONDecodeError:
        raise RequestError("Request body must be valid JSON")

    public_key = body.get("public_key")
    if not isinstance(public_key, str) or not public_key.strip():
        raise RequestError("public_key is required")
    try:
        public_key_bytes, comment = parse_ssh_ed25519_public_key(public_key)
    except ValueError as exc:
        raise RequestError(str(exc)) from exc

    duration = body.get("requested_duration_seconds", DEFAULT_DURATION_SECONDS)
    if not isinstance(duration, int):
        raise RequestError("requested_duration_seconds must be an integer")
    if duration <= 0:
        raise RequestError("requested_duration_seconds must be positive")

    return {
        "public_key": public_key_bytes,
        "comment": comment,
        "duration": duration,
    }


def authorised_principals(cognito_sub, email, requested_duration):
    mapping = load_mapping(cognito_sub)
    principals = (mapping or {}).get("ssh_principals") or []
    if not mapping or not mapping.get("enabled") or not principals:
        audit(cognito_sub, email, [], 0, "", "denied")
        raise AuthError(403, "not_authorised", "No SSH login mapping exists for this user")

    max_mapping_duration = int(mapping.get("max_duration_seconds") or DEFAULT_DURATION_SECONDS)
    if requested_duration > min(max_mapping_duration, MAX_DURATION_SECONDS):
        raise RequestError("Requested certificate duration is too long")
    return principals


def load_mapping(cognito_sub):
    table_name = os.environ["USER_MAPPING_TABLE"]
    table = dynamodb().Table(table_name)
    response = table.get_item(Key={"cognito_sub": cognito_sub})
    return response.get("Item")


def load_ca_key():
    global _ca_key
    if _ca_key is not None:
        return _ca_key

    parameter_name = os.environ["CA_PRIVATE_KEY_PARAMETER_NAME"]
    response = ssm().get_parameter(Name=parameter_name, WithDecryption=True)
    value = response["Parameter"]["Value"]
    if "BEGIN OPENSSH PRIVATE KEY" not in value:
        raise ValueError("CA private key parameter is not an OpenSSH private key")
    _ca_key = parse_openssh_ed25519_private_key(value)
    return _ca_key


def audit(cognito_sub, email, principals, serial, valid_until, outcome):
    table_name = os.environ["AUDIT_EVENTS_TABLE"]
    table = dynamodb().Table(table_name)
    issued_at = iso8601(int(time.time()))
    table.put_item(
        Item={
            "cognito_sub": cognito_sub,
            "issued_at_serial": f"{issued_at}#{serial}",
            "email": email,
            "principals": principals,
            "serial": serial,
            "valid_until": valid_until,
            "outcome": outcome,
        }
    )


def dynamodb():
    global _dynamodb
    if _dynamodb is None:
        import boto3

        _dynamodb = boto3.resource("dynamodb")
    return _dynamodb


def ssm():
    global _ssm
    if _ssm is None:
        import boto3

        _ssm = boto3.client("ssm")
    return _ssm


def sign_user_certificate(user_public_key, comment, ca_key, serial, key_id, principals, valid_after, valid_before):
    nonce = secrets.token_bytes(32)
    cert_without_signature = b"".join(
        [
            ssh_string(CERT_KEY_TYPE.encode()),
            ssh_string(nonce),
            ssh_string(user_public_key),
            struct.pack(">Q", serial),
            struct.pack(">I", CERT_TYPE_USER),
            ssh_string(key_id.encode()),
            ssh_string(b"".join(ssh_string(p.encode()) for p in principals)),
            struct.pack(">Q", valid_after),
            struct.pack(">Q", valid_before),
            ssh_string(b""),
            ssh_string(default_extensions()),
            ssh_string(b""),
            ssh_string(ca_key["public_blob"]),
        ]
    )

    signature = ca_key["signing_key"].sign(cert_without_signature).signature
    signature_blob = ssh_string(ED25519_KEY_TYPE.encode()) + ssh_string(signature)
    cert_blob = cert_without_signature + ssh_string(signature_blob)
    cert_b64 = base64.b64encode(cert_blob).decode("ascii")
    cert_comment = comment or "trustssh"
    return f"{CERT_KEY_TYPE} {cert_b64} {cert_comment}"


def parse_ssh_ed25519_public_key(public_key_line):
    parts = public_key_line.strip().split()
    if len(parts) < 2:
        raise ValueError("public_key must be an OpenSSH public key")
    if parts[0] != ED25519_KEY_TYPE:
        raise ValueError("Only ssh-ed25519 public keys are supported")

    try:
        blob = base64.b64decode(parts[1].encode("ascii"), validate=True)
    except Exception as exc:
        raise ValueError("public_key has invalid base64") from exc

    key_type, public_key = parse_public_key_blob(blob)
    if key_type != ED25519_KEY_TYPE or len(public_key) != 32:
        raise ValueError("public_key is not a valid ssh-ed25519 key")

    comment = " ".join(parts[2:]) if len(parts) > 2 else ""
    return public_key, comment


def parse_public_key_blob(blob):
    reader = SSHReader(blob)
    key_type = reader.read_string().decode("ascii")
    public_key = reader.read_string()
    reader.expect_end()
    return key_type, public_key


def parse_openssh_ed25519_private_key(private_key_text):
    payload = "".join(
        line.strip()
        for line in private_key_text.strip().splitlines()
        if not line.startswith("-----")
    )
    raw = base64.b64decode(payload)
    marker = b"openssh-key-v1\x00"
    if not raw.startswith(marker):
        raise ValueError("CA private key is not in OpenSSH private key format")

    reader = SSHReader(raw[len(marker) :])
    cipher_name = reader.read_string().decode("ascii")
    kdf_name = reader.read_string().decode("ascii")
    reader.read_string()
    key_count = reader.read_uint32()
    if cipher_name != "none" or kdf_name != "none":
        raise ValueError("Encrypted CA private keys are not supported by the MVP Lambda")
    if key_count != 1:
        raise ValueError("CA private key must contain exactly one key")

    ca_public_blob = reader.read_string()
    private_block = SSHReader(reader.read_string())
    check_1 = private_block.read_uint32()
    check_2 = private_block.read_uint32()
    if check_1 != check_2:
        raise ValueError("CA private key checkints do not match")

    key_type = private_block.read_string().decode("ascii")
    if key_type != ED25519_KEY_TYPE:
        raise ValueError("Only ssh-ed25519 CA private keys are supported")

    public_key = private_block.read_string()
    private_key = private_block.read_string()
    if len(public_key) != 32 or len(private_key) != 64:
        raise ValueError("CA private key has invalid Ed25519 key material")

    if private_key[32:] != public_key:
        raise ValueError("CA private key public component does not match")

    return {
        "public_blob": ca_public_blob,
        "signing_key": signing_key(private_key[:32]),
    }


class SSHReader:
    def __init__(self, data):
        self.data = data
        self.offset = 0

    def read_uint32(self):
        self._require(4)
        value = struct.unpack(">I", self.data[self.offset : self.offset + 4])[0]
        self.offset += 4
        return value

    def read_string(self):
        length = self.read_uint32()
        self._require(length)
        value = self.data[self.offset : self.offset + length]
        self.offset += length
        return value

    def expect_end(self):
        if self.offset != len(self.data):
            raise ValueError("Unexpected trailing data in SSH structure")

    def _require(self, length):
        if self.offset + length > len(self.data):
            raise ValueError("Unexpected end of SSH structure")


def ssh_string(value):
    return struct.pack(">I", len(value)) + value


def default_extensions():
    names = [
        b"permit-X11-forwarding",
        b"permit-agent-forwarding",
        b"permit-port-forwarding",
        b"permit-pty",
        b"permit-user-rc",
    ]
    return b"".join(ssh_string(name) + ssh_string(b"") for name in names)


def json_response(status_code, body):
    return {
        "statusCode": status_code,
        "headers": {"content-type": "application/json"},
        "body": json.dumps(body),
    }


def error_response(status_code, code, message):
    return json_response(status_code, {"error": code, "message": message})


def iso8601(timestamp):
    return time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime(timestamp))


def signing_key(seed):
    from nacl.signing import SigningKey

    return SigningKey(seed)
