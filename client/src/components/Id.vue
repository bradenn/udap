<script lang="ts" setup>

import {inject, onMounted, reactive, watchEffect} from "vue";
import IdHash from "@/components/IdHash.vue"
import type {Attribute, Entity, Preferences, Status} from "@/types";
import type {Remote} from "@/remote";

let state = reactive({
    menu: false,
    reloading: true,
    connected: false,
    zoneEntity: {} as Entity,
    zoneAttribute: {} as Attribute,
    status: {} as Status
})

let preferences: Preferences = inject("preferences") as Preferences
let remote: Remote = inject('remote') as Remote
let system: any = inject('system')

onMounted(() => {
    update()
    state.reloading = false
})

watchEffect(() => {
    state.connected = remote.connected
    update()
    return state.zoneAttribute
})

function update() {
    let entity = remote.entities.find(e => e.name === 'faces')
    if (!entity) return
    state.zoneEntity = entity

    let attr = remote.attributes.find(e => e.key === 'deskFace')
    if (!attr) return
    state.zoneAttribute = attr

    let stat = JSON.parse(attr.value) as Status
    if (!stat) return

    state.status = stat

}


</script>

<template>
    <div class="context context-id">
        <IdHash></IdHash>

    </div>


</template>

<style lang="scss" scoped>

.context-id {
  background-color: rgba(255, 255, 255, 0.3);
}

</style>
