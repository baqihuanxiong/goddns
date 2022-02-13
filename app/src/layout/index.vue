<template>
  <div>
    <v-app-bar app clipped-left color="primary" dark>
      <v-app-bar-nav-icon @click="toggle = !toggle">
        <v-icon>mdi-menu</v-icon>
      </v-app-bar-nav-icon>
      <v-toolbar-title>
        <router-link to="/overview" style="font-size: 22px; color: white; text-decoration: none">GoDDNS Dashboard</router-link>
      </v-toolbar-title>
      <v-spacer></v-spacer>
      <user-menu></user-menu>
    </v-app-bar>

    <v-navigation-drawer v-model="toggle" app clipped>
      <v-list nav>
        <v-list-item link to="/overview">
          <v-list-item-icon>
            <v-icon>mdi-home</v-icon>
          </v-list-item-icon>
          <v-list-item-title>Overview</v-list-item-title>
        </v-list-item>
        <v-list-item link to="/domains">
          <v-list-item-icon>
            <v-icon>mdi-layers</v-icon>
          </v-list-item-icon>
          <v-list-item-title>Domains</v-list-item-title>
        </v-list-item>
        <v-list-item link to="/records">
          <v-list-item-icon>
            <v-icon>mdi-format-list-bulleted</v-icon>
          </v-list-item-icon>
          <v-list-item-title>Records</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

    <v-main>
      <v-container fluid>
        <router-view></router-view>
      </v-container>
    </v-main>
  </div>
</template>

<script>
import UserMenu from '@/components/UserMenu';

export default {
  name: "index",

  components: {
    UserMenu
  },

  data() {
    return {
      toggle: true
    }
  }
}
</script>

<style scoped>

</style>