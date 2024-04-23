<!-- Copyright (c) 2024 Braden Nicholson -->

<script lang="ts" setup>

import List from "udap-ui/components/List.vue";
import Element from "udap-ui/components/Element.vue";
import ElementHeader from "udap-ui/components/ElementHeader.vue";
import ElementLink from "udap-ui/components/ElementLink.vue";
</script>

<template>

  <Element>
    <List>
      <ElementHeader title="Physics"></ElementHeader>
      <List>
        <ElementLink icon="􀎖" title="Torque" to="/apps/utilities/torque"></ElementLink>
      </List>
    </List>
  </Element>
</template>

<style lang="scss" scoped>

</style>