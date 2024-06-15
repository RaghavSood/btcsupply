#!/bin/bash
set -e

# Run litestream with your app as the subprocess.
exec /app/bin/litestream replicate -exec "/app/bin/btcsupply"
