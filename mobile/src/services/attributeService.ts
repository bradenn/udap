// Copyright (c) 2022 Braden Nicholson


import type {Attribute} from "@/types";
import request from "@/services/request";


export default {
    async request(attribute: Attribute): Promise<void> {
        const url = `/entities/${attribute.entity}/attributes/${attribute.key}/request`
        return await request.post(url, attribute.request)
    }
}
