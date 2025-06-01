# 🧬 Swagger Merger

A simple CLI tool to merge multiple Swagger (OpenAPI) YAML and JSON files into a single spec.

### 🚀 Features

* Merge two or more Swagger files into one
* Supports overwrite logic based on file order
* Ideal for modular API definitions


### 🛠️ Installation

Make sure you have [Go installed](https://golang.org/doc/install), then run:

```bash
go get github.com/g3co/go-swagger-merger
```


### 📦 Usage

Merge two or more Swagger files into one:

```bash
go-swagger-merger -o /data/swagger.yaml /data/swagger1.yaml /data/swagger2.json
```

You can add as many input files as needed — supported formats include `.yaml`, `.yml`, and `.json`:

```bash
go-swagger-merger -o /data/swagger.json /data/swagger1.yaml /data/swagger2.yaml /data/swagger3.yaml
```


> ⚠️ **Note:** File order matters. Later files will overwrite conflicting fields from earlier ones.


### 🧪 Quick Test

Test assets are available in the repository to quickly try out the merger

