<script>

import Header from "../components/Header.vue";
import Device from "../components/device/Device.vue";

export default {
  name: "Network",
  components: {Device, Header},
  data() {
    return {
      header: {
        icon: "􀌆",
        name: "Settings"
      }
    }
  },
  watch: {},
  computed: {
    manifest: function () {
      let sys = this.$root.session.metadata.system
      return this.$root.devices.map(c => {
        let d = c
        if (d.ipv4 === sys.ipv4) {
          d.controller = true
          d.mac = sys.mac
        } else {
          d.controller = false
        }
        return d
      })
    }
  },
  created() {

  },
  methods: {},
}


</script>


<template>
  <div class="d-flex justify-content-center">
    <div>
      <div style="min-width:25vw;">
        <Header :target="{
        icon: '􀆪',
        name: 'Networks'
        }"></Header>
        <div class="context-container context-container-md gap-1">

          <div v-for="network in this.$root.networks">
            <div class="element">
              <div class="h-bar justify-content-start align-items-center align-content-center">
                <div class="label-xxs label-o2 label-w600">􀉣</div>
                <div class="label-xxs label-o4 label-w500 text-capitalize">&nbsp;&nbsp;{{ network.name }}</div>
                <div class="fill"></div>
                <div class="h-bar gap label-xxs label-o3 label-w400 px-1">
                  <div class="label-xxs label-o3">{{ network.mask }}</div>
                </div>
              </div>
              <div class="h-bar justify-content-start align-items-center align-content-center">
                <div class="label-ys label-o3 label-w400">DNS</div>
                <div class="label-ys label-o2 label-w400 text-capitalize">&nbsp;&nbsp;{{
                    network.dns.replace(",", ", ")
                  }}
                </div>
              </div>

            </div>


          </div>
        </div>
      </div>
    </div>
    <div>
      <Header :target="{
        icon: '􀞿',
        name: 'Devices'
        }"></Header>
      <div class="context-container context-container-md gap-1">
        <div v-for="device in manifest">
          <Device :device="device"></Device>
          <div class="label-xs label-o6"></div>
          <div></div>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>


</style>
