// Copyright (c) 2022 Braden Nicholson


import request from "@/services/request";
import type {Endpoint} from "@/types";

function getEndpoint(id: string, path: string): string {
    return `/endpoints/${id}${path}`
}

export default {
    async createEndpoint(body: Endpoint): Promise<void> {
        const url = getEndpoint('', `create`)
        return await request.post(url, body)
    },
}
