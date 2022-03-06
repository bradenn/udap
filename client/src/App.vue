<!-- Copyright (c) 2022 Braden Nicholson -->

<script setup>
import {reactive} from "vue";
import {Attr, Preference} from "./preferences";


let local = reactive({
  background: new Preference(Attr.Background).get(),
  theme: new Preference(Attr.Theme).get(),
})
</script>

<template>
  <div :class="`theme-${local.theme}`" class="root blurs-8">
    <img :src="`/public/custom/${local.background}@4x.png`" alt="Background" class="backdrop"/>
    <router-view/>
  </div>
</template>


<style scoped>

</style>