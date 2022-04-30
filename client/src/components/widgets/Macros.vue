<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>


import Plot from "@/components/plot/Plot.vue";
import Radio from "@/components/plot/Radio.vue";
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Attribute, Entity, Remote} from "@/types";

let remote = inject("remote") as Remote
let preferences = inject('preferences')

interface Zone {
  name: string
  entities: string[]
}

let state = reactive({
  lights: {} as Entity[],
  targets: ["8c1494c3-6515-490b-8f23-1c03b87bde27", "9a3347a7-7e19-4be5-976c-22384c59142a", "c74d427b-5046-4aeb-8195-2efd05d794f8"] as string[],
  loading: true,
  colorMenu: false,
})

onMounted(() => {
  state.loading = true
  handleUpdates(remote)
})

watchEffect(() => handleUpdates(remote))

function handleUpdates(remote: Remote) {
  state.lights = remote.entities.filter((entity: Entity) => state.targets.includes(entity.id))
  state.loading = false
  return remote.entities
}

function setAttributes(key: string, value: string) {
  remote.attributes.filter((a: Attribute) => a.key == key && state.targets.includes(a.entity)).forEach(v => {
    let copy = v
    copy.request = value
    remote.nexus.requestId("attribute", "request", copy, v.entity)

  })
}

</script>

<template>
  <div v-if="!state.loading" class="w-100">
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