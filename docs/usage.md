# Usage

*Docker Flow Swarm Listener* exposes a API to query series and to send notifications.

## Get Services

The *Get Services* endpoint is used to query all running services with the `DF_NOTIFY_LABEL` label. A `GET` request to **[SWARM_IP]:[SWARM_PORT]/v1/docker-flow-swarm-listener/get-services** returns a json representation of these services.

## Notify Services

*DFSL* normally sends out notifcations when a service is created, updated, or removed. The *Notify Services* endpoint will force *DFSL* to send out notifications for all running services with the `DF_NOTIFY_LABEL` label. A `GET` request to **[SWARM_IP]:[SWARM_PORT]/v1/docker-flow-swarm-listener/notify-services** sends out the notifications.
