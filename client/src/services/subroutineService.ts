// Copyright (c) 2022 Braden Nicholson


import type {SubRoutine} from "@/types";
import request from "@/services/request";


export default {
    async triggerManual(subroutine: SubRoutine): Promise<void> {
        const url = `/subroutines/${subroutine.id}/run`
        return await request.post(url)
    },
    async createSubroutine(subroutine: SubRoutine): Promise<void> {
        const url = `/subroutines/create`
        return await request.post(url, subroutine)
    },
    async deleteSubroutine(id: string): Promise<void> {
        const url = `/subroutines/${id}/delete`
        return await request.post(url)
    }
}
