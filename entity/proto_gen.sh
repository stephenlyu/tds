#!/usr/bin/env bash

protoc -I . entity.proto --go_out=plugins=grpc:.
