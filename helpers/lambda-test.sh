#!/usr/bin/env bash
set -Eeuo pipefail

ROOT_DIR="$(pwd)"
VENV_DIR="$ROOT_DIR/helpers/pytest-env"
REQUIREMENTS_FILE="$ROOT_DIR/helpers/requirements.txt"

TEST_TARGETS=(
    "lambda/users/add"
    "lambda/users/remove"
    "lambda/signer"
)

echo "Setting up Python virtual environment"

if [[ ! -d "$VENV_DIR" ]]; then
    python3 -m venv "$VENV_DIR"
fi

source "$VENV_DIR/bin/activate"

python3 -m pip install --upgrade pip
python3 -m pip install -r "$REQUIREMENTS_FILE"

export AWS_ACCESS_KEY_ID="test"
export AWS_SECRET_ACCESS_KEY="test"
export AWS_DEFAULT_REGION="eu-west-2"

echo "Running Lambda unit tests"

for target in "${TEST_TARGETS[@]}"; do
    echo
    echo "Testing: $target"

    cd "$ROOT_DIR/$target"
    python3 -m unittest test_app.py
done

echo
echo "All tests passed"