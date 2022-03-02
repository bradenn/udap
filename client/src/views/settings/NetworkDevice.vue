<template>
  <div :class="active?'surface':'element'">
    <div class="h-bar justify-content-start align-items-center align-content-center pb-1">
      <div class="label-xxs label-o2 label-w600">􀄭</div>
      <div class="label-xxs label-o4 label-w500">&nbsp;&nbsp;{{ device.name === "" ? "Unnamed" : device.name }}</div>
      <div class="fill"></div>
      <div class="h-bar gap label-xxs label-o3 label-w400 px-2">
        <div class="label-xxs label-o3">{{ subtitle }}</div>
      </div>
    </div>


    <div class="v-bar  g-a px-1 py-1">
      <div class="h-bar align-items-center justify-content-between">
        <div class="label-xs label-o4">{{ device.ipv4 }}&nbsp;&nbsp;</div>
      </div>
      <div class="label-xs label-o2">{{
          device.mac
        }}
      </div>
    </div>

    <div class="h-bar gap ">
      <div class="toolbar-button label-o4 px-2 flex-fill">􀒗&nbsp;Mitigate</div>
      <Context class="toolbar-button label-o4 px-2 flex-fill" icon="􀣋" name="Configure">
        <div class="w-25"></div>
        <div class="w-75">
          <div class="v-bar" style="padding-top: 4.5rem;">
            <Header :target="{name: 'Device', icon: '􀉤'}"></Header>
            <div class="v-bar gap flex-shrink-1" style="width: 16rem">
              <div v-for="attribute in Object.values(this.session)"
                   :key="attribute.key"
                   class="element" v-on:click.stop>
                <div class="h-bar justify-content-start align-items-center align-content-center">
                  <div class="label-xxs label-o2 label-w600">{{ attribute.icon }}</div>
                  <div class="label-xxs label-o4 label-w500">&nbsp;&nbsp;{{ attribute.name }}</div>
                  <div class="fill"></div>
                  <div class="h-bar gap label-xxs label-o3 label-w400 px-2">
                    <div class="label-xxs label-o3">{{ attribute.value }}</div>
                  </div>
                </div>
                <div class="d-flex gap">
                  <Dock v-if="attribute.type==='select'" class="surface w-100 mt-1" small>
                    <div v-for="option in attribute.options"
                         :class="`${attribute.value === option.value?'router-link-active':''}`"
                         class="dock-link" href="#">
                      {{ option.name }}
                    </div>
                  </Dock>
                  <div v-else-if="attribute.key==='ipv4'" class="w-100">
                    <div class="d-flex flex-row justify-content-between align-items-end px-0 pt-1">
                      <div class="label-o3 px-1">Lease: 6 days</div>
                      <div class="toolbar-button label-o4 px-2 py-1 pt-1 ">Reserve Address</div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </Context>
    </div>
  </div>
</template>
<script>
import Context from "../../components/Context.vue";
import Header from "../../components/Header.vue";
import Dock from "../../components/Dock.vue";

export default {
  name: 'NetworkDevice',
  components: {Header, Context, Dock},
  props: {
    device: {},
    subtitle: "",
    active: false
  },
  computed: {
    session: function () {
      return {
        hostname: {
          name: "Hostname",
          icon: "􀑏",
          mutable: true,
          key: "hostname",
          value: this.device.hostname || "set",
          type: "text"
        },
        ipv6: {
          name: "IPv6",
          icon: "􀰗",
          key: "IPv6",
          value: this.device.ipv6 || '::/0',
          type: "text"
        },
        ipv4: {
          name: "IPv4",
          icon: "􀐗",
          key: "ipv4",
          value: this.device.ipv4 || '10.0.1.0/24',
          type: "other",
        },
        mac: {
          name: "MAC Address",
          icon: "􀓔",
          key: "mac",
          value: this.device.mac,
          type: "text",
        },
        first: {
          name: "First Seen",
          icon: "􀓔",
          key: "first",
          value: this.$timeSince(new Date(this.device.created)) + " ago",
          type: "text",
        },
        last: {
          name: "Last Seen",
          icon: "􀓔",
          key: "last",
          value: this.$timeSince(new Date(this.device.updated)) + " ago",
          type: "text",
        },
      }
    }
  }
}
</script>
<style lang="scss">
code {
}
</style>
