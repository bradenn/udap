// Copyright (c) 2023 Braden Nicholson

import {reactive} from "vue";
import {Nexus, Target} from "@/nexus";
import type {
    Attribute,
    AttributeLog,
    Device,
    Endpoint,
    Entity,
    Identifiable,
    Log,
    Macro,
    Metadata,
    Module,
    Network,
    RemoteRequest,
    SubRoutine,
    TerminalDiagnostics,
    Timing,
    Trigger,
    User,
    Zone
} from "./types";

import {PreferenceTypes} from "udap-ui/types";
import {Preference} from "@/preferences";


function connectionString(): string {
    let nexus = new Preference(PreferenceTypes.Controller).get()
    let token = new Preference(PreferenceTypes.Token).get()
    if (token === "unset") {
        return ""
    }
    return `wss://api.udap.app/socket/${token}`
}

export interface Remote {
    connecting: boolean,
    attempts: number,
    connected: boolean,
    interval: number,
    nextAttempt: number,
    metadata: Metadata,
    entities: Entity[],
    subroutines: SubRoutine[],
    macros: Macro[],
    triggers: Trigger[],
    attributes: Attribute[],
    attributeLogs: AttributeLog[],
    users: User[],
    devices: Device[],
    networks: Network[],
    modules: Module[],
    endpoints: Endpoint[],
    timings: Timing[],
    zones: Zone[],
    logs: Log[],
    nexus: Nexus,
    size: string,
    diagnostics: TerminalDiagnostics

    connect(): void,

    disconnect(): void,
}

let remote = reactive<Remote>({
    connecting: false,
    interval: 0,
    attempts: 0,
    nextAttempt: 0,
    connected: false,
    metadata: {} as Metadata,
    entities: [] as Entity[],
    attributes: [] as Attribute[],
    attributeLogs: [] as AttributeLog[],
    devices: [] as Device[],
    networks: [] as Network[],
    endpoints: [] as Endpoint[],
    users: [] as User[],
    timings: [] as Timing[],
    modules: [] as Module[],
    zones: [] as Zone[],
    subroutines: [] as SubRoutine[],
    macros: [] as Macro[],
    triggers: [] as Trigger[],
    logs: [] as Log[],
    nexus: {} as Nexus,
    size: "" as string,
    diagnostics: {
        queue: [] as RemoteRequest[],
        updates: new Map<string, number>(),
        connected: false,
        maxRSS: 0,
        lastTarget: "",
        lastUpdate: 0,
        objects: 0
    } as TerminalDiagnostics,
    connect: connect,
    disconnect: disconnect,
})

const state = reactive({
    ready: false,
    ws: {} as WebSocket
})

interface Message {
    endpoint: string;
    id: string;
    status: string;
    operation: string;
    body: any;
}


function connect(): void {
    if (remote.connected || remote.connecting) return
    // console.log("connecing")
    remote.connecting = true

    state.ws = new WebSocket(connectionString())
    state.ws.onopen = onOpen
    state.ws.onclose = onClose
    state.ws.onmessage = onMessage
    state.ws.onerror = onError
}

function disconnect(): void {
    if (state.ready) {
        state.ws.close(1000, "Disconnecting")
        remote.connected = false
        remote.connecting = false
    }
}

export {
    connect,
    disconnect
}

function onOpen(_: Event) {
    state.ready = true
    remote.connected = true
    remote.connecting = false
    remote.attempts = 0

}

function beginCountdown() {
    clearInterval(remote.interval)
    remote.nextAttempt = 2000
    remote.attempts++
    //@ts-ignore
    remote.interval = setInterval(tick, 33)
}

function tick() {
    if (remote.nextAttempt >= 33) {
        remote.nextAttempt -= 33
    } else {
        remote.nextAttempt = 0
        clearInterval(remote.interval)
        remote.connect()
    }
}

function onClose(_: CloseEvent) {
    state.ready = false
    remote.connected = false
    remote.connecting = false
    remote.disconnect()
    beginCountdown()
}

function onMessage(event: MessageEvent) {
    state.ready = true
    remote.connected = true
    let msg: Message = JSON.parse(event.data) as Message
    if (!msg) {
        console.log("Invalid JSON received")
        return
    }
    if (msg.status !== "success") return;
    let operation: string = msg.operation
    let target: Target = operation as Target
    handleMessage(target, msg.body)
    remote.connected = true
}

function onError(err: Event) {
    console.log(err)
    if (state.ws.readyState == WebSocket.CLOSED) {
        state.ready = false
        remote.connected = false
        remote.connecting = false
    }
    state.ready = true
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
    remote.diagnostics.lastUpdate = new Date().valueOf()
    let dx = 0;
    switch (target) {
        case Target.Close:
            setTimeout(connect, 5000)
            return

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
        case Target.AttributeLog:
            remote.attributeLogs = createOrUpdate(remote.attributeLogs, data)
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

    let prev = remote.diagnostics.updates.get(target) || 0
    remote.diagnostics.updates.set(target, prev + dx);
    let session = {
        target: target,
        time: new Date().valueOf(),
        operation: "update",
        payload: data,
        id: (data as Identifiable).id
    } as RemoteRequest
    remote.diagnostics.queue.push(session)
    remote.diagnostics.lastTarget = target
    if (remote.diagnostics.queue.length >= 10) {
        remote.diagnostics.queue = remote.diagnostics.queue.slice(0, remote.diagnostics.queue.length - 2)
    }

    // remote.diagnostics.maxRSS = memorySizeOf(remote) as number
}


export default remote as Remote



