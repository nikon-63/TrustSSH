import os
import subprocess
import tempfile
import unittest

import app


class OpenSSHCertificateTests(unittest.TestCase):
    def test_sign_user_certificate_generates_valid_openssh_certificate(self):
        with tempfile.TemporaryDirectory() as tmp:
            ca_key = os.path.join(tmp, "ca")
            user_key = os.path.join(tmp, "user")
            cert_path = user_key + "-cert.pub"

            subprocess.run(["ssh-keygen", "-t", "ed25519", "-f", ca_key, "-N", "", "-C", "ca"], check=True, capture_output=True)
            subprocess.run(["ssh-keygen", "-t", "ed25519", "-f", user_key, "-N", "", "-C", "user"], check=True, capture_output=True)

            with open(ca_key, "r", encoding="utf-8") as fh:
                ca_key_material = fh.read()
            ca_key = app.parse_openssh_ed25519_private_key(ca_key_material)
            with open(user_key + ".pub", "r", encoding="utf-8") as fh:
                user_public_key = fh.read()
            user_public_key, comment = app.parse_ssh_ed25519_public_key(user_public_key)

            certificate = app.sign_user_certificate(
                user_public_key,
                comment,
                ca_key,
                serial=12345,
                key_id="trustssh:test:1",
                principals=["ubuntu"],
                valid_after=0,
                valid_before=4102444800,
            )

            with open(cert_path, "w", encoding="utf-8") as fh:
                fh.write(certificate + "\n")

            result = subprocess.run(["ssh-keygen", "-L", "-f", cert_path], check=True, capture_output=True, text=True)
            self.assertIn("Type: ssh-ed25519-cert-v01@openssh.com user certificate", result.stdout)
            self.assertIn("Key ID: \"trustssh:test:1\"", result.stdout)
            self.assertIn("Serial: 12345", result.stdout)
            self.assertIn("ubuntu", result.stdout)


if __name__ == "__main__":
    unittest.main()
