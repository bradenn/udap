<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import List from "udap-ui/components/List.vue";
import ElementHeader from "udap-ui/components/ElementHeader.vue";
import ElementLink from "udap-ui/components/ElementLink.vue";
import ElementPair from "udap-ui/components/ElementPair.vue";
import core from "udap-ui/core";
import {Device} from "udap-ui/device";

const device: Device = core.device()

</script>

<template>
  <List>

    <ElementHeader title="Web APIs"></ElementHeader>
    <List>
      <ElementLink :cb="() => device.notifications.requestPermission()" :icon="device.notifications.granted ? '􀁣' : '􀀳'" button
                   class="d-flex gap-2 py-3" title="Push Notifications">

        <div>

          {{ device.notifications.granted ? 'Granted' : (device.notifications.supported ? 'Denied' : 'Unsupported') }}
        </div>
      </ElementLink>
      <ElementLink :cb="() => device.notifications.appBadge.setBadge(2)" :icon="device.notifications.appBadge.supported ? '􀁣' : '􀀳'" button
                   class="d-flex gap-2 py-3" title="App Badge">

        <div>

          {{
            device.notifications.appBadge.supported ? (device.notifications.granted ? 'Granted' : 'Supported') : 'Unsupported'
          }}
        </div>
      </ElementLink>
      <ElementLink :cb="() => device.notifications.push.requestPermission()" :icon="device.notifications.push.supported ? '􀈠' : '􀀳'" button
                   class="d-flex gap-2 py-3" title="Send Test Notification">

        <div>


        </div>
      </ElementLink>
      <ElementPair :value="device.notifications.push.granted?'Subscribed':'Unconfirmed'" copyable icon="􀅳"
                   title="Push Subscription"></ElementPair>
      <ElementPair :value="device.notifications.push.vapidUrl" copyable icon="􀅳" title="Vapid URL"></ElementPair>


    </List>


  </List>

</template>

<style scoped>

</style>