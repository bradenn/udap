<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import Element from "udap-ui/components/Element.vue"

const props = defineProps<{
  aspect?: number
  secondary?: boolean
}>()
</script>

<template>
  <div>
    <div class="d-flex flex-row gap-1">
      <Element :foreground="true" :link="true" :mutable="true" class="d-flex align-items-center justify-content-center gap-1"
               style="height: 2.8rem" to="/home/dashboard">
        <div class="sf-icon label-c5">􀎟</div>
        <div class="label-c5">Home</div>
      </Element>
      <Element :foreground="true" :link="true" :mutable="true" class="d-flex align-items-center justify-content-center gap-1"
               to="/home/macros">
        <div class="sf-icon label-c5">􀞿</div>
        <div class="label-c5">Macros</div>
      </Element>
      <Element :foreground="true" :link="true" :mutable="true" class="d-flex align-items-center justify-content-center gap-1"
               to="/home/settings">
        <div class="sf-icon label-c5">􀍟</div>
        <div class="label-c5">Settings</div>
      </Element>
    </div>
  </div>
</template>

<style scoped>


</style>