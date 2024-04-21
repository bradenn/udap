<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import Element from "udap-ui/components/Element.vue";
import ElementLink from "udap-ui/components/ElementLink.vue";
import List from "udap-ui/components/List.vue";
</script>

<template>
  <List>

    <Element>
      <List>
        <ElementLink icon="􀨤" title="Subdomains" to="/apps/network/subdomains"></ElementLink>
        <ElementLink icon="􀨤" title="Graph" to="/apps/network/graph"></ElementLink>
      </List>
    </Element>
  </List>
</template>

<style lang="scss" scoped>

</style>