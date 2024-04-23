<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>


import Element from "udap-ui/components/Element.vue";
import List from "udap-ui/components/List.vue";
import {onMounted, reactive, watchEffect} from "vue";
import {Attribute, Entity, Module} from "@/types";
import core from "@/core";
import ElementHeader from "udap-ui/components/ElementHeader.vue";
import attributeService from "@/services/attributeService";
import ElementDivider from "@/components/ElementDivider.vue";

interface EntityGroup {
  entity: Entity
  attributes: Attribute[]
}

interface ModuleGroup {
  module: Module,
  entities: EntityGroup[]
}

const state = reactive({
  attributes: [] as Attribute[],
  entities: [] as Entity[],
  modules: [] as ModuleGroup[]
})

const remote = core.remote();
watchEffect(() => {
  update()
  return remote.attributes
})


onMounted(() => {

  update()


  //
  // state.entities = remote.entities.sort((a: Entity, b: Entity) => a.module - b.module)
  // state.attributes = remote.attributes
})

function update() {
  let modules = [] as ModuleGroup[]
  for (let module of remote.modules) {
    let mg = {} as ModuleGroup
    mg.module = module
    mg.entities = [] as EntityGroup[]
    for (let filterElement of remote.entities.filter(e => e.module == module.name)) {
      mg.entities.push({
        entity: filterElement,
        attributes: remote.attributes.filter(a => a.entity == filterElement.id)
      } as EntityGroup)
    }
    modules.push(mg)
  }

  state.modules = modules
}

function deleteAttribute(id: string) {
  attributeService.delete(id).then(() => {
    console.log("deleted")
  }).catch(() => {
  })
}

</script>

<template>

  <List v-if="state.attributes" class="scrollable-fixed" scroll-y>

    <List v-for="module in state.modules.filter(m => m.entities.length > 0)" :key="module.module.id">
      <ElementHeader :title="module.module.name"></ElementHeader>
      <div v-for="entity in module.entities.filter(e => e.attributes.length > 0)" :key="entity.entity.id">


        <List class="mb-1">
          <ElementDivider :text="entity.entity.alias ? entity.entity.alias : entity.entity.name"
                          monogram="E"></ElementDivider>

          <List v-for="attribute in entity.attributes" :key="attribute.id" class="mb-1">
            <List row>
              <Element :to="`/settings/attribute/${module.module.id}`"
                       class="d-flex gap-2 align-items-center px-2 py-2 flex-grow-1 justify-content-between"
                       foreground
                       mutable>

                <div class="d-flex flex-column gap-1">
                  <div class="label-monospace lh-1">{{ attribute.key }}</div>


                </div>

                <div class="d-flex flex-column align-items-end">
                  <div class="label-monospace lh-1 label-o3 label-c5 label label-truncate">{{ attribute.value }}</div>
                  <div class="label-monospace lh-1 label-o2 label-c6">{{ attribute.type }}</div>
                </div>


              </Element>
              <Element :cb="() => deleteAttribute(attribute.id)" class="flex-shrink-1 sf-icon text-danger" foreground
                       style="max-width: 4rem;">􀈒
              </Element>

            </List>
            <Element v-for="key in Object.keys(JSON.parse(attribute.value))" v-if="attribute.value.startsWith('{')"
                     :cb="() => deleteAttribute(attribute.id)"
                     class="flex-shrink-0"
                     foreground
                     style="margin-left: 0.5rem; margin-right: 4.25rem ; padding: 0.25rem 1rem !important;">
              <div class="d-flex justify-content-between">
                <div>{{ key }}</div>
                <div class="label-truncate label-o3">{{ JSON.parse(attribute.value)[key] }}</div>
              </div>
            </Element>


          </List>
        </List>
      </div>
    </List>

  </List>
</template>

<style lang="scss" scoped>
.label-truncate {

  max-width: 12rem;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}
</style>