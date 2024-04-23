<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import Element from "udap-ui/components/Element.vue";
import List from "udap-ui/components/List.vue";
import ElementToggle from "udap-ui/components/ElementToggle.vue";
import {onMounted, reactive, watchEffect} from "vue";
import core from "@/core";
import {Attribute} from "@/types";
import attributeService from "@/services/attributeService";

import ElementHeader from "udap-ui/components/ElementHeader.vue";

// function requestNotificationPermission() {
//   Notification.requestPermission((permission) => {
//     if (permission === 'granted') {
//       // You can now use the Badging API
//       console.log("granted!!!!")
//     }
//   })
//
//   console.log("Sent?")
// }

interface Item {
  itemId: string
  name: string
  done: boolean
  delete: boolean
}

interface List {
  items: Item[]
}

const remote = core.remote()
const state = reactive({
  listAttribute: {} as Attribute,
  list: {} as List,
  todo: [] as Item[],
  done: [] as Item[],
  selected: "",
  selectedItem: {} as Item,
  message: "",
  ready: false
})

onMounted(() => {
  updateList()
})

watchEffect(() => {
  updateList()
  return remote.attributes
})

watchEffect(() => {
  if (state.selected != "") {
    state.selectedItem.name = state.message
    state.todo = state.todo.map(t => t.itemId === state.selectedItem.itemId ? state.selectedItem : t)
  }
  return state.message
})

function updateList() {
  let listAttribute = remote.attributes.find(a => a.key === "todo-list")
  if (!listAttribute) return

  state.listAttribute = listAttribute
  let list = JSON.parse(listAttribute.value) as List
  if (!list) return;
  state.list = list
  state.todo = state.list.items.filter(i => !i.done)
  state.done = state.list.items.filter(i => i.done)
  state.ready = true
}

function addItem(item: string) {
  if (item == "") return
  let val = JSON.stringify({
    name: item,
    done: false,
    itemId: "",
    delete: false
  })
  let proto = {
    entity: state.listAttribute.entity,
    key: "todo-push",
    value: val,
    request: val
  } as Attribute
  attributeService.request(proto).then(() => {
    state.message = ""
  })
}

function toggleSelection(id: string) {
  if (state.selected == id) {
    state.selected = ""
  } else {
    state.selected = id
    console.log(state.selected)
    let selectedItem = state.list.items.find(l => l.itemId === id)
    if (!selectedItem) return
    state.selectedItem = selectedItem
    state.message = selectedItem.name || ""
  }
}

function changeItem(item: Item) {
  if (item.itemId == "") return
  let val = JSON.stringify({
    name: item.name,
    done: item.done,
    itemId: item.itemId,
    delete: item.delete
  })
  let proto = {
    entity: state.listAttribute.entity,
    key: "todo-push",
    value: val,
    request: val
  } as Attribute
  attributeService.request(proto).then(() => {
    state.selectedItem = {} as Item
    state.selected = ""
    state.message = ""
  })
}

function toggleItem(item: Item) {
  if (item.itemId == "") return
  let val = JSON.stringify({
    name: item.name,
    done: !item.done,
    itemId: item.itemId,
    delete: false
  })
  let proto = {
    entity: state.listAttribute.entity,
    key: "todo-push",
    value: val,
    request: val
  } as Attribute
  attributeService.request(proto).then(() => {
    state.message = ""
  })
}

function deleteItem(item: Item) {
  if (item.itemId == "") return
  let val = JSON.stringify({
    name: item.name,
    done: item.done,
    itemId: item.itemId,
    delete: true
  })
  let proto = {
    entity: state.listAttribute.entity,
    key: "todo-push",
    value: val,
    request: val
  } as Attribute
  console.log("deleting")
  attributeService.request(proto).then(() => {
    state.message = ""
  })
}

function send() {
  if (state.selected == "") {
    addItem(state.message)
  } else {
    state.selectedItem.name = state.message
    changeItem(state.selectedItem)
  }
}


</script>

<template>
  <Element v-if="state.ready" class="h-100" style="height: 42rem !important;">
    <List>

      <ElementHeader style="padding-bottom: 0; margin-bottom: 0 !important;" title="Todo"></ElementHeader>
      <List>
        <ElementToggle v-for="item in state.todo" :key="item.itemId" :accent="state.selected === item.itemId"
                       :cb="() => toggleItem(item)"
                       :long-cb="() => toggleSelection(item.itemId)"
                       :selected="item.done" :title="item.name" class="todo-item"
                       foreground
                       icon="􀑇"></ElementToggle>
        <Element v-if="state.todo.length == 0" class="align-items-center px-3 d-flex label-o3 label-500 label-m"
                 foreground
                 style="height: 3.25rem">
          Nothing to do
        </Element>
      </List>
      <ElementHeader style="padding-bottom: 0; margin-bottom: 0 !important;" title="Complete"></ElementHeader>
      <List>
        <ElementToggle v-for="item in state.done" :key="item.itemId" :cb="() => toggleItem(item)"
                       :class="`${item.delete?'todo-delted':'todo-item'}`"
                       :long-cb="() => deleteItem(item)"
                       :selected="item.done" :title="item.name" class="todo-item"
                       foreground
                       icon="􀑇"></ElementToggle>

        <Element v-if="state.done.length == 0"
                 class="align-items-center px-3 d-flex label-c5 label-o3 label-500 label-m"
                 foreground style="height: 3.25rem">
          Completed tasks will be shown here
        </Element>
      </List>
      <div
          class="align-items-center px-2 d-flex label-c6 label-o3 label-500 label-m">
        <div v-if="state.done.length > 0" class="todo-item">Long press to delete completed tasks</div>
        <div v-else>&nbsp;</div>
      </div>
      <List>
        <ElementHeader style="padding-bottom: 0; margin-bottom: 0 !important;" title="Add Item"></ElementHeader>
        <Element class=" p-0" foreground style="padding: 0.25rem !important; height: 3rem">
          <List class="h-100" row>
            <input id="cypher" v-model="state.message" autocapitalize="off" autocomplete="off"
                   class=" message w-100 px-2" placeholder="Description"
                   style="" type="text">
            <Element :cb="() => send()" class="d-flex align-items-center justify-content-center"
                     foreground
                     style="width: 4rem;">
              <div class="sf label-c6">􀈠</div>
            </Element>
          </List>
        </Element>
      </List>
    </List>
  </Element>
</template>

<style>
.todo-item {
  animation: todoAnimate 250ms ease-out forwards;
}

@keyframes todoAnimate {
  from {
    transform: scale(0.96);
  }
  todo {
    transform: scale(1);
  }
}

.message {
  background-color: transparent;
  outline: none !important;
  border: none;

  color: rgba(255, 255, 255, 0.5)

}
</style>