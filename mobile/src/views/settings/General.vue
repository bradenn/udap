<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import ElementStatus from "udap-ui/components/ElementStatus.vue";
import List from "udap-ui/components/List.vue";
import ElementHeader from "udap-ui/components/ElementHeader.vue";
</script>

<template>
  <List>

    <ElementHeader title="Internal"></ElementHeader>
    <List>
      <ElementStatus icon="􀁣" title="Core"></ElementStatus>
      <ElementStatus icon="􀁣" title="Proxmox"></ElementStatus>
      <ElementStatus icon="􀁣" title="Router"></ElementStatus>
      <ElementStatus icon="􀁣" title="Reverse Proxy"></ElementStatus>
      <ElementStatus icon="􀁣" title="DNS Resolution"></ElementStatus>
    </List>


    <ElementHeader title="External"></ElementHeader>
    <List>
      <ElementStatus icon="􀁣" title="Internet"></ElementStatus>
      <ElementStatus icon="􀁣" title="Out-of-Band Ad-Hoc"></ElementStatus>
      <ElementStatus icon="􀁣" title="Out-of-Band Mux"></ElementStatus>
    </List>

  </List>

</template>

<style scoped>

</style>