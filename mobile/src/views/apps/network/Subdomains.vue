<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import Element from "udap-ui/components/Element.vue";
import List from "udap-ui/components/List.vue";
import Navbar from "@/components/Navbar.vue";
import ElementPair from "udap-ui/components/ElementPair.vue";
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

onMounted(() => {
  updateRoutes(remote.attributes)
})

watchEffect(() => {
  return updateRoutes(remote.attributes.filter(r => r.key == "static-routes"))
})
// A little bit simplified version
const groupBy = <T, K extends keyof any>(arr: T[], key: (i: T) => K) =>
    arr.reduce((groups, item) => {
      (groups[key(item)] ||= []).push(item);
      return groups;
    }, {} as Record<K, T[]>);

function sortBySubnet(a: Route, b: Route) {
  let arrA = a.address.split(".")
  let arrB = b.address.split(".")
  return (parseInt(arrA[2]) - parseInt(arrB[2]))
}

function updateRoutes(routes: Attribute[]) {
  let attribute = routes.find(r => r.key == "static-routes")
  if (!attribute) {
    return
  }

  let obj = JSON.parse(attribute.value)
  let domains = Object.keys(obj)
  let addresses = Object.values(obj)
  state.routes = state.routes.filter(f => false)
  for (let i = 0; i < domains.length; i++) {
    state.routes.push({name: domains[i], address: addresses[i]} as Route)
  }

  state.routes = state.routes.filter(r => r.name.endsWith(".app")).sort(sortBySubnet)


  return state.routes
}

</script>

<template>
  <List>
    <Element>
      <Navbar title="subdomains"></Navbar>
    </Element>
    <Element>
      <List scroll-y style=" max-height: 76vh">
        <ElementPair v-for="route in state.routes" :key="route.name" :title="route.name" :value="route.address"
                     icon="􀨤"></ElementPair>
      </List>
    </Element>
  </List>

</template>

<style lang="scss" scoped>

</style>