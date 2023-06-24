<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import Element from "udap-ui/components/Element.vue";
import ElementLink from "udap-ui/components/ElementLink.vue";
import List from "udap-ui/components/List.vue";
</script>

<template>
  <div class="d-flex flex-column gap-3 px-2">

    <Element>
      <List>
        <ElementLink icon="􀣋" title="General" to="/home/settings/general"></ElementLink>
        <ElementLink icon="􀉣" title="Connection" to="/home/settings/connection"></ElementLink>
      </List>
    </Element>

    <Element>
      <List>
        <ElementLink icon="􀏅" title="Background" to="/home/settings/background"></ElementLink>
      </List>
    </Element>
    <Element>
      <List>
        <ElementLink icon="􀟗" title="Entities" to="/home/settings/entities"></ElementLink>
      </List>
    </Element>
  </div>
</template>

<style scoped>
</style>