<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>


import Plot from "@/components/plot/Plot.vue";
import {reactive} from "vue";
import Subplot from "@/components/plot/Subplot.vue";

let state = reactive({
  voice: "default"
})


const voices = [
  {
    name: "Default",
    path: "default"
  },
  {
    name: "slt",
    path: "slt",
  },

  {
    name: "aew",
    path: "aew",
  },
  {
    name: "ahw",
    path: "ahw",
  },
  {
    name: "aup",
    path: "aup",
  },
  {
    name: "awb",
    path: "awb",
  },
  {
    name: "axb",
    path: "axb",
  },
  {
    name: "bdl",
    path: "bdl",
  },
  {
    name: "clb",
    path: "clb",
  },
  {
    name: "eey",
    path: "eey",
  },
  {
    name: "fem",
    path: "fem",
  },
  {
    name: "gka",
    path: "gka",
  },
  {
    name: "jmk",
    path: "jmk",
  },
  {
    name: "ksp",
    path: "ksp",
  },
  {
    name: "ljm",
    path: "ljm",
  },
  {
    name: "rms",
    path: "rms",
  },
  {
    name: "rxr",
    path: "rxr",
  },
  {
    name: "slt_hts",
    path: "slt_hts",
  },


]


function setVoice(path: string) {
  state.voice = path
}

</script>

<template>
  <div class="d-flex">

    <div class="w-100">
      <div class="label-lg mb-1">Voice Selection</div>
      <Plot cols="9" rows="2">
        <Subplot v-for="voice in voices" :active="state.voice === voice.path" :fn="() => setVoice(voice.path)"
                 :name="voice.name"></Subplot>
      </Plot>
    </div>


  </div>
</template>

<style lang="scss" scoped>
</style>