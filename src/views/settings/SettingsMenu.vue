<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import Element from "udap-ui/components/Element.vue";
import ElementLink from "udap-ui/components/ElementLink.vue";
import List from "udap-ui/components/List.vue";
import Title from "udap-ui/components/Title.vue";
</script>

<template>
  <div class="d-flex flex-column gap-3">

    <Element>
      <List class="gap-3" scroll-y style="height: 72vh">
        <div>
          <Title title="App"></Title>
          <List>
            <ElementLink icon="􀣋" title="General" to="/home/settings/general"></ElementLink>
            <ElementLink icon="􀉣" title="Connection" to="/home/settings/connection"></ElementLink>
            <ElementLink icon="􁇵" title="Configuration" to="/home/settings/config"></ElementLink>
          </List>
        </div>
        <div>
          <Title title="Appearance"></Title>
          <List>
            <ElementLink icon="􀦳" title="Home" to="/home/settings/home"></ElementLink>
            <ElementLink icon="􀏅" title="Background" to="/home/settings/background"></ElementLink>
          </List>
        </div>
        <div>
          <Title title="System"></Title>
          <List>
            <ElementLink icon="􁓽" title="Entities" to="/home/settings/entities"></ElementLink>
            <ElementLink icon="􁚊" title="Modules" to="/home/settings/modules"></ElementLink>
            <ElementLink icon="􀏭" title="Endpoints" to="/home/settings/endpoints"></ElementLink>
            <ElementLink icon="􀟼" title="Zones" to="/home/settings/zones"></ElementLink>
          </List>
        </div>
      </List>
    </Element>

  </div>
</template>

<style scoped>
</style>