<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import ElementLink from "udap-ui/components/ElementLink.vue";
import List from "udap-ui/components/List.vue";
import ElementHeader from "udap-ui/components/ElementHeader.vue";


const services = [
  {
    name: "Core",
    services: [
      {
        name: "App",
        alt: "Mobile PWA",
        address: "https://udap.app"
      },
      {
        name: "API",
        alt: "Application Program Interface",
        address: "https://api.udap.app"
      },
      {
        name: "CDN",
        alt: "Content Delivery Network",
        address: "https://static.udap.app"
      }
    ]
  },
  {
    name: "Infrastructure",
    services: [
      {
        name: "Proxmox",
        alt: "Virtual Environment",
        address: "https://proxmox.udap.app"
      },

    ]
  },
  {
    name: "Network",
    services: [
      {
        name: "Access Point",
        alt: "Wi-Fi Access Point",
        address: "https://10.0.1.8"
      },
    ]
  },
  {
    name: "Entertainment",
    services: [
      {
        name: "Minecraft Server Map",
        alt: "Udap-MC Dynmap",
        address: "https://mc-map.udap.app"
      },
    ]
  },
  {
    name: "Out of Band",
    services: [
      {
        name: "Root Certificate",
        alt: "Udap root certificate",
        address: "http://10.0.2.1/UDAP_CA.der"
      },
    ]
  }
]


function openOutside(link: string) {
  if (link.includes('https://udap.app')) {
    alert("You are already here.")
    return
  }
  if (("standalone" in window.navigator)) {
    window.open(link, '_blank');
  }
}

</script>

<template>
  <List class="">


    <List v-for="group in services">
      <ElementHeader :title="group.name"></ElementHeader>
      <ElementLink v-for="service in group.services" :alt="service.alt" :cb="() => openOutside(service.address)"
                   :title="service.name" icon="􀉣">
        <div class="mono label-c5 label-o3">{{ service.address.replace("https://", "").replace("http://", "") }}</div>
      </ElementLink>

    </List>


  </List>

</template>

<style scoped>

</style>