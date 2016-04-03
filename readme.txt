index: curl -i http://localhost:8080/locations

create: curl -i -H "Content-Type: application/json" -d '{"name": "Shanghai"}' http://localhost:8080/location

delete: curl -i -X DELETE "http://localhost:8080/location/Shanghai"

