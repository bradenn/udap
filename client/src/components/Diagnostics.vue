<script lang="ts" setup>

import {inject, reactive} from "vue";
import Loader from "@/components/Loader.vue";

let remote: any = inject("remote")
let ui: any = inject("ui")

let state = reactive({
  reloading: false
})

function reload() {
  state.reloading = true
  document.location.reload()
}

</script>

<template>
  <div v-if="ui.context" class="context d-flex flex-column pt-4 gap-4" @click="ui.context = false">
    <div class="d-flex justify-content-between align-items-center" @click.stop>
      <div>
        <h2 class="label-w600">Diagnostics</h2>
      </div>
      <div class="app-icon element" @click="reload">
        <div v-if="state.reloading">
          <Loader size="sm"></Loader>
        </div>
        <div v-else>
          <i class="fa-solid fa-arrow-rotate-right"></i>
        </div>

      </div>
    </div>
    <div class="d-flex flex-row justify-content-around gap-3 w-100">
      <div>
        <h3>Entities</h3>
        <div v-for="ent in remote.entities" v-if="remote.entities.length > 0"
             class="label-c1 d-flex flex-column gap-0">
          <span class="label-w600 label-o5">{{ ent.name }} - {{ ent.module }}</span>

          <div v-for="attr in remote.attributes.filter((att: any) => att.entity === ent.id)"
               :key="attr.id"
               class="d-flex label-o4 label-c4">
            <div style="max-width: 11rem; text-overflow: ellipsis; overflow: clip;">{{ attr.key }} {{ attr.value }}
            </div>
          </div>
        </div>
        <div v-else class="label-o4 label-c2">
          No entities
        </div>
      </div>
      <div>
        <h3>Network</h3>
        <div v-for="net in remote.networks" v-if="remote.networks.length > 0"
             class="label-o4 label-c2 d-flex flex-column gap-0">
          <div class="element p-2">
            <div class="">{{ net.name }}</div>
            <div class="">{{ net.dns }}</div>
            <div class="">{{ net.range }}</div>
          </div>
        </div>
        <div v-else class="label-o5 label-c2">
          No networks
        </div>
      </div>
      <div>
        <h3>Devices</h3>
        <div v-for="net in remote.devices" v-if="remote.devices.length > 0" class="label-c1 d-flex flex-column gap-0">
          <span class="label-w600 label-o5">{{ net.name }} - {{ net.ipv4 }}</span>
        </div>
        <div v-else class="label-o4 label-c2">
          No devices
        </div>
      </div>
      <div>
        <h3>Endpoints</h3>
        <div v-for="net in remote.endpoints" v-if="remote.endpoints.length > 0"
             class="label-c1 d-flex flex-column gap-0">
          <span class="label-w600 label-o5">{{ net.name }} - {{ net.connected }}</span>
        </div>
        <div v-else class="label-o4 label-c2">
          No endpoints
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss">

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

.spin {
  animation: rotate 100ms linear !important;
}

@keyframes rotate {
  0% {
    transform: rotateY(0deg);
  }
  100% {
    transform: rotateY(359deg);
  }
}
</style>
