#!/bin/bash

set -e

# Directory to store sorted YAML files by namespace
SORTED_DIR="sorted_by_namespace"
mkdir -p "$SORTED_DIR"

# Find all YAML files in subdirectories
YAML_FILES=$(find . -type f -name "*.yaml")

# Extract resources by namespace and write them to separate files with proper separation
echo "Extracting resources by namespace..."
for namespace in $(echo "$YAML_FILES" | xargs yq eval '.metadata.namespace' | sort | uniq); do
  if [ "$namespace" != "null" ]; then
    first_resource=true
    for yaml_file in $YAML_FILES; do
      if yq eval "select(.metadata.namespace == \"$namespace\")" "$yaml_file" > /dev/null; then
        if [ "$first_resource" = true ]; then
          first_resource=false
        else
          # Add document separator between resources
          echo "---" >> "$SORTED_DIR/$namespace.yaml"
        fi
        # Append the resource content
        yq eval "select(.metadata.namespace == \"$namespace\")" "$yaml_file" >> "$SORTED_DIR/$namespace.yaml"
      fi
    done
    echo "Resources for namespace '$namespace' saved to '$SORTED_DIR/$namespace.yaml'"
  else
    echo "Skipping resources with no namespace."
  fi
done

# Check if any YAML files were created
if [ -z "$(ls -A $SORTED_DIR)" ]; then
  echo "No namespace-specific YAML files found. Exiting."
  exit 1
fi

# Run kube-score for each namespace YAML file
echo "Running kube-score on each namespace..."
exit_code=0
for file in $SORTED_DIR/*.yaml; do
  namespace=$(basename "$file" .yaml)
  echo "Running kube-score for namespace: $namespace..."
  
  # Run kube-score and capture the exit code
  kube-score score --output-format ci "$file" --exit-one-on-warning || exit_code=1
done

# Exit with the appropriate code
exit $exit_code
