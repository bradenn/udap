// Copyright (c) 2022 Braden Nicholson


import type {Macro} from "udap-ui/types";
import request from "@/services/request";


export default {
    async createMacro(macro: Macro): Promise<void> {
        const url = `/macros/create`
        return await request.post(url, macro)
    },
    async runMacro(id: string): Promise<void> {
        const url = `/macros/${id}/run`
        return await request.post(url)
    },
    async deleteMacro(id: string): Promise<void> {
        const url = `/macros/${id}/delete`
        return await request.post(url)
    },
    async updateMacro(macro: Macro): Promise<void> {
        const url = `/macros/${macro.id}/update`
        return await request.post(url, macro)
    },
}
