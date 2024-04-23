<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import Element from "udap-ui/components/Element.vue";
import List from "udap-ui/components/List.vue";
import Navbar from "@/components/Navbar.vue";
import core from "@/core";
import {onMounted, reactive, watchEffect} from "vue";
import {Attribute} from "udap-ui/types";


const remote = core.remote()

interface Route {
  name: string
  address: string
}

const state = reactive({
  name: "" as string,
  routes: [] as Route[],
})

</script>

<template>
  <List>

  </List>

</template>

<style lang="scss" scoped>

</style>