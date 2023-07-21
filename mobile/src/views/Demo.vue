<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import Element from "udap-ui/components/Element.vue";
import ElementHeader from "udap-ui/components/ElementHeader.vue";
import List from "udap-ui/components/List.vue";
</script>

<template>
  <Element v-if="true">
    <ElementHeader title="Words">
      <template v-slot:description>Running</template>
    </ElementHeader>
    <List>
      <Element :foreground="true" :mutable="true">
        <ElementHeader title="Demo Item">
          <template v-slot:description>Running</template>
        </ElementHeader>
        <div class="px-2">
          Hello
        </div>
      </Element>
      <Element :foreground="true" :mutable="true">
        <ElementHeader title="Demo Item">
          <template v-slot:description>Running</template>
        </ElementHeader>
        <div class="px-2">
          Hello
        </div>
      </Element>

      <Element :foreground="true" :mutable="true">
        <ElementHeader title="Demo Item">
          <template v-slot:description>Running</template>
        </ElementHeader>
      </Element>
    </List>
  </Element>

</template>

<style lang="scss" scoped>

</style>