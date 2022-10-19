// Copyright (c) 2022 Braden Nicholson

import type {Remote} from "@/types";
import {Target} from "@/views/terminal/nexus";

interface Message {
    endpoint: string;
    id: string;
    status: string;
    operation: string;
    body: any;
}

let ws: WebSocket
let remote: Remote
let ready: boolean

// provide("remote", nil)

function connect(url: string): boolean {
    ws = new WebSocket(url)
    ws.onopen = onOpen
    ws.onclose = onClose
    ws.onmessage = onMessage
    ws.onerror = onError
}

function onOpen(event: Event) {
    ready = true
}

function onClose(event: CloseEvent) {
    ready = false
}

function onMessage(event: MessageEvent) {
    let msg: Message = JSON.parse(event.data) as Message
    if (msg.status !== "success") return;
    classifyObject(msg.operation, msg.body)
}

function onError(event: Event) {
    ready = false
}

function classifyObject(target: string, data: any) {
    switch (target) {
        case Target.Close:
            break
        case Target.Metadata:
            break
        case Target.Timing:
            break
        case Target.Entity:
            break
        case Target.Macro:
            break
        case Target.SubRoutine:
            break
        case Target.Trigger:
            break
        case Target.Attribute:
            break
        case Target.User:
            break
        case Target.Device:
            break
        case Target.Network:

            break
        case Target.Endpoint:

            break
        case Target.Module:

            break
        case Target.Zone:

            break
        case Target.Log:

            break
    }
}