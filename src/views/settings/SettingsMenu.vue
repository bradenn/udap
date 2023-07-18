<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import ElementLink from "udap-ui/components/ElementLink.vue";
import List from "udap-ui/components/List.vue";
import ElementHeader from "udap-ui/components/ElementHeader.vue";
</script>

<template>

  <List class="gap-3" scroll-y style="max-height: 80vh">
    <div>
      <ElementHeader title="Appearance"></ElementHeader>
      <List>
        <ElementLink icon="􀦳" title="Home" to="/home/settings/home"></ElementLink>
        <ElementLink icon="􀏅" title="Background" to="/home/settings/background"></ElementLink>
      </List>
    </div>
    <div>
      <ElementHeader title="App"></ElementHeader>
      <List>
        <ElementLink icon="􀣋" title="General" to="/home/settings/general"></ElementLink>
        <ElementLink icon="􀉣" title="Connection" to="/home/settings/connection"></ElementLink>
        <ElementLink icon="􁇵" title="Configuration" to="/home/settings/config"></ElementLink>
      </List>
    </div>

    <div>
      <ElementHeader title="Network"></ElementHeader>
      <List>
        <ElementLink icon="􀉣" title="Web Services" to="/home/settings/webservices"></ElementLink>
      </List>
    </div>
    <div>
      <ElementHeader title="System"></ElementHeader>
      <List>
        <ElementLink icon="􁓽" title="Entities" to="/home/settings/entities"></ElementLink>
        <ElementLink icon="􁚊" title="Modules" to="/home/settings/modules"></ElementLink>
        <ElementLink icon="􀏭" title="Endpoints" to="/home/settings/endpoints"></ElementLink>
        <ElementLink icon="􀟼" title="Zones" to="/home/settings/zones"></ElementLink>
        <ElementLink icon="􀐬" title="Timings" to="/home/settings/timings"></ElementLink>
      </List>
    </div>
  </List>


</template>

<style scoped>
</style>