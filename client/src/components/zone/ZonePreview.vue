<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import Plot from "@/components/plot/Plot.vue";
import {inject} from "vue";
import type {Zone} from "@/types";
import axios from "axios";
import type {Remote} from "@/remote";

interface ZoneProps {
    zone: Zone
}

let props = defineProps<ZoneProps>()

let remote = inject('remote') as Remote

function restoreZone(id: string) {
    axios.post(`http://localhost:3020/zones/${id}/restore`).then(res => {
    }).catch(err => console.log(err))
    // remote.nexus.requestId("zone", "restore", "", id)
}

function deleteZone(id: string) {
    axios.post(`http://localhost:3020/zones/${id}/delete`).then(res => {
    }).catch(err => console.log(err))
    // remote.nexus.requestId("zone", "delete", "", id)
}

</script>

<template>

    <Plot :cols="2" :rows="3">
        <div class="subplot subplot-inline">
            <div class="d-flex label-sm label-o4">{{ props.zone.name }}</div>
        </div>
        <div class="d-flex label-sm label-o4 subplot justify-content-end px-1 subplot-inline"
             @click="props.zone.deleted?restoreZone(props.zone.id):deleteZone(props.zone.id)">
            <div>
                {{ props.zone.deleted ? 'Restore' : 'Delete' }}
            </div>

        </div>

        <div v-for="entity in props.zone.entities" class="subplot">
            <div v-if="entity" class="d-flex justify-content-start align-items-center flex-row  w-100">
                <div class="label-w500 label-o3 label-c1">{{ entity.icon }}&nbsp;</div>
                <div class="label-w500 label-c1">{{ entity.name }}</div>
                <div class="flex-grow-1"></div>
            </div>
        </div>

        <div>
        </div>
    </Plot>
</template>

<style lang="scss" scoped>


.app-container {
  display: flex;
  justify-content: center;
  flex-direction: column;
  align-items: center;

}

.app-icon::before {

}

.app-icon {

  width: 2.5rem;
  height: 2.5rem;
  display: flex;
  align-self: center;
  align-items: center;
  justify-content: center;
  aspect-ratio: 1/1 !important;
  border-radius: 0.6rem;
  font-size: 0.9rem;
  color: rgba(200, 200, 200, 0.6);
}

.app-container:active {


}

.app-name {
  font-size: 0.6rem;
  font-family: "SF Pro Display", sans-serif;
  font-weight: 400;
  text-shadow: 0 0 1px rgba(0, 0, 0, 0.2);
  color: rgba(255, 255, 255, 0.55)
}
</style>