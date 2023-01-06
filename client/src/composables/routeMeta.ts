// Copyright (c) 2023 Braden Nicholson

import {useRouter} from "vue-router";

interface RouteMeta {
    icon: string
    order: number
}

export function useRouteMeta(): RouteMeta {

    const router = useRouter()

    let matched = router.currentRoute.value.matched

    if (matched.length >= 3) {
        let local = matched[2].meta
        return {
            order: local.order as number | 0,
            icon: local.icon as string | "?"
        }
    }
    return {} as RouteMeta
}



