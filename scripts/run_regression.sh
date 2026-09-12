#!/usr/bin/env bash
# ==============================================================================
# Cafe ERP System - Full Regression Test Suite Orchestrator (Bash/Linux/CI)
# ==============================================================================

set -eo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
BACKEND_DIR="${PROJECT_ROOT}/backend"
FRONTEND_DIR="${PROJECT_ROOT}/frontend"

START_TIME=$(date +%s)

echo "======================================================================"
echo "🚀 CAFE ERP SYSTEM - REGRESSION TESTING SUITE (CI/CD)"
echo "======================================================================"
echo "Timestamp: $(date)"
echo ""

run_step() {
    local name="$1"
    shift
    echo "----------------------------------------------------------------------"
    echo "▶ Running ${name}..."
    local step_start=$(date +%s)
    
    if "$@"; then
        local step_end=$(date +%s)
        local dur=$((step_end - step_start))
        echo "✅ ${name} PASSED (${dur}s)"
    else
        local step_end=$(date +%s)
        local dur=$((step_end - step_start))
        echo "❌ ${name} FAILED (${dur}s)"
        exit 1
    fi
}

# 1. Smoke Testing
run_step "1. Smoke & Sanity Testing" bash -c "cd '${BACKEND_DIR}' && go test -count=1 -v ./tests/smoke/..."

# 2. Unit Testing
run_step "2. Unit Testing Suite" bash -c "cd '${BACKEND_DIR}' && go test -count=1 -v ./pkg/crypto/..."

# 3. Integration Testing
run_step "3. Integration Testing Suite" bash -c "cd '${BACKEND_DIR}' && go test -count=1 -v ./tests/integration/..."

# 4. E2E Testing
run_step "4. End-to-End (E2E) Testing Suite" bash -c "cd '${BACKEND_DIR}' && go test -count=1 -v ./tests/e2e/..."

# 5. Frontend Quality Gate
run_step "5. Frontend Quality Gate (Typecheck & Build)" bash -c "cd '${FRONTEND_DIR}' && npx vue-tsc --noEmit && npm run build"

END_TIME=$(date +%s)
TOTAL_DUR=$((END_TIME - START_TIME))

echo ""
echo "======================================================================"
echo "🎉 ALL REGRESSION SUITES PASSED in ${TOTAL_DUR}s"
echo "======================================================================"
exit 0

