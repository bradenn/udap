// Copyright (c) 2022 Braden Nicholson

// Configuration Types

import type {Nexus} from "@/views/terminal/nexus";

export interface Trigger {
    created: string;
    updated: string;
    id: string;
    name: string;
    type: string;
    description: string;
    lastTrigger: string;
}


export enum TaskType {
    String,
    Number,
    Object,
    List,
    Radio,
    Icon,
}

export interface TaskOption {
    title: string
    description: string
    value: any
}

export interface Task {
    title: string
    description: string
    type: TaskType
    options?: TaskOption[]
    preview: string
    value: any
}

export interface SubRoutine {
    created: string;
    updated: string;
    id: string;
    revertAfter: number;
    icon: string;
    lastRun: string;
    triggerId: string;
    macros: Macro[];
    description: string;
}

export interface Macro {
    created: string;
    updated: string;
    id: string;
    name: string;
    description: string;
    zone: string;
    type: string;
    value: string;
}

export interface ApiRateLimit {
    remaining: number,
    limit: number,
    reset: number,
}

export interface Landmarks {
    rightEye: {
        xa: number;
        ya: number;
        xb: number;
        yb: number;
    };
    leftEye: {
        xa: number;
        ya: number;
        xb: number;
        yb: number;
    };
    nose: {
        x: number;
        y: number;
    };
}

export interface Prediction {
    name: string;
    top: number;
    right: number;
    bottom: number;
    left: number;
    distance: number;
    landmarks: Landmarks;
}

export interface Status {
    zone: string;
    predictions: Prediction[];
    updated: string;
}

export interface Preferences {
    ui: {
        screensaver: {
            enabled: boolean
            countdown: number
            selection: string
        }
        background: {
            image: string,
            blur: boolean
        }
        accent: string
        theme: string
        mode: string
        blur: number
        brightness: number
        grid: boolean
        watermark: boolean
        night: boolean
        outlines: boolean
    }
}

export interface Calendar {
    description: string;
    summary: string;
    location: string;
    rule: string;
    start: string;
    days: string;
    end: string;
}

export interface Controller {
    name: string
    address: string
    status?: boolean
}

export interface TerminalDiagnostics {
    connected: boolean
    maxRSS: number
    objects: number
    lastUpdate: number
    lastTarget: string
    queue: RemoteRequest[]
    updates: Map<string, number>
}

export interface Remote {
    connecting: boolean,
    connected: boolean,
    metadata: Metadata,
    entities: Entity[],
    subroutines: SubRoutine[],
    macros: Macro[],
    triggers: Trigger[],
    attributes: Attribute[],
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
}

export interface RemoteRequest {
    time: number
    target: string,
    operation: string,
    payload: string,
    id: string
}

export interface ModuleVariable {
    name: string;
    description: string;
    default: string;
}

export interface Module {
    created: string;
    updated: string;
    id: string;
    name: string
    uuid: string
    config: string
    variables: string
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
    startNano: number;
    stop: string;
    stopNano: number;
    delta: number;
    frequency: number;
    complete: boolean;
    depth: number;
    id: string;
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

export interface Compute {
    fanSpeed: number;
    temperature: {
        current: number;
        throttle: number;
        target: number;
        max: number;
    };
    utilization: {
        gpu: number;
        memory: number;
    };
    power: {
        draw: number;
        max: number;
    };
    memory: {
        used: number;
        reserved: number;
        total: number;
    };
    clocks: {
        graphics: {
            current: number;
            max: number;
        };
        streaming: {
            current: number;
            max: number;
        };
        memory: {
            current: number;
            max: number;
        };
        video: {
            current: number;
            max: number;
        };
    };
    processes: {
        name: string;
        pid: string;
        memory: number;
    }[];
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

export interface Log {
    group: string
    level: string
    event: string
    time: string
    message: string
    id: string
}

export interface Identifiable {
    id: string
    deleted: boolean
}

export function memorySizeOf(obj: any) {
    var bytes = 0;

    function sizeOf(obj: any) {
        if (obj !== null && obj !== undefined) {
            switch (typeof obj) {
                case 'number':
                    bytes += 8;
                    break;
                case 'string':
                    bytes += obj.length * 2;
                    break;
                case 'boolean':
                    bytes += 4;
                    break;
                case 'object':
                    var objClass = Object.prototype.toString.call(obj).slice(8, -1);
                    if (objClass === 'Object' || objClass === 'Array') {
                        for (var key in obj) {
                            if (!obj.hasOwnProperty(key)) continue;
                            sizeOf(obj[key]);
                        }
                    } else bytes += obj.toString().length * 2;
                    break;
            }
        }
        return bytes;
    }


    return sizeOf(obj);
}

export function formatByteSize(bytes: number) {
    if (bytes < 1024) return bytes + " bytes";
    else if (bytes < 1048576) return (bytes / 1024).toFixed(3) + " KB";
    else if (bytes < 1073741824) return (bytes / 1048576).toFixed(3) + " MB";
    else return (bytes / 1073741824).toFixed(3) + " GB";
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

export interface Utilization {
    memory: {
        total: number;
        used: number;
    };
    network: {
        hostname: string;
        ipv4: string;
        mac: string;
    };
    cpu: {
        cores: number;
        usage: number[];
    };
    disk: {
        total: number;
        used: number;
    };
    compute: Compute[]
}

export interface Device {
    created: string;
    updated: string;
    state: string;
    id: string;
    networkId: string;
    entityId: string;
    name: string;
    hostname: string;
    utilization: Utilization;
    isQueryable: boolean
    lastSeen: string;
    latency: number;
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
    pinned: boolean;
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

export interface Session {
    user: User,
    screensaver: boolean
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


