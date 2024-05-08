all: validate-project setup-tools generate-grpc # generate-http | To run everything :) | Я тут заигрался и сделал генерацию для dart и web (ts), но чтобы не мешало закомментил

# Base Paths
PROTO_DIR=api/proto
PROTO_FILES=$(shell find $(PROTO_DIR) -name '*.proto')
OPENAPI_OUT_DIR=api/openapiv2
OUT_DIR=pkg/api
GRPC_OUT_DIR=$(OUT_DIR)/grpc
HTTP_OUT_DIR=$(OUT_DIR)/http


# Generated Paths
GO_OUT_DIR=golang
DART_OUT_DIR=dart
TYPESCRIPT_OUT_DIR=typescript


# Vendoring Paths
VENDOR_DIR=vendor.protobuf
GOOGLE_APIS_DIR=$(VENDOR_DIR)/google
GRPC_GATEWAY_DIR=$(VENDOR_DIR)/grpc-gateway




CREATE_GRPC_OUT_DIR:
	@echo Creating directory $(GRPC_OUT_DIR)
	@mkdir -p $(GRPC_OUT_DIR)

CREATE_GO_GRPC_OUT_DIR:
	@echo "Creating directory $(GRPC_OUT_DIR)/$(GO_OUT_DIR)"
	@mkdir -p $(GRPC_OUT_DIR)/$(GO_OUT_DIR)

CREATE_DART_GRPC_OUT_DIR:
	@echo Creating directory $(GRPC_OUT_DIR)/$(DART_OUT_DIR)
	@mkdir -p $(GRPC_OUT_DIR)/$(DART_OUT_DIR)

CREATE_OPENAPI_OUT_DIR:
	@echo Creating directory $(OPENAPI_OUT_DIR)
	@mkdir -p $(OPENAPI_OUT_DIR)




# Installing Tools For Code Generation; Using Golang
setup-golang-tools:
	@echo "Setup golang tools"
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	@go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

# Installing Tools For Code Generation; Using Dart
setup-dart-tools:
	@echo "Setup dart tools"
	@dart pub global activate protoc_plugin


# Installing Tools For Code Generation; Using NPM
setup-npm-tools:
	@echo "Setup openapi generator"
	@npm install -g @openapitools/openapi-generator-cli

# Setup All Tools
setup-tools: validate-project vendor setup-golang-tools #  setup-dart-tools setup-npm-tools

# Vendoring Dependencies
vendor: .vendor-reset .vendor-googleapis .vendor-google-protobuf .vendor-protoc-gen-openapiv2 .vendor-protovalidate .vendor-tidy
#	@echo "Vendoring dependencies"
#	@echo "After vendoring completes, if you use jetbrains ide, enhance IDE dependency recognition with this quick setup: Go to Settings > Languages & Frameworks > Protocol Buffers > Import Paths, and add the path to your new dependency folder."
#	@if [ ! -d "$(VENDOR_DIR)" ]; then \
#  		echo "Creating vendor directory"; \
#		mkdir -p $(VENDOR_DIR); \
#	fi
#	@if [ ! -d "$(GOOGLE_APIS_DIR)" ] || [ -z "$$(ls -A $(GOOGLE_APIS_DIR) 2>/dev/null)"  ]; then \
#  		echo "Vendoring googleapis"; \
##		git clone --depth=1 https://github.com/googleapis/googleapis.git $(VENDOR_DIR)/googleapis; \
##		mv $(VENDOR_DIR)/googleapis/google $(VENDOR_DIR); \
##		rm -rf $(VENDOR_DIR)/googleapis; \
#		echo "Done vendoring googleapis"; \
#	fi
#	@if [ ! -d "$(GRPC_GATEWAY_DIR)" ] || [ -z "$$(ls -A $(GRPC_GATEWAY_DIR) 2>/dev/null)" ]; then \
#  		echo "Vendoring grpc-gateway"; \
#		git clone https://github.com/grpc-ecosystem/grpc-gateway.git $(GRPC_GATEWAY_DIR); \
#		mv $(VENDOR_DIR)/grpc/google $(VENDOR_DIR); \
#        rm -rf $(VENDOR_DIR)/googleapis; \
#	fi
#	@echo "Dependencies vendored"

