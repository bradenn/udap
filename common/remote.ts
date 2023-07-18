// Copyright (c) 2023 Braden Nicholson

import {reactive} from "vue";
import type {
    Attribute,
    Device,
    Endpoint,
    Entity,
    Identifiable,
    Log,
    Macro,
    Metadata,
    Module,
    Network,
    SubRoutine,
    Timing,
    Trigger,
    User,
    Zone
} from "./types";
import {config} from "./config";

export enum Target {
    Metadata = "metadata",
    Module = "module",
    Zone = "zone",
    Entity = "entity",
    User = "user",
    Attribute = "attribute",
    AttributeLog = "attributelog",
    Macro = "macro",
    SubRoutine = "subroutine",
    Trigger = "trigger",
    Device = "device",
    Network = "network",
    Endpoint = "endpoint",
    Timing = "timing",
    Log = "log",
    Close = "close",
}

function connectionString(endpoint: string, token: string): string {
    return `wss://${endpoint}/socket/${token}`
}

export interface Client {
    socket: WebSocket
    url: string

    ready: boolean
    connecting: boolean
    connected: boolean
    closed: boolean

    interval: number
    attempts: number
    nextAttempt: number

    connect(): boolean

    disconnect(): void
}


export interface Remote {
    client: Client,

    attributes: Attribute[],
    devices: Device[],
    endpoints: Endpoint[],
    entities: Entity[],
    logs: Log[],
    macros: Macro[],
    metadata: Metadata,
    modules: Module[],
    networks: Network[],
    subroutines: SubRoutine[],
    triggers: Trigger[],
    timings: Timing[],
    users: User[],
    zones: Zone[],


}

let remote = reactive<Remote>({
    client: {
        socket: {} as WebSocket,

        url: "",

        ready: false,
        connecting: false,
        connected: false,
        closed: false,

        interval: 0,
        attempts: 0,
        nextAttempt: 0,

        connect: connect,
        disconnect: disconnect,
    } as Client,
    attributes: [] as Attribute[],
    devices: [] as Device[],
    endpoints: [] as Endpoint[],
    entities: [] as Entity[],
    logs: [] as Log[],
    macros: [] as Macro[],
    metadata: {} as Metadata,
    modules: [] as Module[],
    networks: [] as Network[],
    subroutines: [] as SubRoutine[],
    timings: [] as Timing[],
    triggers: [] as Trigger[],
    users: [] as User[],
    zones: [] as Zone[],
})

interface Message {
    endpoint: string;
    id: string;
    status: string;
    operation: string;
    body: any;
}

function setup(): boolean {

    let token = localStorage.getItem("token")
    if (token == null) {
        console.log("token not set")
        return false
    }

    let endpoint = localStorage.getItem("endpoint")
    if (endpoint == null) {
        endpoint = config.controllers.find(c => c.name === "Production")?.address || "api.udap.app"
    }

    remote.client.url = connectionString(endpoint, token)

    return true
}

function connect(): boolean {

    if (!setup()) {
        return false;
    }

    if (remote.client.connected || remote.client.connecting) return false
    remote.client.connecting = true

    remote.client.socket = new WebSocket(remote.client.url)
    remote.client.socket.onopen = onOpen
    remote.client.socket.onclose = onClose
    remote.client.socket.onmessage = onMessage
    remote.client.socket.onerror = onError

    return true;
}

function disconnect(): void {
    if (remote.client.ready) {
        remote.client.socket.close(1000, "Disconnecting")
        remote.client.connected = false
        remote.client.connecting = false
    }
}

export {
    connect,
    disconnect
}

function onOpen(_: Event) {
    remote.client.ready = true
    remote.client.connected = true
    remote.client.connecting = false
    remote.client.attempts = 0

}

function beginCountdown() {
    clearInterval(remote.client.interval)
    remote.client.nextAttempt = 2000
    remote.client.attempts++
    //@ts-ignore
    remote.client.interval = setInterval(tick, 100)
}

function tick() {
    if (remote.client.nextAttempt >= 100) {
        remote.client.nextAttempt -= 100
    } else {
        remote.client.nextAttempt = 0
        clearInterval(remote.client.interval)
        remote.client.connect()
    }
}

