<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import Range from "./Range.vue";
import Toggle from "./Toggle.vue";
import {inject} from "vue";
import type {Attribute} from "@/types";
import attributeService from "@/services/attributeService";

// Define the props for this component
let props = defineProps<{
  attribute: Attribute
  entityId: string
  small: boolean,
}>()

// Inject the remote reference
let remote: any = inject('remote')

// Apply changes made to an attribute
function commitChange(attribute: Attribute) {
  attributeService.request(attribute).then(res => {
    console.log(res)
  }).catch(err => {
    console.log(err)
  })
}

</script>

<template>
  <div>
    <Range v-if="attribute.type === 'range'" :attribute="attribute" :commit="commitChange"
           :small="small"></Range>
    <Toggle v-else-if="attribute.type === 'toggle'" :attribute="attribute" :commit="commitChange"
            :small="small"></Toggle>
  </div>
</template>

<style scoped>

</style>
