<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>

import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Attribute, Entity, Remote} from "@/types";
import Plot from "@/components/plot/Plot.vue";
import Subplot from "@/components/plot/Subplot.vue";

interface Response {
  result: {
    end: number;
    start: number;
    conf: number;
    word: string;
  }[];
  text: string;
}

let remote = inject("remote") as Remote

let state = reactive({
  entity: {} as Entity,
  buffer: {} as Attribute,
  status: {} as Attribute,
  response: {} as Response,
  decoded: {} as Status,
  loaded: false,
})

onMounted(() => {
  load(remote)
})

interface Status {
  recognizer: string
  synthesizer: string
}

watchEffect(() => load(remote))

function load(r: Remote) {
  let target = r.entities.find((entity: Entity) => entity.name === "atlas")
  if (!target) {
    target = {} as Entity
  }
  state.entity = target
  let buffer = r.attributes.find((attribute: Attribute) => attribute.entity === state.entity.id && attribute.key === 'buffer')
  if (!buffer) {
    buffer = {} as Attribute
  }
  state.buffer = buffer
  state.response = JSON.parse(buffer.value) as Response


  let status = r.attributes.find((attribute: Attribute) => attribute.entity === state.entity.id && attribute.key === 'status')
  if (!status) {
    status = {} as Attribute
  }
  state.status = status
  state.decoded = JSON.parse(state.status.value) as Status
  state.loaded = true
  return r
}


</script>

<template>
  <div class="d-flex justify-content-center align-items-start h-100 w-100">

    <div v-if="state.loaded" class="mt-4">
      <Plot :cols="1" :rows="2" style="width: 13rem;" title="Status">
        <Subplot
            :alt="state.decoded.recognizer.toUpperCase()"
            :fn="() => {}" name="Recognizer"></Subplot>
        <Subplot
            :alt="state.decoded.synthesizer.toUpperCase()"
            :fn="() => {}" name="Synthesizer"></Subplot>
        <div>
        </div>
      </Plot>
      <div v-if="state.loaded" class="d-flex">
        <div v-for="word in state.response.result">
          <div :style="`color:rgba(${255-Math.round(word.conf * 255)}, ${Math.round(word.conf * 255)}, 0, 1);`">
            {{ word.word }}&nbsp;
          </div>
        </div>

      </div>
    </div>
    <div v-else>
      Loading
    </div>

  </div>
</template>

<style lang="scss" scoped>
</style>