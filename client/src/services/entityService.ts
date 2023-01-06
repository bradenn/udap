// Copyright (c) 2022 Braden Nicholson


import request from "@/services/request";

function getEndpoint(id: string, path: string): string {
    return `/entities/${id}${path}`
}

export default {
    async setAlias(id: string, alias: string): Promise<void> {
        const url = getEndpoint(id, "/alias")
        return await request.post(url, alias)
    },
    async setIcon(id: string, icon: string): Promise<void> {
        const url = getEndpoint(id, "/icon")
        return await request.post(url, icon)
    },
    async toggleNeural(id: string, neural: boolean): Promise<void> {
        const url = getEndpoint(id, "/neural")
        return await request.post(url, neural)
    },
    async setPosition(id: string, x: number, y: number): Promise<void> {
        const url = getEndpoint(id, "/position")

        interface pos {
            x: number
            y: number
        }

        let pso = {
            x: x,
            y: y
        } as pos

        return await request.post(url, pso)
    },
    async setEnergyUsage(id: string, energy: number): Promise<void> {
        const url = getEndpoint(id, "/energy")
        return await request.post(url, energy)
    }
}
