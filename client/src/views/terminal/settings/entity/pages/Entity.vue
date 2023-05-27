<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {onMounted, reactive, watchEffect} from "vue";
import type {Attribute, Entity} from "@/types";
import EntityDom from "@/components/Entity.vue";
import AttributeDom from "@/components/Attribute.vue";
import Button from "@/components/Button.vue";
import core from "@/core";
import moment from "moment";
import type {Notify} from "@/notifications";
import FixedScroll from "@/components/scroll/FixedScroll.vue";
import Confirm from "@/components/plot/Confirm.vue";

const router = core.router()
const remote = core.remote()

interface Log {
  key: string
  time: string
  date: string
  old: string
  request: string
  who: string
}

const state = reactive({
  entity: {} as Entity,
  attributes: [] as Attribute[],
  showEntity: false,
  logs: [] as Log[],
  loaded: false,
})

onMounted(() => {
  load()
})

function load() {
  const id = router.currentRoute.value.params["entityId"];
  if (!id) return;
  const sr = remote.entities.find(s => s.id === id)
  if (!sr) return
  state.entity = sr
  const attributes = remote.attributes.filter(s => s.entity === state.entity.id)
  if (!attributes) return
  state.attributes = attributes
  sync()
  state.loaded = true;
}

interface KeyPayload {
  icon: string
  unit: string
}

const typeMap = new Map<string, KeyPayload>([
  ["on", {icon: "􀆨", unit: ""}],
  ["dim", {icon: "􀇮", unit: "%"}],
  ["cct", {icon: "􀆭", unit: "K"}],
  ["hue", {icon: "􀟗", unit: "°"}],
  ["api", {icon: "􁉢", unit: ""}],
  ["online", {icon: "􀩲", unit: ""}]])

watchEffect(() => {
  load()

  return remote.attributes
})

watchEffect(() => {

  return remote.attributeLogs
})

function sync() {
  let logs = [] as Log[]
  const logRepo = remote.attributeLogs.sort((a, b) => new Date(b.time).valueOf() - new Date(a.time).valueOf())
  if (remote.attributes.length <= 0) return
  for (let i = 0; i < remote.attributeLogs.length; i++) {
    for (let j = 0; j < state.attributes.length; j++) {
      if (logRepo[i].attribute === state.attributes[j].id) {
        let a = state.attributes[j]
        let l = logRepo[i]
        if (typeMap.has(a.key)) {
          let data = typeMap.get(a.key) as KeyPayload
          logs.push({
            key: a.key,
            old: l.from + data.unit,
            request: l.to + data.unit,
            date: moment(l.time).format("dddd, HH:MM:ss").replace(moment().format("dddd"), "Today"),
            time: moment(l.time).fromNow()
          } as Log)
        }
      }
    }
  }

  state.logs = logs
}

function goBack() {
  router.push("/terminal/settings/entities")
}

const notify: Notify = core.notify()

function timeSince(time: string): string {
  return moment(time).fromNow()
}

function showEntity() {

}

function calcDuration(min: number) {
  return `${Math.floor(min / 60)}h ${min % 60}m `
}


function parseDate(dt: string) {


  return moment().from(dt).toString()
}

</script>

