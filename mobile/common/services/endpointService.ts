// Copyright (c) 2022 Braden Nicholson


import request from "./request";


function getEndpoint(id: string, path: string): string {
    return `/endpoints/${id}${path}`
}

export default {
    async registerPush(id: string, body: string): Promise<void> {
        const url = getEndpoint(id, `/push`)
        return await request.post(url, body)
    },
}
