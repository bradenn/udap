<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import Element from "udap-ui/components/Element.vue";
import List from "udap-ui/components/List.vue";
import {onMounted, reactive, watchEffect} from "vue";
import core from "@/core";

import ElementHeader from "udap-ui/components/ElementHeader.vue";


const remote = core.remote()
const state = reactive({

  ready: false
})

onMounted(() => {

})

watchEffect(() => {

  return remote.attributes
})


</script>

<template>
  <Element v-if="state.ready" class="h-100">
    <List>

      <ElementHeader style="padding-bottom: 0; margin-bottom: 0 !important;" title="Todo"></ElementHeader>
      <List>

      </List>

    </List>
  </Element>
</template>

<style>

</style>