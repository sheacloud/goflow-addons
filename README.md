# goflow-addons
Additional modules for Cloudflares goflow tool

Mostly designed to create a standalone flow ingestion tool rather than leveraging goflows native transport of Kafka + separate consumers for processing/storing the data.
This is not intended to be a better solution at scale, but rather provide a less complex solution for small deployments such as for home networks.

Currently adds:
 - Ability to have multiple targets for the flow messages (i.e. send to Cloudwatch logs and Kafka)
 - Cloudwatch Logs target
   - Supports batching the upload requests by time and size in order to reduce the total number of API calls
 - Flowlog enrichment
   - GeoIP information via MaxMind IP databases
   - Flow direction + client/server differentiation based on configured local CIDRs and src/dst port comparisons
   - Reverse DNS lookups for source and destination IPs

## Extended modules
This package implements a wrapper/extension of goflows Transport interface. The ExtendedWrapper Transport adheres to the original Transport interface, but can be configured with a list of ExtendedTransports to which it publishes each message. This allows you to have multiple targets for the flows, such as sysout and cloudwatch logs, helpful for troubleshooting.

The ExtendedTransport interface is nearly identical to the original Transport interface, except that it's Publish method takes in a list of ExtendedFlowMessages. The ExtendedWrapper takes care of converting the FlowMessages produced by goflow into these ExtendedFlowMessages before sending them to each of the configured ExtendedTransports.

The ExtendedWrapper also can be configured with a list of Enrichers, which are modules that add additional metadata to each flow message, such as GeoIP information or reverse DNS lookup info.
