response=$(curl -X GET -s -w "%{http_code}" http://localhost:8091/api/v1/ping)
if [ "$response" = "200" ]; then
    exit 0
else
    exit 1
fi