// Copyright (c) 2022 Braden Nicholson


import request from "./request";

function getEndpoint(id: string, path: string): string {
    return `/modules/${id}${path}`
}

export declare interface TraceRequest {
    to: number;
    from: number;
    window: number;
    mode: string;
    labels: string[];
}

export declare interface Trace {
    name: string;
    time: number[];
    data: number[];
    labels: any
}

export declare interface TraceResults {
    traces: Trace[];
}

export default {
    async trace(trace: TraceRequest): Promise<TraceResults> {

        return await request.post("/trace", trace).then(res => res.data)
    },

}
