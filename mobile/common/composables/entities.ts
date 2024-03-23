// Copyright (c) 2023 Braden Nicholson
import type {UnwrapNestedRefs} from 'vue'
import {onMounted, reactive, watchEffect} from 'vue'
import {Entity} from "../types";
import {Remote} from "../remote";

export interface EntityGroup {
    name: string
    entities: Entity[]
}

interface EntitiesState {
    filter: string,
    groupBy: string,
    entities: Entity[],
    loaded: boolean,
    lastUpdated: number,
    groups: EntityGroup[],
    labels: string[]
    toggleGroup: () => void
    toggleFilter: () => void
}

const groupTypes = ["module", "type", "none"]
const filterTypes = ["alphabetical", "updated", "created"]

export default function useEntities(remote: Remote): EntitiesState {
    const state: UnwrapNestedRefs<EntitiesState> = reactive<EntitiesState>({
        filter: "alphabetical" as string,
        groupBy: "module" as string,
        entities: [],
        groups: [] as EntityGroup[],
        labels: [] as string[],
        loaded: false,
        lastUpdated: Date.now().valueOf(),
        toggleGroup: toggleGroupBy,
        toggleFilter: toggleFilter
    })


    function sortByModule() {
        let groups: EntityGroup[] = []
        for (let i = 0; i < state.entities.length; i++) {
            let e = state.entities[i]

            let fd = groups.find(g => g.name == e.module)
            if (!fd) {
                groups.push({
                    name: e.module,
                    entities: [e]
                })
            } else {
                groups[groups.indexOf(fd)].entities.push(e)
            }
        }
        state.groups = groups

    }

    function sortByNone() {
        let groups: EntityGroup[] = []

        groups.push({
            name: "All",
            entities: state.entities
        })
        state.groups = groups
        state.loaded = true
    }

    function sortByType() {
        let groups: EntityGroup[] = []
        for (let i = 0; i < state.entities.length; i++) {
            let e = state.entities[i]

            let fd = groups.find(g => g.name == e.type)
            if (!fd) {
                groups.push({
                    name: e.type,
                    entities: [e]
                })
            } else {
                groups[groups.indexOf(fd)].entities.push(e)
            }
        }
        state.groups = groups
        state.loaded = true
    }

    function applyGroup() {
        switch (state.groupBy) {
            case "module":
                sortByModule()
                break;
            case "type":
                sortByType()
                break;
            case "none":
                sortByNone()
                break;
            default:
                break
        }
    }

    function toggleGroupBy() {
        let id = groupTypes.indexOf(state.groupBy)
        if (id < 0) {
            return
        }
        let nid = (id + 1) % groupTypes.length
        state.groupBy = groupTypes[nid]
        updateLocal()
    }

    function sortAlphabetical() {
        // state.entities = state.entities.sort((a: Entity, b: Entity) => (a.alias ? a.alias : a.name).localeCompare((b.alias ? b.alias : b.name)))
        state.groups.forEach(group => {
            group.entities = group.entities.sort((a, b) => (a.alias ? a.alias : a.name || "").localeCompare((b.alias ? b.alias : b.name || "")))
        })
    }

    function sortUpdated() {
        state.groups.forEach(group => {
            group.entities = group.entities.sort((a, b) => new Date(a.updated).valueOf() - new Date(b.updated).valueOf())
        })
    }

    function sortCreated() {
        state.groups.forEach(group => {
            group.entities = group.entities.sort((a, b) => new Date(a.created).valueOf() - new Date(b.created).valueOf())
        })
    }

    function applySort() {
        switch (state.filter) {
            case "alphabetical":
                sortAlphabetical()
                break;
            case "updated":
                sortUpdated()
                break;
            case "created":
                sortCreated()
                break;
            default:
                break
        }
    }

    function toggleFilter() {
        let id = filterTypes.indexOf(state.filter)
        if (id < 0) {
            return
        }
        let nid = (id + 1) % filterTypes.length
        state.filter = filterTypes[nid]
        updateLocal()
    }

    onMounted(() => {
        updateLocal()
    })

    watchEffect(() => {
        updateLocal()
        return remote.entities
    })

    function updateLocal() {
        state.entities = remote.entities
        state.entities = state.entities.sort((a: Entity, b: Entity) => (new Date(b.updated).valueOf() - new Date(a.updated).valueOf()))
        state.lastUpdated = Date.now().valueOf()
        applyGroup()
        applySort()
        state.loaded = true
    }


    return state
}

