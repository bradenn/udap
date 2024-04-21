// Copyright (c) 2023 Braden Nicholson


export interface Sensor {
    ax: number[],
    ay: number[],
    az: number[],
    fast: boolean,

    gx: number[],
    gy: number[],
    gz: number[],
    time: number[],

    core?: number[],
    ambient: number[],
    rssi?: number,
    temperature?: number,
}

export interface ProcessedSensorData {
    core: number,
    ambient: number,
    magnitude: number[]
    time: number[]
}

