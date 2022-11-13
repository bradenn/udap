// Copyright (c) 2022 Braden Nicholson


import request from "@/services/request";
import type {Zone} from "@/types";

function getEndpoint(id: string, path: string): string {
    return `/zones/${id}${path}`
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
    },
    async createZone(zone: Zone): Promise<void> {
        return await request.post('/zones/create', zone)
    },
    async updateZone(zone: Zone): Promise<void> {
        return await request.post('/zones/update', zone)
    },
}
