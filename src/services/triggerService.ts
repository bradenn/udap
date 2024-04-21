// Copyright (c) 2023 Braden Nicholson


import request from "@/services/request";


export default {
    async trigger(id: string): Promise<void> {
        const url = `/triggers/${id}/trigger`
        return await request.post(url)
    }
}