# Generating Libraries From Protobuf
generate-grpc: CREATE_GRPC_OUT_DIR CREATE_OPENAPI_OUT_DIR CREATE_GO_GRPC_OUT_DIR # CREATE_DART_GRPC_OUT_DIR
	@echo "Generating gRPC libraries..."
	@protoc -I$(PROTO_DIR) \
	-I$(VENDOR_DIR) \
	--go_out $(GRPC_OUT_DIR)/$(GO_OUT_DIR) \
	--go_opt paths=source_relative \
	--go-grpc_out $(GRPC_OUT_DIR)/$(GO_OUT_DIR) \
	--go-grpc_opt paths=source_relative \
	--grpc-gateway_out $(GRPC_OUT_DIR)/$(GO_OUT_DIR) \
	--grpc-gateway_opt paths=source_relative \
	--openapiv2_out $(OPENAPI_OUT_DIR) \
	--openapiv2_opt use_go_templates=true,preserve_rpc_order=true \
	$(PROTO_FILES)

#	--dart_out $(GRPC_OUT_DIR)/$(DART_OUT_DIR) \



# Generating Libraries From Openapi
generate-http:
	@echo "Generating HTTP Libraries From Openapi...";
	@SWAGGER_FILES=$$(find $(OPENAPI_OUT_DIR) -name "*.swagger.json"); \
	for file in $$SWAGGER_FILES; do \
		name=$$(basename $$file .swagger.json); \
		echo "Generating for $$file..."; \
		echo "Generating typescript library"; \
		openapi-generator-cli generate \
			-i $$file \
			-g typescript-axios \
			-o $(HTTP_OUT_DIR)/$(TYPESCRIPT_OUT_DIR)/$$name \
			--additional-properties usePromises=true \
			--additional-properties useES6=true; \
		echo "Generating dart library"; \
		openapi-generator-cli generate \
			-i $$file \
			-g dart-dio \
			-o $(HTTP_OUT_DIR)/$(DART_OUT_DIR)/$$name; \
		echo "Generating golang library"; \
		openapi-generator-cli generate \
			-i $$file \
			-g  go \
			-o $(HTTP_OUT_DIR)/$(GO_OUT_DIR)/$$name; \
	done



# Validate api Folder, wich
validate-project:
	@echo "Validating project..."
	@if [ ! -d $(PROTO_DIR) ]; then \
		echo "Protobuf directory not found !!!!!!"; \
		echo "Please, place your proto files in the $(PROTO_DIR) directory !!!!!!"; \
		exit 1; \
	fi
	@echo "Project Is Valid!"


# Vendor zone

.vendor-reset:
	@echo "Reset vendor"
	@rm -rf $(VENDOR_DIR)
	@mkdir -p $(VENDOR_DIR)

.vendor-tidy:
	@echo "Vendor Tidy"
	@find $(VENDOR_DIR) -type f ! -name "*.proto" -delete
	@find $(VENDOR_DIR) -empty -type d -delete

# Устанавливаем proto описания google/protobuf
.vendor-google-protobuf:
	@echo "Vendoring google-protobuf"
	@git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf $(VENDOR_DIR)/protobuf &&\
	cd $(VENDOR_DIR)/protobuf &&\
	git sparse-checkout set --no-cone src/google/protobuf &&\
	git checkout
	@mkdir -p $(VENDOR_DIR)/google
	@mv $(VENDOR_DIR)/protobuf/src/google/protobuf $(VENDOR_DIR)/google
	@rm -rf $(VENDOR_DIR)/protobuf

# Устанавливаем proto описания validate
.vendor-protovalidate:
	@echo "Vendoring protovalidate"
	@git clone -b main --single-branch --depth=1 --filter=tree:0 \
		https://github.com/bufbuild/protovalidate $(VENDOR_DIR)/protovalidate && \
	cd $(VENDOR_DIR)/protovalidate
	@git checkout
	@mv $(VENDOR_DIR)/protovalidate/proto/protovalidate/buf $(VENDOR_DIR)
	@rm -rf $(VENDOR_DIR)/protovalidate

# Устанавливаем proto описания google/api
.vendor-googleapis:
	@echo "Vendoring googleapis"
	@git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/googleapis/googleapis $(VENDOR_DIR)/googleapis &&\
	cd $(VENDOR_DIR)/googleapis &&\
	git checkout
	@mv $(VENDOR_DIR)/googleapis/google $(VENDOR_DIR)
	@rm -rf $(VENDOR_DIR)/googleapis

# Устанавливаем proto описания protoc-gen-openapiv2/options
.vendor-protoc-gen-openapiv2:
	@echo "Vendoring protoc-gen-openapiv2"
	@git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/grpc-ecosystem/grpc-gateway $(VENDOR_DIR)/grpc-gateway && \
 	cd $(VENDOR_DIR)/grpc-gateway && \
	git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
	git checkout
	@mkdir -p $(VENDOR_DIR)/protoc-gen-openapiv2
	@mv $(VENDOR_DIR)/grpc-gateway/protoc-gen-openapiv2/options $(VENDOR_DIR)/protoc-gen-openapiv2
	@rm -rf $(VENDOR_DIR)/grpc-gateway
