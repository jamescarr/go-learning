FROM learningstaircase/rabbitmq-cluster

COPY plugins /provided-plugins

RUN cp /provided-plugins/* /usr/lib/rabbitmq/lib/rabbitmq_server-3.6.6/plugins

CMD rabbitmq-plugins enable --offline rabbitmq_delayed_message_exchange rabbitmq_consistent_hash_exchange rabbitmq_sharding && rabbitmq-cluster
