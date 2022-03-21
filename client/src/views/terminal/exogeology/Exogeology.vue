<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import Sidebar from "@/components/sidebar/Sidebar.vue";
import SidebarItem from "@/components/sidebar/SidebarItem.vue";</script>

<template>

  <div class="d-flex flex-row align-content-start align-items-start gap">
    <Sidebar class="flex-shrink-0" icon="meteor" name="Exogeology">
      <router-link class="macro-icon-default" draggable="false" to="/terminal/exogeology/earth/">
        <SidebarItem icon="earth-americas" name="Earth"></SidebarItem>
      </router-link>
      <router-link class="macro-icon-default" draggable="false" to="/terminal/exogeology/moon/">
        <SidebarItem icon="moon" name="Moon"></SidebarItem>
      </router-link>
      <router-link class="macro-icon-default" draggable="false" to="/terminal/exogeology/sol/">
        <SidebarItem icon="sun" name="Sol"></SidebarItem>
      </router-link>
    </Sidebar>
    <div class="w-100">
      <router-view/>
    </div>
  </div>
</template>

<style scoped>

</style>