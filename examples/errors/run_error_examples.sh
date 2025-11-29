#!/bin/bash

# Script to run all Beeflang error examples
# This demonstrates the error handling system

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$SCRIPT_DIR/../.."

echo "=================================================="
echo "Beeflang Error Examples"
echo "=================================================="
echo ""
echo "These examples intentionally contain errors to"
echo "demonstrate the error handling system."
echo ""

# Array of error examples with descriptions
declare -a examples=(
  "type_mismatch.beef:Type Mismatch (INTEGER + BOOLEAN)"
  "undefined_variable.beef:Undefined Variable"
  "unknown_operator.beef:Unknown Operator (BOOLEAN + BOOLEAN)"
  "invalid_negation.beef:Invalid Negation (-BOOLEAN)"
  "string_type_mismatch.beef:String Type Mismatch (STRING + INTEGER)"
)

for example in "${examples[@]}"; do
  IFS=':' read -r file description <<< "$example"

  echo "--------------------------------------------------"
  echo "Example: $description"
  echo "File: $file"
  echo "--------------------------------------------------"

  # Run the example and capture output (expecting it to fail)
  if ! go run "$PROJECT_ROOT/main.go" "$SCRIPT_DIR/$file" 2>&1; then
    echo ""
  fi

  echo ""
done

echo "=================================================="
echo "All error examples completed!"
echo "=================================================="
echo ""
echo "Notice how each error includes:"
echo "  • Line number"
echo "  • Column number"
echo "  • Clear error message"
echo "  • Execution stops immediately"
echo ""
