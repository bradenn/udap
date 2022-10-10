<!-- Copyright (c) 2022 Braden Nicholson -->

<script setup>

import Atom from "@/views/terminal/periodic/Atom.vue";
import {reactive} from "vue";

const state = reactive({
  atom: {
    protons: 12,
    neutrons: 12,
    electrons: 64,
  }
})

</script>

<template>
  <div class="d-flex h-100 justify-content-center">
    <Atom :electrons="state.atom.electrons" :neutrons="state.atom.neutrons" :protons="state.atom.protons"></Atom>
  </div>
</template>

<style scoped>

</style>