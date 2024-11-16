set -e

# Configure logging
LOG_FILE="/app/logs/app.log"
touch $LOG_FILE

# Start the application with logging
exec /app/apiserver 2>&1 | tee -a $LOG_FILE

---
# scripts/manage.sh
#!/bin/bash

function show_help {
    echo "Usage: ./manage.sh [command]"
    echo "Commands:"
    echo "  start       - Start all containers"
    echo "  stop        - Stop all containers"
    echo "  restart     - Restart all containers"
    echo "  logs        - Show logs from all containers"
    echo "  status      - Show container status"
    echo "  clean       - Remove all containers and images"
    echo "  build       - Rebuild all containers"
    echo "  test        - Run tests"
    echo "  debug       - Enter debug mode"
}

case "$1" in
    start)
        docker-compose up -d
        ;;
    stop)
        docker-compose down
        ;;
    restart)
        docker-compose down
        docker-compose up -d
        ;;
    logs)
        docker-compose logs -f
        ;;
    status)
        docker-compose ps
        ;;
    clean)
        docker-compose down -v
        docker system prune -af
        ;;
    build)
        docker-compose build --no-cache
        ;;
    test)
        docker-compose run --rm app go test ./...
        ;;
    debug)
        docker-compose exec app /bin/sh
        ;;
    *)
        show_help
        ;;
esac