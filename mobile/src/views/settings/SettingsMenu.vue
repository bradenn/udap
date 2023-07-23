<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import ElementLink from "udap-ui/components/ElementLink.vue";
import List from "udap-ui/components/List.vue";
import ElementHeader from "udap-ui/components/ElementHeader.vue";
import core from "@/core";
import {onMounted, reactive} from "vue";

const router = core.router()

interface Route {
  name: string
  icon: string
  path: string
  category: string
}

interface Group {
  name: string,
  children: Route[]
}

const state = reactive({
  routes: [] as Route[],
  groups: [] as Group[]
})

onMounted(() => {
  let routes = router.getRoutes()
  state.routes = routes.filter(r => r.path.includes("/settings") && r.meta["category"]).map(r => {
    let category: string = (r.meta["category"] || 'ungrouped') as string
    let route: Route = {
      name: r.meta["name"] || r.name,
      icon: r.meta["icon"] || '?',
      path: r.path,
    } as Route
    let target = state.groups.find(g => g.name == category)
    if (!target) {
      state.groups.push({
        name: category,
        children: [route]
      })
    } else {
      target.children.push(route)
      state.groups = state.groups.map(g => g.name == category ? target : g)
    }
    return route
  })
})

</script>

<template>

  <List class="gap-3" scroll-y style="max-height: 80vh">

    <div v-for="group in state.groups">
      <ElementHeader :title="group.name"></ElementHeader>
      <List>
        <ElementLink v-for="path in group.children" :icon="path.icon" :title="path.name" :to="path.path">
        </ElementLink>
      </List>
    </div>

  </List>


</template>

<style scoped>
</style>