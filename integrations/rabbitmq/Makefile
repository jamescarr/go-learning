DELAYED_EXCH_FILE := rabbitmq_delayed_message_exchange-0.0.1.ez
DELAYED_EXCH_DOWNLOAD_URL := https://bintray.com/rabbitmq/community-plugins/download_file?file_path=${DELAYED_EXCH_FILE}

get_delayed_exch_plugin:
	wget ${DELAYED_EXCH_DOWNLOAD_URL} -O plugins/${DELAYED_EXCH_FILE}

build:
	docker-compose build

setup: get_delayed_exch_plugin build

up:
	docker-compose up -d

down:
	docker-compose down -d
