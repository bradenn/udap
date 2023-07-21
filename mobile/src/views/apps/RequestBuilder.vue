<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>


import Element from "udap-ui/components/Element.vue";
import ElementHeader from "udap-ui/components/ElementHeader.vue";
import List from "udap-ui/components/List.vue";
</script>

<template>
  <Element>
    <ElementHeader title="Selected"></ElementHeader>
    <List>
      <List row>
        <Element foreground></Element>
        <Element foreground></Element>
        <Element foreground></Element>
        <Element foreground></Element>
      </List>
      <Element foreground>
      </Element>
    </List>
    <ElementHeader title="Request"></ElementHeader>
    <List>
      <List row>
        <Element foreground></Element>
        <Element foreground></Element>
        <Element foreground></Element>
        <Element foreground></Element>
      </List>
      <Element foreground>
      </Element>
    </List>
  </Element>

</template>

<style>

</style>