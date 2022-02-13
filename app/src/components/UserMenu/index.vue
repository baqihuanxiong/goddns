<template>
  <v-menu :close-on-content-click="false">
    <template v-slot:activator="{ on, attrs }">
      <v-btn icon v-bind="attrs" v-on="on">
        <v-icon>mdi-account-circle</v-icon>
      </v-btn>
    </template>
    <v-card>
      <v-card-text v-if="userInfo">
        <v-icon>mdi-account-circle</v-icon>
        {{ userInfo.username }}
      </v-card-text>
      <v-card-actions>
        <v-btn v-if="userInfo" text block color="error" plain @click="showLogoutDialog">logout</v-btn>
        <v-btn v-else text block color="primary" plain @click="gotoLogin">login</v-btn>
      </v-card-actions>
    </v-card>
  </v-menu>
</template>

<script>
import api from '@/api';
import Dialog from '@/components/Dialog';

export default {
  name: "user-menu",

  data () {
    return {
      userInfo: null
    }
  },

  created () {
    this.getCurrentUser();
  },

  methods: {
    getCurrentUser() {
      api.user.getCurrentUser().then(res => {
        this.userInfo = res;
      })
    },

    showLogoutDialog() {
      Dialog({
        title: 'Tips',
        content: 'Are you sure to logout?',
        cancelText: 'cancel',
        confirmText: 'ok',
        cancel: () => {},
        confirm: () => {
          window.localStorage.removeItem('jwt');
          this.gotoLogin();
        }
      })
    },

    gotoLogin() {
      this.$router.push('Login');
    }
  }
}
</script>

<style>

</style>