function onClose(_: CloseEvent) {
    remote.client.ready = false
    remote.client.connected = false
    remote.client.connecting = false
    remote.client.disconnect()
    beginCountdown()
}

function onMessage(event: MessageEvent) {
    remote.client.ready = true
    remote.client.connected = true
    let msg: Message = JSON.parse(event.data) as Message
    if (!msg) {
        console.log("Invalid JSON received")
        return
    }
    if (msg.status !== "success") return;
    let operation: string = msg.operation
    let target: Target = operation as Target
    handleMessage(target, msg.body)
    remote.client.connected = true
}

function onError(err: Event) {
    console.log(err)
    if (remote.client.socket.readyState == WebSocket.CLOSED) {
        remote.client.ready = false
        remote.client.connected = false
        remote.client.connecting = false
    }
    remote.client.ready = true
}


function createOrUpdate(target: any[], data: Identifiable): any[] {
    if (!target) return []
    if (target.find((e: Identifiable) => e.id === data.id)) {
        if (data.deleted) {
            return target.filter((a: Identifiable) => a.id !== data.id)
        } else {
            return target.map((a: Identifiable) => a.id === data.id ? data : a)
        }
    } else {
        if (!data.deleted) {
            target.push(data)
            return target
        }
    }
    return target
}

// const sys = inject("system")

// Handle and route incoming messages to the local cache
function handleMessage(target: Target, data: any) {
    // remote.diagnostics.lastUpdate = new Date().valueOf()
    let dx = 0;
    switch (target) {
        case Target.Close:
            disconnect()
            break
        case Target.Metadata:
            remote.metadata = data as Metadata
            dx = 1
            break
        case Target.Timing:
            if (remote.timings.find((e: Timing) => e.pointer === data.pointer)) {
                remote.timings = remote.timings.map((a: Timing) => a.pointer === data.pointer ? data : a)
                dx = 0
            } else {
                remote.timings.push(data)
                dx = 1
            }
            break
        case Target.Entity:
            remote.entities = createOrUpdate(remote.entities, data)
            break
        case Target.Macro:
            remote.macros = createOrUpdate(remote.macros, data)
            break
        case Target.SubRoutine:
            remote.subroutines = createOrUpdate(remote.subroutines, data)
            break
        case Target.Trigger:
            remote.triggers = createOrUpdate(remote.triggers, data)
            let trigger = data as Trigger
            let last = new Date(trigger.lastTrigger)
            if (new Date().valueOf() - last.valueOf() < 2000) {
                // notifications.show(`Trigger: ${trigger.name}`, trigger.description, 1, 2500)
            }
            break
        case Target.Attribute:
            remote.attributes = createOrUpdate(remote.attributes, data)
            break
        case Target.User:
            remote.users = createOrUpdate(remote.users, data)
            break
        case Target.Device:
            remote.devices = createOrUpdate(remote.devices, data)
            break
        case Target.Network:
            remote.networks = createOrUpdate(remote.networks, data)
            break
        case Target.Endpoint:
            remote.endpoints = createOrUpdate(remote.endpoints, data)
            break
        case Target.Module:
            remote.modules = createOrUpdate(remote.modules, data)
            break
        case Target.Zone:
            remote.zones = createOrUpdate(remote.zones, data)
            break
        case Target.Log:
            remote.entities = createOrUpdate(remote.entities, data)
            break
    }

    // let prev = remote.diagnostics.updates.get(target) || 0
    // remote.diagnostics.updates.set(target, prev + dx);
    // let session = {
    //     target: target,
    //     time: new Date().valueOf(),
    //     operation: "update",
    //     payload: data,
    //     id: (data as Identifiable).id
    // } as RemoteRequest
    // remote.diagnostics.queue.push(session)
    // remote.diagnostics.lastTarget = target
    // if (remote.diagnostics.queue.length >= 10) {
    //     remote.diagnostics.queue = remote.diagnostics.queue.slice(0, remote.diagnostics.queue.length - 2)
    // }

    // remote.diagnostics.maxRSS = memorySizeOf(remote) as number
}


export default remote as Remote



