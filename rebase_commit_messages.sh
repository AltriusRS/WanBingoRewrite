#!/bin/bash
set -euo pipefail

if [ $# -lt 1 ]; then
  echo "Usage: $0 commits.csv"
  exit 1
fi

CSV_FILE="$1"

if [ ! -f "$CSV_FILE" ]; then
  echo "Error: File not found: $CSV_FILE"
  exit 1
fi

# Create a backup branch before rewriting
git branch "backup-before-rewrite-$(date +%s)"

TMP_MAP="/tmp/commit-map.txt"
> "$TMP_MAP"

# Read CSV safely (skip header line)
tail -n +2 "$CSV_FILE" | while IFS=',' read -r commit_hash original_message new_message; do
  # Handle possible quoted fields and commas
  commit_hash=$(echo "$commit_hash" | sed 's/^"//; s/"$//; s/^[[:space:]]*//; s/[[:space:]]*$//')
  new_message=$(echo "$new_message" | sed 's/^"//; s/"$//; s/^[[:space:]]*//; s/[[:space:]]*$//')

  if [ -n "$commit_hash" ] && [ -n "$new_message" ]; then
    echo "$commit_hash|$new_message" >> "$TMP_MAP"
  fi
done

echo "✅ Commit message map built: $(wc -l < "$TMP_MAP") entries."

# Rewrite commit messages using the map
git filter-branch --msg-filter '
while IFS="|" read -r commit_hash new_message; do
  if [ "$GIT_COMMIT" = "$commit_hash" ]; then
    echo "$new_message"
    exit 0
  fi
done < /tmp/commit-map.txt
cat
' -- --all

echo "✅ Done! Commit messages rewritten."
echo "💡 If anything looks wrong, you can revert with:"
echo "   git reset --hard backup-before-rewrite-*"
