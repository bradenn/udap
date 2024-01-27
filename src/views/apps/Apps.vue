<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>


import Element from "udap-ui/components/Element.vue";
import Navbar from "@/components/Navbar.vue";
import core from "@/core";
import {RouteMeta} from "vue-router";
import {AppMeta} from "@/views/apps/index";

const router = core.router();

function isApps(): boolean {
  return router.currentRoute.value.fullPath.endsWith("apps") || appName() == "Unknown"
}

const appName: () => string = () => {
  let meta: RouteMeta = router.currentRoute.value.meta
  if (!meta) {
    return "Unknown"
  }
  let appMeta: AppMeta = meta as unknown as AppMeta
  if (!appMeta || Object.keys(appMeta).length == 0) {
    return "Unknown"
  }

  return appMeta.name
}

const apps = [
  {
    name: "Monitor",
    link: "/monitor",
    icon: '􁂥'
  },
  {
    name: "Thermostat",
    link: "/thermostat",
    icon: '􁁋'
  },
  {
    name: "Lasers",
    link: "/sentry",
    icon: '􀑃'
  },
  {
    name: "Lights",
    link: "/lights",
    icon: '􁓼'
  },
  {
    name: "Todo",
    link: "/apps/todo",
    icon: '􀷾'
  },
  {
    name: "Message",
    link: "/apps/message",
    icon: '􀈠'
  },
  {
    name: "Sensors",
    link: "/apps/sensors",
    icon: '􁔊'
  },
  {
    name: "Actions",
    link: "/apps/actions",
    icon: '􀒘'
  },
  {
    name: "Conversions",
    link: "/apps/conversions",
    icon: '􀅮'
  },
  {
    name: "Movies",
    link: "/apps/movies",
    icon: '􀜤'
  },
  {
    name: "Network",
    link: "/apps/network",
    icon: '􁆬'
  },
]

</script>

<template>
  <div class="d-flex gap-1 flex-column scrollable-fixed">
    <Element v-if="!isApps()" class="d-flex justify-content-end">
      <Navbar :title="appName()" back="/apps">
      </Navbar>
    </Element>
    <div class="scrollable-fixed">
      <router-view></router-view>
    </div>
  </div>
</template>

<style lang="scss">

.app-list {
  height: 100%;
  display: grid;
  padding: 1rem;
  grid-template-rows: repeat(6, minmax(1rem, 1fr));
  grid-template-columns: repeat(4, minmax(2rem, 6rem));
  grid-gap: 1.5rem;
}
</style>