<template>
  <div class="h-100">

    <div v-if="state.loaded" class="d-flex gap-1 h-100">
      <div class="d-flex justify-content-start" style="height: 1.6rem">
        <div class="label-w500 label-c1 text-accent element subplot py-1 px-3"
             @click="() => {goBack()}">􀆉 Back
        </div>

      </div>
      <div class="d-flex flex-column gap-1 flex-grow-1">

        <div class="d-flex align-items-center">
          <div class="label-xs label-w300 label-o6 px-1">
            {{ state.entity.icon }}
          </div>
          <div class="label-sm label-w600 label-o6">
            {{
              state.entity.alias || state.entity.name.slice(0, 24)
            }}
          </div>
        </div>
        <EntityDom :entity="state.entity" noselect></EntityDom>
        <div>
          <Button :active="true" class="element flex-grow-1" icon="􀍟"
                  style="height: 1.8rem" text="Manage"
                  :to="`/terminal/settings/entities/${state.entity.id}/manage`"></Button>
        </div>
        <div class="d-flex gap-1">
          <Button :active="true" class="element flex-grow-1" icon="􀈑"
                  style="height: 1.8rem" text="Delete"
                  @click="() => {}"></Button>
        </div>
        <div class="element d-flex flex-column gap-1 p-2">
          <div class="d-flex justify-content-between">
            <div class="label-c2 label-w500 label-o4">Module</div>
            <div class="label-c2 label-w500 label-o4">
              {{ state.entity.module }}
            </div>
          </div>
          <div class="d-flex justify-content-between">
            <div class="label-c2 label-w500 label-o4">Registered</div>
            <div class="label-c2 label-w500 label-o4">
              {{ timeSince(state.entity.created) }}
            </div>
          </div>
          <div class="d-flex justify-content-between">
            <div class="label-c2 label-w500 label-o4">Updated</div>
            <div class="label-c2 label-w500 label-o4">
              {{ timeSince(state.entity.updated) }}
            </div>
          </div>
        </div>
      </div>

      <div class="d-flex flex-column gap-1 flex-grow-1" style="max-width: 40%">
        <div class="d-flex align-items-center">
          <div class="label-sm label-w200 label-o6 px-1"></div>
          <div class="label-sm label-w600 label-o6">
            Attributes
          </div>
        </div>
        <FixedScroll style="overflow-y: scroll;">
          <div class="d-flex flex-column gap-1">
            <div
                v-for="attr in state.attributes.filter(a => moment(a.lastUpdated).isAfter(moment().subtract(5, `minutes`))).sort((a, b) => a.order - b.order)">
              <AttributeDom :attribute="attr" :noselect="true"></AttributeDom>
            </div>
          </div>
        </FixedScroll>
      </div>
      <div class="d-flex flex-column gap-1 flex-grow-1 h-100">
        <div class="d-flex align-items-end justify-content-between">
          <div class="label-sm label-w600 label-o6">
            Logs
          </div>

          <div class="d-flex align-items-center justify-content-end"
          >
            <div
                class="label-c2 label-w500 label-o3 px-0 lh-1 flex-grow-1 text-end"
                style="width: 6rem">
              {{ state.logs.length }} logs
            </div>
            <div class="v-sep"></div>
            <Confirm
                :active="true"
                :fn="() => {}"
                class="label-w500 label-c1 text-accent element subplot py-1 px-3"
                icon="􀈑"
                style="height: 1.6rem" title="Clear"></Confirm>
          </div>
        </div>
        <FixedScroll
            style="overflow-y: scroll !important; height: 100%; max-height: 100% !important;">
          <div class="d-flex flex-column gap-1">
            <div v-for="log in state.logs">
              <div class="element p-1 d-flex align-items-center gap-1"
                   style="height: 2.2rem">
                <div
                    class="label-c1 label-o2 label-w500 px-1 d-flex justify-content-center"
                    style="width: 1.5rem">
                  {{ typeMap.get(log.key)?.icon }}
                </div>
                <div>
                  <div class="label-c2 label-o4 label-w700 lh-1">
                    {{ log.key }}
                  </div>
                  <div class="label-c3 label-o3 label-w400 lh-sm d-flex gap-1"
                       style="max-width: 7rem;  white-space: nowrap ; overflow: hidden; text-overflow: ellipsis; padding-right: 0.125rem">
                    <div class="label-c3 label-o3 "><code
                    >{{
                        log.old
                      }}</code>
                    </div>
                    <div class="label-c3 label-o2">
                      􀄫
                    </div>
                    <div class="label-c3 label-o3 "><code
                    >{{
                        log.request
                      }}</code>
                    </div>

                  </div>
                </div>
                <div class="flex-fill"></div>
                <div
                    class="d-flex flex-column align-items-end px-1 gap-1">
                  <div class="px-1 label-c2 label-o3 label-w500 lh-1">
                    {{ log.time }}
                  </div>
                  <div class="px-1 label-c3 label-o2 label-w500 lh-1">
                    {{ log.date }}
                  </div>
                </div>

              </div>

            </div>
          </div>
        </FixedScroll>
      </div>
      <!--      <ManageEntity v-if="state.showEntity"-->
      <!--                    :done="() => {state.showEntity = false}"-->
      <!--                    :entity="state.entity"></ManageEntity>-->
      <!--    <EditMacro v-if="state.editMacro" :done="doneEditing"-->
      <!--               :macro="state.macro"-->
      <!--               @click="() => {state.editMacro = false;}"></EditMacro>-->
      <!--    <ShowMacros v-if="state.showMacros"-->
      <!--                :done="() => {state.showMacros = false;}"-->
      <!--                :subroutine="state.subroutine"></ShowMacros>-->

    </div>

  </div>

</template>

<style scoped>

.preview {
  width: 9rem;
}

.grid-element {
  width: 100%;
  height: 100%;
  display: grid;
  grid-column-gap: 0.25rem;
  grid-row-gap: 0.25rem;
  grid-template-rows: repeat(1, 1fr);
  grid-template-columns: repeat(4, 1fr);
}

.generic-grid {
  display: flex;
  justify-content: center;
  grid-column-gap: 0.25rem;
}
</style>