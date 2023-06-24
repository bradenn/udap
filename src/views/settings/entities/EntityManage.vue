<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import core from "@/core";
import Element from "udap-ui/components/Element.vue";
import ElementHeader from "udap-ui/components/ElementHeader.vue";
import List from "udap-ui/components/List.vue";
import type {Attribute, Entity} from "udap-ui/types"
import {onMounted, reactive} from "vue";

const router = core.router()

const state = reactive({
  entityId: "",
  entity: {} as Entity,
  attributes: [] as Attribute[]
})

const remote = core.remote()

onMounted(() => {
  let eid = router.currentRoute.value.params["entityId"]
  if (!eid) return
  state.entityId = eid as string
  let e = remote.entities.find(e => e.id == state.entityId)
  if (!e) return;
  state.entity = e
  state.attributes = remote.attributes.filter(a => a.entity == state.entityId)

})

</script>

<template>
  <div class="d-flex flex-column gap-2 mx-2">
    <Element class="p-3">
      <div class="label-c3 lh-1 label-o5 d-flex">
        <div class="sf-icon label-c5" style="width: 15px; margin-right: 6px">
          {{ state.entity.icon }}
        </div>
        {{ state.entity.alias ? `${state.entity.alias}ª` : state.entity.name }}
      </div>
      <div class="label-c5 lh-1 label-o3">{{ state.entity.type }}</div>
    </Element>
    <ElementHeader class="px-3 pb-0 m-0" title="Attributes"></ElementHeader>
    <Element>
      <List>
        <Element v-for="attr in state.attributes" :foreground="true" :mutable="true">
          <div class="label-c3 lh-1 label-o5 d-flex">
            {{ attr.key }}
          </div>
          <div class="label-c5 lh-1 label-o3">{{ attr.value }}</div>
        </Element>
      </List>
    </Element>
  </div>
</template>

<style lang="scss" scoped>
.entity {

}
</style>