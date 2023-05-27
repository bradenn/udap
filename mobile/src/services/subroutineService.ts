// Copyright (c) 2022 Braden Nicholson


import type {SubRoutine} from "@/types";
import request from "@/services/request";


export default {
    async triggerManual(id: string): Promise<void> {
        const url = `/subroutines/${id}/run`
        return await request.post(url)
    },
    async addMacro(id: string, macroId: string): Promise<void> {
        const url = `/subroutines/${id}/macros/${macroId}/add`
        return await request.post(url)
    },
    async removeMacro(id: string, macroId: string): Promise<void> {
        const url = `/subroutines/${id}/macros/${macroId}/remove`
        return await request.post(url)
    },
    async createSubroutine(subroutine: SubRoutine): Promise<void> {
        const url = `/subroutines/create`
        return await request.post(url, subroutine)
    },
    async deleteSubroutine(id: string): Promise<void> {
        const url = `/subroutines/${id}/delete`
        return await request.post(url)
    },
    async updateSubroutine(param: SubRoutine) {
        const url = `/subroutines/${param.id}/update`
        return await request.post(url, param)
    }
}
