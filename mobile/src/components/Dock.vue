<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import Element from "udap-ui/components/Element.vue"
import List from "udap-ui/components/List.vue"

const props = defineProps<{
  aspect?: number
  secondary?: boolean
}>()
</script>

<template>
  <div>
    <List row>
      <Element :link="true" :mutable="true" class="d-flex align-items-center justify-content-center gap-1 "
               foreground
               style="height: 3rem" to="/home/dashboard">
        <div class="sf-icon label-c5">􀎟</div>

      </Element>

      <Element :link="true" :mutable="true" class="d-flex align-items-center justify-content-center gap-1 "
               foreground
               to="/home/apps">
        <div class="sf-icon label-c5">􀙣</div>

      </Element>
      <Element :link="true" :mutable="true" class="d-flex align-items-center justify-content-center gap-1 "
               foreground
               to="/home/settings">
        <div class="sf-icon label-c5">􀍟</div>

      </Element>
    </List>
  </div>
</template>

<style scoped>


</style>