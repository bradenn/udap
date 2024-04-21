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
      <Navbar title="graphs"></Navbar>
    </Element>
    <div class="pane pane-base">
      <div class="pane">
        <div class="pane-title">Internet</div>
      </div>
      <div class="pane pane-row">
        <div class="pane">
          <div class="pane-title">FireWall</div>
        </div>
        <div class="pane">
          <div class="pane-title">Udap OOB</div>
        </div>
      </div>
      <div class="pane">
        <div class="pane-title">Router</div>
        <div class="pane">Hello</div>
        <div class="pane">Hello<br>me</div>
      </div>
      <div class="pane">Hello</div>
    </div>
  </List>

</template>

<style lang="scss" scoped>
.pane-base {
  backdrop-filter: blur(9px);
}

$pane-padding: 6px;
$pane-step: 6px;
$pane-radius: 16px;

.pane-center {

}

.pane-title {
  font-family: "SF Pro Display", sans-serif;
  padding: 0px 4px;
  line-height: 18px;
  font-weight: 600;
  font-size: 16px;
}

.pane {
  background-color: rgba(255, 255, 255, 0.05);
  box-shadow: inset 0 0 1px 0.3px rgba(255, 255, 255, 0.09);
  border-radius: $pane-radius;
  padding: $pane-padding;
  display: flex;
  width: 100%;
  flex-direction: column;
  gap: $pane-padding/2;
}

.pane-row {
  flex-direction: row;
}

.pane > .pane {
  background-color: rgba(255, 255, 255, 0.05);
  padding: $pane-step;
  border-radius: $pane-radius - $pane-step;
}

.pane > .pane > .pane {
  background-color: rgba(255, 255, 255, 0.05);
  padding: $pane-step * 2;
  border-radius: $pane-radius - $pane-step * 2;
}
</style>