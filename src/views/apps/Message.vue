<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import Element from "udap-ui/components/Element.vue";
import List from "udap-ui/components/List.vue";
import core from "@/core";
import attributeService from "@/services/attributeService";
import {onMounted, reactive} from "vue";

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

const remote = core.remote()

const state = reactive({
  selected: [] as string[],
  from: "",
  message: ""
})

onMounted(() => {
  let token = parseJwt(localStorage.getItem("token") || "") as { id: string }
  state.from = remote.endpoints.find(e => e.id == token.id)?.name || "Unknown"
})

function parseJwt(token: string): any {
  let base64Url = token.split('.')[1];
  let base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
  let jsonPayload = decodeURIComponent(atob(base64).split('').map(function (c) {
    return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
  }).join(''));

  return JSON.parse(jsonPayload);
}

function sendNotification(endpointId: string) {
  let webpush = remote.attributes.find(a => a.key === 'notify')

  if (!webpush) {
    alert("No Web push!")
    return
  }
  webpush.request = JSON.stringify({
    endpointId: endpointId,
    title: `Message from ${state.from}`,
    message: state.message
  })
  attributeService.request(webpush).then((a) => {
    // alert("sent!")
  })
}

function send() {
  state.selected.forEach(s => {
    sendNotification(s)
  })

}

function toggleSelection(id: string) {
  if (state.selected.includes(id)) {
    state.selected = state.selected.filter(s => s != id)
  } else {
    state.selected.push(id)
  }
}

</script>

<template>
  <Element>
    <List>
      <Element v-for="endpoint in remote.endpoints.filter(e => e.notifications)" :key="endpoint.id"
               :accent="state.selected.includes(endpoint.id)"
               :cb="() => toggleSelection(endpoint.id)" foreground
               mutable>
        {{ endpoint.name }}
      </Element>
      <Element class=" p-0" foreground style="padding: 0.25rem !important;">
        <List row>
          <input id="cypher" v-model="state.message" autocapitalize="off" autocomplete="off"
                 class=" message w-100 px-2" placeholder="Message"
                 style="" type="text">
          <Element :cb="send" class="d-flex align-items-center justify-content-center" foreground
                   style="width: 4rem; height: 2rem;">
            <div class="sf label-c6">􀈠</div>
          </Element>
        </List>
      </Element>

    </List>

  </Element>
</template>

<style>
.message {
  background-color: transparent;
  outline: none !important;
  border: none;

  color: rgba(255, 255, 255, 0.5)

}
</style>