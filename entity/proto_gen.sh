#!/usr/bin/env bash

protoc -I . entity.proto --go_out=plugins=grpc:.
python3 -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. entity.proto