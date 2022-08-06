<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>


import Plot from "@/components/plot/Plot.vue";
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Attribute, Entity, Remote, Zone} from "@/types";
import Light from "@/components/widgets/Light.vue";
import Widget from "@/components/widgets/Widget.vue";
import Select from "@/components/plot/Select.vue";
import Subplot from "@/components/plot/Subplot.vue";

let remote = inject("remote") as Remote
let preferences = inject('preferences')

let zones = [
  {
    name: "All",
    key: "all"
  },
  {
    name: "Bedroom",
    key: "bedroom"
  },
  {
    name: "Kitchen",
    key: "kitchen"
  },
  {
    name: "Lor",
    key: "lor"
  }
]


let state = reactive({
  zone: {} as Zone,
  zones: [] as Zone[],
  lights: {} as Entity[],
  targets: ["8c1494c3-6515-490b-8f23-1c03b87bde27", "9a3347a7-7e19-4be5-976c-22384c59142a", "c74d427b-5046-4aeb-8195-2efd05d794f8"] as string[],
  loading: true,
  colorMenu: false,
})

onMounted(() => {
  state.loading = true
  handleUpdates(remote)
  switchZone("all")
})


watchEffect(() => handleUpdates(remote))

function sortZones(a: Zone, b: Zone): number {
  if (a.pinned && !b.pinned) {
    return -1;
  } else if (!a.pinned && b.pinned) {
    return 1;
  } else if (a.pinned && b.pinned && a.name > b.name) {
    return 2;
  } else if (a.pinned && b.pinned && a.name < b.name) {
    return -2;
  }
  return 0
}

function handleUpdates(remote: Remote) {
  state.lights = remote.entities.filter((entity: Entity) => state.targets.includes(entity.id))
  state.zones = remote.zones.filter((zone: Zone) => !zone.deleted).sort(sortZones)
  state.loading = false

  if (!state.zone.name) {
    switchZone("all")
  }

  return remote.entities
}

function switchZone(name: string) {
  const zone = remote.zones.find(z => z.name === name)
  if (!zone) return
  state.zone = zone
  state.targets = zone.entities.map(e => e.id)
  handleUpdates(remote)
}

function setAttributes(key: string, value: string) {
  remote.attributes.filter((a: Attribute) => a.key == key && state.targets.includes(a.entity)).forEach(v => {
    remote.nexus.requestAttribute(v.entity, key, value)

  })
}


</script>

<template>
  <div v-if="!state.loading" class="d-flex flex-column gap-1" style="width: 11rem;">
    <div>
      <Select :selected="`${state.zone.name?.charAt(0).toUpperCase()}${state.zone.name?.substring(1)}`">
        <div v-for="zone in state.zones" :class="state.zone.name !== zone.name?'subplot-inline':''" class="subplot"
             @click="() => switchZone(zone.name)">
          <div class="d-flex align-content-center align-items-center justify-content-center">
            <span class="label-c4 label-o2 lh-1" style="width: 0.75rem; margin-top: 4px;"><span
                v-if="zone.pinned">􀎧</span></span>
            <div class="text-capitalize lh-1">
              {{ zone.name }}
            </div>
          </div>
          <div class="label-c3 label-o3">
            {{ zone.entities.length }} 􀛮
          </div>

        </div>
      </Select>
    </div>
    <Widget v-if="!state.loading" :cols="4" :rows="5" class="d-flex flex-column" size="sm">
      <Light v-for="light in state.lights.slice(0, 5)"
             :key="light.id"
             :entity="light"></Light>
    </Widget>
    <Widget v-if="!state.loading" :cols="4" :rows="3" class="d-flex flex-column" size="sm">

      <Plot :cols="4" :rows="1" style="width: 100%">
        <Subplot :active="false" :fn="() => setAttributes('on', 'false')" name="OFF"></Subplot>
        <Subplot :active="false" :fn="() => setAttributes('on', 'true')" name="ON"></Subplot>
        <Subplot :active="false" :fn="() => setAttributes('dim', '20')" name="􀆫"></Subplot>
        <Subplot :active="false" :fn="() => setAttributes('dim', '60')" name="􀆮"></Subplot>
      </Plot>
    </Widget>

  </div>

</template>


<style lang="scss" scoped>

.color-menu {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  opacity: 1;

}
</style>