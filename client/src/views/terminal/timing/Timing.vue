<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import Sidebar from "@/components/sidebar/Sidebar.vue";
import SidebarItem from "@/components/sidebar/SidebarItem.vue";</script>
<template>
  <div class="h-75">
    <Sidebar class="flex-shrink-0 my-1" horizontal icon="stopwatch" name="Timing" small>
      <div class="w-100"></div>
      <router-link class="macro-icon-default" draggable="false" to="/terminal/timing/stopwatch">
        <SidebarItem icon="stopwatch" name="Stopwatch"></SidebarItem>
      </router-link>
      <router-link class="macro-icon-default" draggable="false" to="/terminal/timing/timer">
        <SidebarItem icon="hourglass" name="Timer"></SidebarItem>
      </router-link>
      <router-link class="macro-icon-default" draggable="false" to="/terminal/timing/conversion">
        <SidebarItem icon="0" name="Conversion"></SidebarItem>
      </router-link>
    </Sidebar>
    <router-view></router-view>
  </div>
</template>

<style scoped>

</style>