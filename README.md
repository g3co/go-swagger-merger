#Swagger merger

Merging few swagger yaml files into the one.

Install the command line tool first.

``
 go get github.com/g3co/go-swagger-merger
``

The command below will merge **/data/swagger1.yaml** **/data/swagger2.yaml** and save result file in the **/data/swagger.yaml**

``
 go-swagger-merger -o /data/swagger.yaml /data/swagger1.yaml /data/swagger2.yaml
``
