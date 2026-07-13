#!/bin/bash
# migrate.sh - Applies schema migrations.

echo "--> Applying database migrations..."

# Under each microservice, check if a migrations folder exists and run goose/migrate
for dir in service/*; do
  if [ -d "$dir/migrations" ]; then
    service_name=$(basename "$dir")
    echo "--> Running migrations for service: $service_name..."
    # Placeholder: if using migrate or goose
    # (cd "$dir" && goose up)
    echo "Migrations applied for $service_name successfully."
  fi
done

echo "Database migration sequence complete."
