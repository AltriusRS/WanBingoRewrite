#!/bin/bash

# Script to automatically rebase and update commit messages
# Run this on the meta-commit-cleanup branch

echo "Starting interactive rebase for the last 8 commits..."

git rebase -i HEAD~8 <<EOF
reword 358cec2 [SERVER] Bundle migrations in container
reword 042a038 [SERVER] [DB] Add migration runner
reword 4a97ddd [DB] Split schema into migration files
reword ba280a5 [SERVER] [FRONTEND] [DB] Implement tile suggestion system
reword a4e3965 [SERVER] [FRONTEND] Enable anonymous bingo gameplay
reword bef1fc1 [FRONTEND] [SERVER] Clean up unused files
reword d35baa6 [CI/CD] Fix build pipeline
reword 1baf6c7 [FRONTEND] [SERVER] Fix build issues and add CORS env
EOF

echo "Rebase completed. Check git log to verify changes."