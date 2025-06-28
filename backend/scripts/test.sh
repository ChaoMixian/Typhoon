#!/bin/bash

# This script runs Go tests for the Typhoon backend.

echo "Changing to backend directory..."
# Assuming this script is run from the repository root.
# If it's already in backend/scripts/, then 'cd ..' might be more appropriate.
# For now, let's assume it could be run from root or backend/scripts.
if [ -d "backend" ]; then
  cd backend || exit 1
else
  # Potentially already in backend directory or backend/scripts
  # Check for go.mod to be more certain we are in the Go project root (backend)
  if [ ! -f "go.mod" ]; then
    echo "Error: go.mod not found. Please run this script from the repository root or the backend directory."
    exit 1
  fi
fi

echo "Running Go tests..."
go test -v -coverprofile=coverage.out ./...

# Check if tests were successful
if [ $? -ne 0 ]; then
  echo "Go tests failed."
  exit 1
fi

echo "Go tests passed."
echo "To view coverage report, run: go tool cover -html=coverage.out"

# Optional: Static analysis checks (can be added here)
# echo "Running go vet..."
# go vet ./...
# if [ $? -ne 0 ]; then
#   echo "go vet found issues."
#   # exit 1 # Optionally exit on vet issues
# fi

# echo "Running staticcheck..."
# staticcheck ./... # Requires staticcheck to be installed
# if [ $? -ne 0 ]; then
#   echo "staticcheck found issues."
#   # exit 1 # Optionally exit on staticcheck issues
# fi

echo "Backend test script completed."
