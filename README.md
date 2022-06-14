# Swagger merger

To merge a few swagger YAML files into one.

Install Go if you don't have one.

	https://golang.org/doc/install

Install the command line tool first.

	go get github.com/g3co/go-swagger-merger


The command below will merge ``/data/swagger1.yaml`` ``/data/swagger2.yaml`` and save result file in the ``/data/swagger.yaml``. The library supports more than two files to merge.

	go-swagger-merger -o /data/swagger.yaml /data/swagger1.yaml /data/swagger2.yaml


Attention. The order of the files is essential, and the following file overwrites the same fields from the previous file.
