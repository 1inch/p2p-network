protobuf-gen:
	protoc \
		--plugin="protoc-gen-es=./node_modules/.bin/protoc-gen-es" \
		--proto_path=./../proto \
		--es_out=./gen \
		--es_opt target=ts \
		./../proto/*.proto
