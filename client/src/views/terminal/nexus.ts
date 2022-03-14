// Copyright (c) 2022 Braden Nicholson

import {Preference} from "@/preferences"

import {PreferenceTypes} from "@/types";

export enum NexusState {
    Connecting,
    Connected,
    Disconnected,
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

function parseJwt(token: string): string {
    let base64Url = token.split('.')[1];
    let base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    let jsonPayload = decodeURIComponent(atob(base64).split('').map(function (c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));

    return JSON.parse(jsonPayload);
}

function connectionString(): string {
    let nexus = new Preference(PreferenceTypes.Controller).get()
    let token = new Preference(PreferenceTypes.Token).get()
    if (token === "unset") {
        return ""
    }
    return `ws://${nexus}/socket/${token}`
}

interface NexusRequest {
    target: string,
    operation: string,
    payload: string,
    id: string
}

export class Nexus {
    ws: WebSocket
    state: NexusState

    constructor(fn: (target: Target, data: any) => void) {
        this.state = NexusState.Connecting
        this.ws = new WebSocket(connectionString())
        this.ws.onopen = this.onOpen
        this.ws.onmessage = (event: MessageEvent) => {
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

    public requestId(target: string, operation: string, data: any, id: string) {
        const request: NexusRequest = {
            target: target,
            operation: operation,
            payload: JSON.stringify(data),
            id: id
        }
        this.request(request)
    }

    private request(request: NexusRequest) {
        this.ws.send(JSON.stringify(request));
    }

    private onOpen(_: Event): void {
        this.state = NexusState.Connected
    }

    private onClose(_: CloseEvent): void {
        this.state = NexusState.Disconnected
    }

}
