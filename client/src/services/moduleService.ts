// Copyright (c) 2022 Braden Nicholson


import request from "@/services/request";

function getEndpoint(id: string, path: string): string {
    return `/modules/${id}${path}`
}

export default {
    async setEnabled(id: string, enabled: boolean): Promise<void> {
        const url = getEndpoint(id, `/${enabled ? "enable" : "disable"}`)

        return await request.post(url)
    },
    async reload(id: string): Promise<void> {
        const url = getEndpoint(id, "/reload")
        return await request.post(url)
    },
}
