# Realtime Broker Scaling Note

Current implementation uses in-memory fan-out suitable for a single API instance.

To scale horizontally:

1. Replace direct `Broker.Publish` fan-out with a shared message bus (Redis Pub/Sub or NATS).
2. Keep `Broker` as the per-instance SSE/WebSocket fan-out sink.
3. Add a subscriber worker in each instance to receive cross-instance order events.

This preserves handler contracts while allowing multi-instance kitchen streams.
