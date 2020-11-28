# General
 - Move configuration to a config file rather than hard coded
 - Add prometheus metrics for the custom transports and enrichers
 - Look into buffering the message publishing. Currently each thread waits for the message to be fully published before reading more data from the socket. UDP socket buffers are pretty small by default on linux, so this could lead to data loss during large spikes in netflow exports, and would allow the buffer size to be managed by the program. Could be done at the ExtendedWrapper layer

# Transports
 - Add elasticsearch transport
 - Maybe add a Kinesis transport as an alternative to Kafka

## Cloudwatch Logs
 - Improve batching code to ensure that we don't try to upload more than 1MB per batch, and break it into multiple requests if we do


# Enrichers
 - Add an IP reputation enricher (preferably via an offline solution rather than querying an API for each IP)
