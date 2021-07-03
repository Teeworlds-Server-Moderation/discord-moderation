# discord-moderation


Discord moderation is yet another module in this whole microservice framework that allows for a cleaner interface and easier processing of events.

Compared to the monolithic v1 of the Teeworlds econ discord moderation bot this bot solely processes events an doe snot handle any log parsing as well as event creation.
This decoupling of event creation and event processing allows to have a cleaner code base.
It might be possible to have a lot of independent event processors that might create a bigger code base, but at the same time allow for an easily maintainable architecture.



We do have RabbitMQ as a central instance of event storage that may receive events from any kind of monitor.
In our case we have a zCatch-monitor that parses the server logs at runtime and creates events and publishes those events at the RabbitMQ broker.

Those events are published at event specific exchanges.
Any new queue can be bound to an exchange in order to receive that event as well allowing for complex event processing topologies.
In our case we do have an exchange for every event type as well as a broadcasting exchange that can be used to fetch server states as well as broadcast command execution requests to all servers.


The discord moderation bot subscribes to an individual queue that is bound to all available exchanges in order to receive all events that we want to process.
Each received event may pass through any number of processor functions that may d whatever they want with that event.
They may:
    - log any event
    - react to a specific event (detect VPNs of joining players, abort kickvotes)

