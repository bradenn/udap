<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import Element from "udap-ui/components/Element.vue";
import ElementStatus from "udap-ui/components/ElementStatus.vue";
import List from "udap-ui/components/List.vue";
import Title from "udap-ui/components/Title.vue";
</script>

<template>
  <div class="d-flex gap-1 flex-column">
    <Element>
      <Title title="Internal"></Title>
      <List>
        <ElementStatus icon="􀁣" title="Core"></ElementStatus>
        <ElementStatus icon="􀁣" title="Proxmox"></ElementStatus>
        <ElementStatus icon="􀁣" title="Router"></ElementStatus>
        <ElementStatus icon="􀁣" title="Reverse Proxy"></ElementStatus>
        <ElementStatus icon="􀁣" title="DNS Resolution"></ElementStatus>
      </List>
    </Element>
    <Element>
      <Title title="External"></Title>
      <List>
        <ElementStatus icon="􀁣" title="Internet"></ElementStatus>
        <ElementStatus icon="􀁣" title="Out-of-Band Ad-Hoc"></ElementStatus>
        <ElementStatus icon="􀁣" title="Out-of-Band Mux"></ElementStatus>
      </List>
    </Element>
  </div>

</template>

<style scoped>

</style>