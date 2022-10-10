// Copyright (c) 2022 Braden Nicholson


import type {Trigger} from "@/types";
import request from "@/services/request";


export default {
    async createTrigger(trigger: Trigger): Promise<void> {
        const url = `/triggers/create`
        return await request.post(url, trigger)
    }
}
