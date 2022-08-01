<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>


import Plot from "@/components/plot/Plot.vue";
import Radio from "@/components/plot/Radio.vue";
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Attribute, Entity, Remote, Zone} from "@/types";
import Subplot from "@/components/plot/Subplot.vue";
import Light from "@/components/widgets/Light.vue";
import Widget from "@/components/widgets/Widget.vue";

let remote = inject("remote") as Remote
let preferences = inject('preferences')


let state = reactive({
  zone: {} as Zone,
  lights: {} as Entity[],
  targets: ["8c1494c3-6515-490b-8f23-1c03b87bde27", "9a3347a7-7e19-4be5-976c-22384c59142a", "c74d427b-5046-4aeb-8195-2efd05d794f8"] as string[],
  loading: true,
  colorMenu: false,
})

onMounted(() => {
  state.loading = true
  switchZone("all")
  handleUpdates(remote)
})

watchEffect(() => handleUpdates(remote))

function handleUpdates(remote: Remote) {
  state.lights = remote.entities.filter((entity: Entity) => state.targets.includes(entity.id))
  state.loading = false
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
  <div v-if="!state.loading" class="d-flex flex-column gap-1" size="sm" style="width: 11rem">
    <Plot :cols="3" :rows="1">

      <Subplot :active="state.zone.name === 'all'" :fn="() => switchZone('all')" name="All" title="All"></Subplot>
      <Subplot :active="state.zone.name === 'bedroom'" :fn="() => switchZone('bedroom')" name="Bedroom"
               title="Bedroom"></Subplot>
      <Subplot :active="state.zone.name === 'lor'" :fn="() => switchZone('lor')" name="Lor"
               title="Lor"></Subplot>
      <Subplot :active="state.zone.name === 'kitchen'" :fn="() => switchZone('kitchen')" name="Kitchen"
               title="Kitchen"></Subplot>
    </Plot>
    <Widget v-if="!state.loading" :cols="4" :rows="5" class="d-flex flex-column" size="sm">
      <Light v-for="light in state.lights.slice(0, 5)"
             :key="light.id"
             :entity="light" style="height: 8rem !important;"></Light>
    </Widget>
    <Widget v-if="!state.loading" :cols="4" :rows="3" class="d-flex flex-column" size="sm">

      <Plot :cols="4" :rows="3">
        <Radio :active="false" :fn="() => setAttributes('on', 'false')" title="OFF"></Radio>
        <Radio :active="false" :fn="() => setAttributes('on', 'true')" title="ON"></Radio>
        <Radio :active="false" :fn="() => setAttributes('dim', '20')" title="􀆫"></Radio>
        <Radio :active="false" :fn="() => setAttributes('dim', '60')" title="􀆮"></Radio>
        <Radio :active="false" :fn="() => setAttributes('cct', '2000')" title="2000K"></Radio>
        <Radio :active="false" :fn="() => setAttributes('cct', '2600')" title="2600K"></Radio>
        <Radio :active="false" :fn="() => setAttributes('cct', '6500')" title="6500K"></Radio>
        <Radio :active="false" :fn="() => setAttributes('cct', '7100')" title="7100K"></Radio>
        <Radio :active="false" :fn="() => setAttributes('hue', '200')" title="Blue"></Radio>
        <Radio :active="false" :fn="() => setAttributes('hue', '50')" title="Yellow"></Radio>
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