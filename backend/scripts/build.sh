#!/bin/bash

# This script builds the Typhoon Go application, moves it to a temporary
# run directory, and executes it. It also handles cleanup.

# Determine the backend directory (where go.mod is)
# This handles being run from repo root or from backend/scripts itself
SCRIPT_DIR="$(cd "\$(dirname "\${BASH_SOURCE[0]}")" &> /dev/null && pwd)"
BACKEND_DIR=""

if [ -f "go.mod" ]; then # If run from backend/ directory
  BACKEND_DIR=$(pwd)
elif [ -f "../go.mod" ]; then # If run from backend/scripts/ directory
  BACKEND_DIR=$(cd ".." && pwd)
elif [ -f "backend/go.mod" ]; then # If run from repo root
  BACKEND_DIR=$(cd "backend" && pwd)
else
  echo "Error: Could not determine backend directory (where go.mod is located)."
  echo "Please run this script from the repository root, the backend/ directory, or the backend/scripts/ directory."
  exit 1
fi

echo "Backend directory identified as: ${BACKEND_DIR}"

# Configuration
BINARY_NAME="Typhoon"
BUILD_OUTPUT_DIR="${BACKEND_DIR}/build_temp" # Temporary build output within backend
RUN_DIR="${BACKEND_DIR}/run_temp"          # Temporary run directory within backend
BINARY_PATH_BUILD="${BUILD_OUTPUT_DIR}/${BINARY_NAME}"
BINARY_PATH_RUN="${RUN_DIR}/${BINARY_NAME}"

# --- Cleanup Function ---
echo "Setting up cleanup trap..."
cleanup() {
  echo "Cleaning up..."
  # Stop Typhoon if it's running (this is tricky, might need PID management or specific stop command)
  # For now, we'll just remove the temp directories.
  # If Typhoon writes a PID file, we could use that here.
  if [ -d "${RUN_DIR}" ]; then
    # Try to kill the process if it's still running by name (platform dependent and risky)
    # pkill -f "${BINARY_NAME}" # Be careful with this, might kill other processes
    echo "Removing temporary run directory: ${RUN_DIR}"
    rm -rf "${RUN_DIR}"
  fi
  if [ -d "${BUILD_OUTPUT_DIR}" ]; then
    echo "Removing temporary build directory: ${BUILD_OUTPUT_DIR}"
    rm -rf "${BUILD_OUTPUT_DIR}"
  fi
  echo "Cleanup finished."
}

# Trap EXIT signal to run cleanup function
# This ensures cleanup happens on script exit, whether normal or due to error (Ctrl+C, etc.)
trap cleanup EXIT

# --- Main Script --- 

# 1. Initial Cleanup (in case previous run didn't clean up properly)
echo "Performing initial cleanup of old temp directories (if any)..."
if [ -d "${RUN_DIR}" ]; then
  rm -rf "${RUN_DIR}"
fi
if [ -d "${BUILD_OUTPUT_DIR}" ]; then
  rm -rf "${BUILD_OUTPUT_DIR}"
fi

# 2. Create temporary directories
echo "Creating temporary directories..."
mkdir -p "${BUILD_OUTPUT_DIR}"
mkdir -p "${RUN_DIR}"

# 3. Build Go project
echo "Building the project in ${BACKEND_DIR}..."
(cd "${BACKEND_DIR}" && go build -v -o "${BINARY_PATH_BUILD}" .)

# Check if build was successful
if [ $? -ne 0 ]; then
  echo "Build failed. Exiting..."
  exit 1 # Cleanup will be triggered by trap
fi
echo "Build successful. Binary created at: ${BINARY_PATH_BUILD}"

# 4. Move binary to run directory
echo "Moving binary to run directory: ${RUN_DIR}"
mv "${BINARY_PATH_BUILD}" "${BINARY_PATH_RUN}"
if [ $? -ne 0 ]; then
  echo "Failed to move binary. Exiting..."
  exit 1
fi
cp -r "${BACKEND_DIR}/tests/mihomo" "${RUN_DIR}/" # Copy mihomo directory


# 5. Copy/Create necessary config files in the RUN_DIR if needed.
# Typhoon's default behavior is to create config.json in the executable's directory if not found.
# It also creates mihomo/config/default/config.yaml relative to executable.
# So, just running it from RUN_DIR should be enough for it to create its default files there.
echo "Typhoon will create its default config files in ${RUN_DIR} if they don't exist."

# 6. Run the program
echo "Running Typhoon from ${RUN_DIR}... (Press Ctrl+C to stop)"
(cd "${RUN_DIR}" && ./${BINARY_NAME})

# Note: The script will wait here until Typhoon is terminated (e.g., by Ctrl+C).
# The 'trap cleanup EXIT' will handle cleanup afterwards.

echo "Typhoon execution finished (or was interrupted)."
# Explicit exit, cleanup will still run due to trap.
exit 0
