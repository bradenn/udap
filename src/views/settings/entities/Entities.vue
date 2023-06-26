<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import core from "@/core";
import useEntities from "udap-ui/composables/entities"
import {Remote} from "udap-ui/remote";
import Element from "udap-ui/components/Element.vue";
import ElementHeader from "udap-ui/components/ElementHeader.vue";
import List from "udap-ui/components/List.vue";
import Entity from "udap-ui/components/EntityView.vue";

const remote: Remote = core.remote() as Remote
const state = useEntities(remote)

</script>

<template>
  <div v-if="state.loaded" class="mx-2">

    <div class="d-flex flex-column gap-2">
      <Element class="d-flex gap-1 justify-content-between">
        <Element :foreground="true" class="">
          <div class="d-flex justify-content-between">
            <div>
              <div class="notches align-items-start">
                <div class="d-flex align-items-center">
                  <div class="notch-split-cont">
                    <div class="notch-split active"></div>
                    <div class="notch-split"></div>
                  </div>
                  &nbsp;
                  <div class="label-c7 ">Registered</div>
                </div>
                <div class="d-flex align-items-center">
                  <div class="notch-split-cont">
                    <div class="notch-split active"></div>
                    <div class="notch-split warn"></div>
                    <div class="notch-split"></div>
                  </div>
                  &nbsp;
                  <div class="label-c7 ">Updating</div>
                </div>
                <div class="d-flex align-items-center">
                  <div class="notch-split-cont">
                    <div class="notch-split"></div>
                    <div class="notch-split warn"></div>
                  </div>
                  &nbsp;
                  <div class="label-c7">Fault</div>
                </div>
              </div>
            </div>
            <div class="d-flex gap-2 align-items-center">

              <!--              <div style="width: 1px; height: 0.7rem; background-color: rgba(255,255,255,0.1)"></div>-->

            </div>
          </div>
        </Element>
        <div class="d-flex gap-1">
          <Element :foreground="true" :mutable="true" style="width: 5rem">
            <div class="d-flex flex-column gap-0 h-100 align-items-start justify-content-center "
                 @touchstart="(e) => {e.preventDefault(); state.toggleGroup()}">
              <div class="label-c7 label-w500 label-o3  lh-1">Group By</div>
              <div class="label-c5 label-w600 text-capitalize  lh-1">
                {{ state.groupBy }}
              </div>
            </div>
          </Element>
          <Element :foreground="true" :mutable="true" style="width: 7rem">
            <div class="d-flex flex-column gap-0 h-100  align-items-start justify-content-center"
                 @touchstart="(e) => {e.preventDefault(); state.toggleFilter()}">
              <div class="label-c7 label-w500 label-o3  lh-1">Sort By</div>
              <div class="label-c5 label-w600 text-capitalize  lh-1">
                {{ state.filter }}
              </div>
            </div>
          </Element>
        </div>
      </Element>
      <div v-for="group in state.groups" :key="group.name">
        <ElementHeader :title="group.name" class="px-3"></ElementHeader>
        <Element>
          <List class="">
            <Entity v-for="entity in group.entities" :key="entity.id" :entity="entity" :foreground="true"
                    :mutable="true"
                    class=""></Entity>
          </List>
        </Element>
      </div>
    </div>

    <Element v-if="false" class="">
      <List class="">
        <Element :foreground="true">
          <div class="d-flex justify-content-between">
            <div>
              <div class="notches align-items-start">
                <div class="d-flex align-items-center">
                  <div class="notch-split-cont">
                    <div class="notch-split active"></div>
                    <div class="notch-split"></div>
                  </div>
                  &nbsp;
                  <div class="label-c7 ">Registered</div>
                </div>
                <div class="d-flex align-items-center">
                  <div class="notch-split-cont">
                    <div class="notch-split active"></div>
                    <div class="notch-split warn"></div>
                    <div class="notch-split"></div>
                  </div>
                  &nbsp;
                  <div class="label-c7 ">Updating</div>
                </div>
                <div class="d-flex align-items-center">
                  <div class="notch-split-cont">
                    <div class="notch-split"></div>
                    <div class="notch-split warn"></div>
                  </div>
                  &nbsp;
                  <div class="label-c7">Faults</div>
                </div>
              </div>
            </div>
            <div class="d-flex gap-2 align-items-center">
              <div class="d-flex gap-1">
                <div class="label-c5 label-w500 label-o3">Sort By</div>
                <div class="label-c5 label-w400"
                     style="color:rgba(255,128,12,0.6); text-decoration: underline; text-underline-offset: 2px;">Module
                </div>
              </div>
              <div style="width: 1px; height: 0.7rem; background-color: rgba(255,255,255,0.1)"></div>
              <div class="d-flex gap-1">
                <div class="label-c5 label-w500 label-o3">Filter</div>
                <div class="label-c5 label-w400 label-o2"
                     style="color:rgba(255,128,12,0.6); text-decoration: underline; text-underline-offset: 2px;">
                  Recent
                </div>
              </div>
            </div>
          </div>
        </Element>
        <Entity v-for="entity in state.entities" :entity="entity" :foreground="true" :mutable="true"
                class=""></Entity>
      </List>
    </Element>
  </div>
</template>

<style lang="scss" scoped>
.entity {

}
</style>