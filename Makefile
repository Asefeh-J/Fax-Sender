GO=go 
ARCH?=amd64 #for 32  bit execute : 386
EXPORT_DIR_NAME=source_code
EXPORT_PATH=$(PWD)/bin/$(EXPORT_DIR_NAME)

BASE_PATH=$(PWD)/build
##docker
DOCKER_GO_FILE_PATH=$(PWD)/bin/go1.21.4.linux-amd64.tar.gz
DOCKER_SOURCE_PATH=$(PWD)/bin/package_src/

# debian 
DOCKER_PATH_DEB=$(BASE_PATH)/docker/deb
DOCKER_NAME_DEB=ubuntu_deb_package

# rpm 
DOCKER_PATH_RPM=$(BASE_PATH)/docker/rpm
DOCKER_NAME_RPM=centos_rpm_package

# network
DOCKER_NETWORK_NAME=network_fax_sender

##########################
## execution 
run_ui_sender:: init
	@$(GO) run $(PWD)/build/main.go run_ui_sender

run_ui:: init
	@$(GO) run $(PWD)/build/main.go run_ui

run:: init
	@$(GO) run $(PWD)/build/main.go run

init::
	@$(GO) run $(PWD)/build/main.go init

package:: 
	@$(GO) run $(PWD)/build/main.go package

test::
	$(GO) test ./test/...

install_rpm:: 
	@sudo dnf install ./bin/print2fax-1.0-1.noarch.rpm -y 
##########################
## deploy 
deploy_linux_daemon:: 
	@export GOOS=linux; export CGO_ENABLED=1 ;$(GO) run $(PWD)/build/main.go deploy_linux_daemon 

deploy_linux_ui:: 
	@export GOOS=linux; export CGO_ENABLED=1 ;$(GO) run $(PWD)/build/main.go deploy_linux_ui  $(ARCH)

deploy_deb:: 
	@$(GO) run $(PWD)/build/main.go deploy_debs $(ARCH)

deploy_rpm::
	@rm -rf ./bin/print2fax-1.0*
	@yum remove print2fax-1.0-1.noarch -y 
	@$(GO) run $(PWD)/build/main.go deploy_rpm $(ARCH)

deploy_windows_ui::
	@$(GO) run $(PWD)/build/main.go deploy_windows_ui $(ARCH)

deploy_windows_installer:: deploy_windows_ui
	@makensis ./build/resources/installer.nsi
	
deploy_windows_daemon::
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc $(GO) build -tags versiontag -mod=vendor -o ./bin/fax_daemon.exe  ./src/server/main.go

deploy_docs::
	@$(GO) install golang.org/x/pkgsite/cmd/pkgsite@latest
	@pkgsite -http=:6060 -cache

deploy_export_sources::
	@rm -rf $(EXPORT_PATH)
	@mkdir -p $(EXPORT_PATH)
	@cp -r build docs src test go.mod go.sum Makefile $(EXPORT_PATH)
	@cd $(EXPORT_PATH); zip -r ../$(EXPORT_DIR_NAME).zip  *
	@rm -rf $(EXPORT_PATH)

##########################################3
## docker 

docker_create_network: 
	@docker network create -d bridge $(DOCKER_NETWORK_NAME)

docker_get_go:
	@if [ -f $(DOCKER_GO_FILE_PATH) ]; then \
		echo "go file exist"; \
	else \
		wget -P $(PWD)/bin/ https://go.dev/dl/go1.21.4.linux-amd64.tar.gz ; \
	fi

docker_copy_source: 
	@rm -rf $(EXPORT_PATH).zip 
	@sudo rm -rf $(DOCKER_SOURCE_PATH)
	@make deploy_export_sources; 
	@mkdir -p $(DOCKER_SOURCE_PATH); 
	@echo $(EXPORT_PATH).zip; 
	@unzip -d $(DOCKER_SOURCE_PATH) $(EXPORT_PATH).zip; 
	@cp vendor -R $(DOCKER_SOURCE_PATH); 
	@rm -rf $(EXPORT_PATH).zip

# deb 
docker_build_deb: docker_get_go docker_copy_source
	@docker build -f $(DOCKER_PATH_DEB)/Dockerfile -t $(DOCKER_NAME_DEB) .

docker_run_deb:
	@docker run -it --rm --network $(DOCKER_NETWORK_NAME) -v $(DOCKER_SOURCE_PATH):/app  --privileged -u root --name $(DOCKER_NAME_DEB)_run $(DOCKER_NAME_DEB)

docker_remove_deb:
	@docker image rm $(DOCKER_NAME_DEB) -f

# rpm 
docker_build_rpm: docker_get_go docker_copy_source
	@docker build -f $(DOCKER_PATH_RPM)/Dockerfile -t $(DOCKER_NAME_RPM) .

docker_run_rpm:
	@docker run -it --rm --network $(DOCKER_NETWORK_NAME) -v $(DOCKER_SOURCE_PATH):/app  --privileged -u root --name $(DOCKER_NAME_RPM)_run $(DOCKER_NAME_RPM)

docker_remove_rpm:
	@docker image rm $(DOCKER_NAME_RPM) -f