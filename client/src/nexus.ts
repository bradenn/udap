// Copyright (c) 2022 Braden Nicholson

import {Attr, Preference} from "./preferences"


enum ConnectionState {
    Connecting = "connecting",
    Connected = "connected",
    Disconnected = "disconnected",
}

export enum Target {
    Metadata = "metadata",
    Entity = "entity",
    Attribute = "attribute",
    Device = "device",
    Network = "network",
    Endpoint = "endpoint",
    Timing = "timing",
}


class Nexus {
    ws: WebSocket
    state: ConnectionState
    handler: (target: Target, data: any) => void

    constructor() {
        this.state = ConnectionState.Disconnected

    }

    private static connectionString() {
        let nexus = new Preference(Attr.Controller).get()
        let token = new Preference(Attr.Token).get()
        if (token === "unset") {
            return
        }
        return `ws://${nexus}/socket/${token}`
    }

    public connect(fn: (target: Target, data: any) => void) {
        this.ws = new WebSocket(Nexus.connectionString())
        this.ws.onopen = this.onOpen
        this.ws.onmessage = (event: WebSocket.MessageEvent) => {
            if (typeof event.data === "string") {
                let data: any = JSON.parse(event.data)
                if (!data) {
                    console.log("Invalid JSON received")
                    return
                }
                let operation: string = data.operation
                let target: Target = operation as Target
                fn(target, data.body)
            }
        }
        this.ws.onclose = this.onClose
    }

    enroll() {

    }

    private onOpen(event: WebSocket.Event) {
        this.state = ConnectionState.Connected
    }

    private onClose(event: WebSocket.CloseEvent) {
        this.state = ConnectionState.Disconnected
    }

}

export {
    Nexus
}