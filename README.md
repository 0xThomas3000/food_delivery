docker run -d --name mysql --privileged=true -e MYSQL_ROOT_PASSWORD="admin" -e MYSQL_USER="food_delivery" -e MYSQL_PASSWORD="12345678" -e MYSQL_DATABASE="food_delivery" -p 3307:3306 mysql:latest

Transport_layer: parse data from request/socker
Business_layer: Do some logic
Storage_layer: Integrate with DB