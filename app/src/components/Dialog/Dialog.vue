<template>
  <v-dialog v-model="isDialogVisible" max-width="320">
    <v-card>
      <v-card-title>{{ title }}</v-card-title>
      <v-card-text>{{ content }}</v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn text color="grey" @click="cancel">{{ cancelText }}</v-btn>
        <v-btn text color="primary" @click="confirm">{{ confirmText }}</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
export default {
  data () {
    return {
      isDialogVisible: false
    }
  },

  props: {
    title: String,
    content: String,
    cancelText: String,
    confirmText: String,
    handleCancel: Function,
    handleConfirm: Function
  },

  methods: {
    cancel() {
      this.isDialogVisible = false;
      this.handleCancel();
    },

    confirm() {
      this.isDialogVisible = false;
      this.handleConfirm();
    }
  }
}
</script>