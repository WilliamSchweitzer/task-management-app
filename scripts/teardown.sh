#!/bin/bash
set -e

echo "⚠️  This will destroy all infrastructure. Are you sure? (yes/no)"
read CONFIRM

if [ "$CONFIRM" != "yes" ]; then
  echo "Cancelled."
  exit 0
fi

cd "$(dirname "$0")/../terraform"

echo "🗑️  Destroying infrastructure..."
terraform destroy -auto-approve

echo "✅ All resources destroyed!"