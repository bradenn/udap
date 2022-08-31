// Copyright (c) 2022 Braden Nicholson


import request from "@/services/request";

function getEndpoint(id: string, path: string): string {
    return `http://10.0.1.2:3020/zones/${id}${path}`
}

export default {
    async setDeleted(id: string, deleted: boolean): Promise<void> {
        const url = getEndpoint(id, `/${deleted ? 'restore' : 'delete'}`)
        return await request.post(url)
    },
    async setPinned(id: string, pinned: boolean): Promise<void> {
        const url = getEndpoint(id, `/${pinned ? 'unpin' : 'pin'}`)
        return await request.post(url)
    },
    async addEntity(id: string, entity: number): Promise<void> {
        const url = getEndpoint(id, `/entity/${entity}/add`)
        return await request.post(url)
    },
    async removeEntity(id: string, entity: number): Promise<void> {
        const url = getEndpoint(id, `/entity/${entity}/remove`)
        return await request.post(url)
    }
}
