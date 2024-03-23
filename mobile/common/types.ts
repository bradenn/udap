// Copyright (c) 2023 Braden Nicholson

// UDAP Core Domain (Alphabetical)

export interface Action {

    name: string;
    entities: string[];
    attribute: string;
    triggerId: string;
    request: string;
    created: string;
    updated: string;
    id: string;
}

export interface Detour {
    time: number
    dose: number
    intensity: number
}

export interface Attribute {
    created: string;
    id: string;
    value: string;
    updated: string;
    request: string;
    requested: string;
    lastUpdated: string;
    entity: string;
    key: string;
    type: string;
    order: number;
}

export interface AttributeLog {
    created: string;
    id: string;
    time: string;
    attribute: string;
    from: string;
    to: string;
    updated: string;
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

export interface Endpoint {
    created: string;
    updated: string;
    id: string;
    name: string;
    type: string;
    notifications: boolean;
    push: string;
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

export interface Module {
    created: string;
    updated: string;
    id: string;
    name: string
    uuid: string
    interval: number
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

export interface Trigger {
    created: string;
    updated: string;
    id: string;
    name: string;
    type: string;
    description: string;
    lastTrigger: string;
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

export interface User {
    created: string;
    updated: string;
    id: string;
    username: string;
    first: string;
    middle: string;
    last: string;
    residency: string;
    classification: string;
    photo: string;
    password: string;
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


export interface Session {
    user: User,
    screensaver: boolean
}


declare interface Persistent {
    created: string;
    updated: string;
    id: string;
}


// Other


export enum TaskType {
    String,
    Number,
    Passcode,
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
    },
    appdata: {
        colors: number[]
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
    targets: Target[]
}

export interface ConfigFile {
    controllers: ControllersEntity[];
    targets: TargetsEntity[];
    preferences: Preferences;
    defaults: Defaults;
}

export interface ControllersEntity {
    name: string;
    address: string;
}

export interface TargetsEntity {
    name: string;
    address: string;
    heartbeat: string;
    description: string;
    required: boolean;
    cacheable: boolean;
}

export interface Preferences {
    background?: (string)[] | null;
}

export interface Defaults {
    background: string;
    theme: string;
    token: string;
}


export interface Target {
    name: string
    address: string
    heartbeat: string
    description: string
    required: boolean
    cacheable: boolean
}

export interface TargetHealth {
    hostname: string
    address: string
    ipv4: string
    ipv6: string
    live: boolean
}

export interface UdapStatus {
    targets: TargetHealth[]
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

// Preferences

export enum PreferenceTypes {
    Controller = "controller",
    Background = "background",
    Screensaver = "screenSaver",
    TouchMode = "touchMode",
    Token = "token",
    Theme = "theme",
    Name = "name",
}

// Udap Models




