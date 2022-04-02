<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import Sidebar from "@/components/sidebar/Sidebar.vue";
import SidebarItem from "@/components/sidebar/SidebarItem.vue";</script>

<template>
  <div class="d-flex flex-row align-content-start align-items-start gap">

    <Sidebar class="flex-shrink-0" icon="gear" name="Settings">
      <router-link class="" draggable="false" to="/terminal/settings/preferences">
        <SidebarItem icon="bars-progress" name="Preferences"></SidebarItem>
      </router-link>
      <router-link class="" draggable="false" to="/terminal/settings/connection">
        <SidebarItem icon="cloud" name="Connection"></SidebarItem>
      </router-link>
    </Sidebar>
    <div class="flex-grow-1">
      <router-view/>
    </div>
  </div>
</template>

<style scoped>

.sidebar-container {
  width: 14rem;
}
</style>