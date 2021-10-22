
CONFIG_FILE=./samconfig.toml

.PHONY: deploy
.SILENT: deploy

deploy:
ifeq ("$(wildcard $(CONFIG_FILE))", "")
	echo "Configuration file does not exist."
	echo "Running with --guided flag"
	sam build
	sam deploy --guided
else
	echo "Configuration file exists."
	echo "Running without the --guided flag"
	sam build
	sam deploy
endif
