// Copyright (c) 2022 Braden Nicholson

// Configuration Types

import type {Nexus} from "@/views/terminal/nexus";


export interface Controller {
    name: string
    address: string
    status?: boolean
}

export interface Remote {
    metadata: Metadata,
    entities: Entity[],
    attributes: Attribute[],
    devices: Device[],
    networks: Network[],
    modules: Module[],
    endpoints: Endpoint[],
    timings: Timing[],
    zones: Zone[],
    nexus: Nexus
}

export interface Module {
    created: string;
    updated: string;
    id: string;
    name: string
    path: string
    type: string
    enabled: boolean
    description: string
    version: string
    author: string
    state: string
}

export interface Timing {
    pointer: string;
    name: string;
    start: string;
    stop: string;
    delta: number;
    frequency: number;
    complete: boolean;
    depth: number;
}

export interface Metadata {
    system: System
}

export interface System {
    name: string;
    version: string;
    environment: string;
    ipv4: string;
    ipv6: string;
    hostname: string;
    mac: string;
    go: string;
    cores: number;
    threads: number[];
}

export interface Defaults {
    background: string
    theme: string
    token: string
    controller: string
    screenSaver: string
    touchMode: string
}

export interface Config {
    controllers: Controller[]
    defaults: Defaults
}

// Remote & Websocket Types

export interface Identifiable {
    id: string
}

// Preferences

export enum PreferenceTypes {
    Controller = "controller",
    Background = "background",
    Screensaver = "screenSaver",
    TouchMode = "touchMode",
    Token = "token",
    Theme = "theme",
}

// Udap Models

export interface Attribute {
    created: string;
    id: string;
    value: string;
    updated: string;
    request: string;
    requested: string;
    entity: string;
    key: string;
    type: string;
    order: number;
}

export interface Device {
    created: string;
    updated: string;
    id: string;
    networkId: string;
    entityId: string;
    name: string;
    hostname: string;
    mac: string;
    ipv4: string;
    ipv6: string;
}

export interface Zone {
    created: string;
    updated: string;
    id: string;
    name: string;
    user: string;
    entities: Entity[];
    deleted: boolean
}

export interface Endpoint {
    created: string;
    updated: string;
    id: string;
    name: string;
    type: string;
    frequency: number;
    connected: boolean;
    key: string;
}

export interface Entity {
    created: string;
    updated: string;
    id: string;
    lastPoll: string;
    name: string;
    alias: string;
    type: string;
    module: string;
    neural: string;
    locked: boolean;
    protocol: string;
    icon: string;
    frequency: number;
    predicted: string;
    state: string;
    config: string;
    position: string;
    live: boolean;
}

export interface Network {
    created: string;
    updated: string;
    id: string;
    name: string;
    dns: string;
    router: string;
    lease: string;
    mask: string;
    range: string;
}

export interface User {
    created: string;
    updated: string;
    id: string;
    username: string;
    first: string;
    middle: string;
    last: string;
    password: string;
}

declare interface Persistent {
    created: string;
    updated: string;
    id: string;
}


