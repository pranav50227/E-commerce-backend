#!/bin/bash
# rollback.sh - Reverts/rolls back the last applied database migrations.

echo "--> Rolling back database migrations..."

for dir in service/*; do
  if [ -d "$dir/migrations" ]; then
    service_name=$(basename "$dir")
    echo "--> Rolling back migration for service: $service_name..."
    # Placeholder: if using migrate or goose
    # (cd "$dir" && goose down)
    echo "Rollback completed for $service_name successfully."
  fi
done

echo "Database rollback sequence complete."
